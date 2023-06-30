package transformer

import (
	"crypto/rand"
	"testing"

	cidtest "github.com/TrueCloudLab/frostfs-sdk-go/container/id/test"
	objectSDK "github.com/TrueCloudLab/frostfs-sdk-go/object"
	"github.com/TrueCloudLab/frostfs-sdk-go/version"
	"github.com/stretchr/testify/require"
)

func TestChannelTarget(t *testing.T) {
	const maxSize = 100

	ch := make(chan *objectSDK.Object, 10)
	tt := new(testTarget)

	chTarget, _ := newPayloadSizeLimiter(maxSize, NewChannelTarget(ch))
	testTarget, _ := newPayloadSizeLimiter(maxSize, tt)

	ver := version.Current()
	cnr := cidtest.ID()
	hdr := objectSDK.New()
	hdr.SetContainerID(cnr)
	hdr.SetType(objectSDK.TypeRegular)
	hdr.SetVersion(&ver)

	payload := make([]byte, maxSize*2+maxSize/2)
	_, _ = rand.Read(payload)

	expectedIDs := writeObject(t, testTarget, hdr, payload)
	actualIDs := writeObject(t, chTarget, hdr, payload)
	_ = expectedIDs
	_ = actualIDs
	//require.Equal(t, expectedIDs, actualIDs)

	for i := range tt.objects {
		select {
		case obj := <-ch:
			// Because of the split ID objects can be different.
			// However, payload and attributes must be the same.
			require.Equal(t, tt.objects[i].Payload(), obj.Payload())
			require.Equal(t, tt.objects[i].Attributes(), obj.Attributes())
		default:
			require.FailNow(t, "received less parts than expected")
		}
	}

	select {
	case <-ch:
		require.FailNow(t, "received more parts than expected")
	default:
	}
}
