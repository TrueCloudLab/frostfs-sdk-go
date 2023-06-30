package transformer

import (
	"crypto/rand"
	"crypto/sha256"
	"testing"

	cid "github.com/TrueCloudLab/frostfs-sdk-go/container/id"
	cidtest "github.com/TrueCloudLab/frostfs-sdk-go/container/id/test"
	objectSDK "github.com/TrueCloudLab/frostfs-sdk-go/object"
	"github.com/TrueCloudLab/frostfs-sdk-go/user"
	"github.com/TrueCloudLab/frostfs-sdk-go/version"
	"github.com/nspcc-dev/neo-go/pkg/crypto/keys"
	"github.com/stretchr/testify/require"
)

func TestTransformer(t *testing.T) {
	const maxSize = 100

	tt := new(testTarget)

	target, pk := newPayloadSizeLimiter(maxSize, tt)

	cnr := cidtest.ID()
	hdr := newObject(cnr)

	var owner user.ID
	user.IDFromKey(&owner, pk.PrivateKey.PublicKey)
	hdr.SetOwnerID(&owner)

	expectedPayload := make([]byte, maxSize*2+maxSize/2)
	_, _ = rand.Read(expectedPayload)

	ids := writeObject(t, target, hdr, expectedPayload)
	require.Equal(t, 4, len(tt.objects)) // 3 parts + linking object

	var actualPayload []byte
	for i := range tt.objects {
		childCnr, ok := tt.objects[i].ContainerID()
		require.True(t, ok)
		require.Equal(t, cnr, childCnr)
		require.Equal(t, objectSDK.TypeRegular, tt.objects[i].Type())
		require.Equal(t, &owner, tt.objects[i].OwnerID())

		payload := tt.objects[i].Payload()
		require.EqualValues(t, tt.objects[i].PayloadSize(), len(payload))
		actualPayload = append(actualPayload, payload...)

		if len(payload) != 0 {
			cs, ok := tt.objects[i].PayloadChecksum()
			require.True(t, ok)

			h := sha256.Sum256(payload)
			require.Equal(t, h[:], cs.Value())
		}

		switch i {
		case 0, 1:
			require.EqualValues(t, maxSize, len(payload))
		case 2:
			require.EqualValues(t, maxSize/2, len(payload))
		case 3:
			parID, ok := tt.objects[i].ParentID()
			require.True(t, ok)
			require.Equal(t, ids.ParentID, &parID)
		}
	}
	require.Equal(t, expectedPayload, actualPayload)

	t.Run("parent checksum", func(t *testing.T) {
		cs, ok := ids.ParentHeader.PayloadChecksum()
		require.True(t, ok)

		h := sha256.Sum256(expectedPayload)
		require.Equal(t, h[:], cs.Value())
	})
}

func newObject(cnr cid.ID) *objectSDK.Object {
	ver := version.Current()
	hdr := objectSDK.New()
	hdr.SetContainerID(cnr)
	hdr.SetType(objectSDK.TypeRegular)
	hdr.SetVersion(&ver)
	return hdr
}

func writeObject(t *testing.T, target ObjectTarget, header *objectSDK.Object, payload []byte) *AccessIdentifiers {
	require.NoError(t, target.WriteHeader(header))

	_, err := target.Write(payload)
	require.NoError(t, err)

	ids, err := target.Close()
	require.NoError(t, err)

	return ids
}

func BenchmarkTransformer(b *testing.B) {
	hdr := newObject(cidtest.ID())

	b.Run("small", func(b *testing.B) {
		benchmarkTransformer(b, hdr, 8*1024)
	})
	b.Run("big", func(b *testing.B) {
		benchmarkTransformer(b, hdr, 64*1024*1024*9/2)
	})
}

func benchmarkTransformer(b *testing.B, header *objectSDK.Object, payloadSize int) {
	const maxSize = 64 * 1024 * 1024

	payload := make([]byte, payloadSize)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		f, _ := newPayloadSizeLimiter(maxSize, benchTarget{})
		if err := f.WriteHeader(header); err != nil {
			b.Fatalf("write header: %v", err)
		}
		if _, err := f.Write(payload); err != nil {
			b.Fatalf("write: %v", err)
		}
		if _, err := f.Close(); err != nil {
			b.Fatalf("close: %v", err)
		}
	}
}

func newPayloadSizeLimiter(maxSize uint64, nextTarget ObjectTarget) (ObjectTarget, *keys.PrivateKey) {
	p, err := keys.NewPrivateKey()
	if err != nil {
		panic(err)
	}

	return NewPayloadSizeLimiter(Params{
		Key:                    &p.PrivateKey,
		NextTarget:             nextTarget,
		NetworkState:           dummyEpochSource(123),
		MaxSize:                maxSize,
		WithoutHomomorphicHash: true,
	}), p
}

type dummyEpochSource uint64

func (s dummyEpochSource) CurrentEpoch() uint64 {
	return uint64(s)
}

type benchTarget struct{}

func (benchTarget) WriteHeader(object *objectSDK.Object) error {
	return nil
}

func (benchTarget) Write(p []byte) (n int, err error) {
	return len(p), nil
}

func (benchTarget) Close() (*AccessIdentifiers, error) {
	return nil, nil
}

type testTarget struct {
	current *objectSDK.Object
	payload []byte
	objects []*objectSDK.Object
}

func (tt *testTarget) WriteHeader(object *objectSDK.Object) error {
	tt.current = object
	return nil
}

func (tt *testTarget) Write(p []byte) (n int, err error) {
	tt.payload = append(tt.payload, p...)
	return len(p), nil
}

func (tt *testTarget) Close() (*AccessIdentifiers, error) {
	tt.current.SetPayload(tt.payload)
	tt.objects = append(tt.objects, tt.current)
	tt.current = nil
	tt.payload = nil
	return nil, nil // AccessIdentifiers should not be used.
}
