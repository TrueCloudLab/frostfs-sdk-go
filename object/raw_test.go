package object

import (
	"crypto/rand"
	"crypto/sha256"
	"testing"

	"github.com/TrueCloudLab/frostfs-api-go/v2/object"
	"github.com/TrueCloudLab/frostfs-sdk-go/checksum"
	cidtest "github.com/TrueCloudLab/frostfs-sdk-go/container/id/test"
	oid "github.com/TrueCloudLab/frostfs-sdk-go/object/id"
	sessiontest "github.com/TrueCloudLab/frostfs-sdk-go/session/test"
	usertest "github.com/TrueCloudLab/frostfs-sdk-go/user/test"
	"github.com/TrueCloudLab/frostfs-sdk-go/version"
	"github.com/stretchr/testify/require"
)

func randID(t *testing.T) oid.ID {
	var id oid.ID
	id.SetSHA256(randSHA256Checksum(t))

	return id
}

func randSHA256Checksum(t *testing.T) (cs [sha256.Size]byte) {
	_, err := rand.Read(cs[:])
	require.NoError(t, err)

	return
}

func randTZChecksum(t *testing.T) (cs [64]byte) {
	_, err := rand.Read(cs[:])
	require.NoError(t, err)

	return
}

func TestObject_SetID(t *testing.T) {
	obj := New()

	id := randID(t)

	obj.SetID(id)

	oID, set := obj.ID()
	require.True(t, set)
	require.Equal(t, id, oID)
}

func TestObject_SetPayload(t *testing.T) {
	obj := New()

	payload := make([]byte, 10)
	_, _ = rand.Read(payload)

	obj.SetPayload(payload)

	require.Equal(t, payload, obj.Payload())
}

func TestObject_SetVersion(t *testing.T) {
	obj := New()

	var ver version.Version
	ver.SetMajor(1)
	ver.SetMinor(2)

	obj.SetVersion(&ver)

	require.Equal(t, ver, *obj.Version())
}

func TestObject_SetPayloadSize(t *testing.T) {
	obj := New()

	sz := uint64(133)
	obj.SetPayloadSize(sz)

	require.Equal(t, sz, obj.PayloadSize())
}

func TestObject_SetContainerID(t *testing.T) {
	obj := New()

	cid := cidtest.ID()

	obj.SetContainerID(cid)

	cID, set := obj.ContainerID()
	require.True(t, set)
	require.Equal(t, cid, cID)
}

func TestObject_SetOwnerID(t *testing.T) {
	obj := New()

	ownerID := usertest.ID()

	obj.SetOwnerID(ownerID)

	require.Equal(t, ownerID, obj.OwnerID())
}

func TestObject_SetCreationEpoch(t *testing.T) {
	obj := New()

	creat := uint64(228)
	obj.SetCreationEpoch(creat)

	require.Equal(t, creat, obj.CreationEpoch())
}

func TestObject_SetPayloadChecksum(t *testing.T) {
	obj := New()
	var cs checksum.Checksum
	cs.SetSHA256(randSHA256Checksum(t))

	obj.SetPayloadChecksum(cs)
	cs2, set := obj.PayloadChecksum()

	require.True(t, set)
	require.Equal(t, cs, cs2)
}

func TestObject_SetPayloadHomomorphicHash(t *testing.T) {
	obj := New()

	var cs checksum.Checksum
	cs.SetTillichZemor(randTZChecksum(t))

	obj.SetPayloadHomomorphicHash(cs)
	cs2, set := obj.PayloadHomomorphicHash()

	require.True(t, set)
	require.Equal(t, cs, cs2)
}

func TestObject_SetAttributes(t *testing.T) {
	obj := New()

	a1 := NewAttribute()
	a1.SetKey("key1")
	a1.SetValue("val1")

	a2 := NewAttribute()
	a2.SetKey("key2")
	a2.SetValue("val2")

	obj.SetAttributes(*a1, *a2)

	require.Equal(t, []Attribute{*a1, *a2}, obj.Attributes())
}

func TestObject_SetPreviousID(t *testing.T) {
	obj := New()

	prev := randID(t)

	obj.SetPreviousID(prev)

	oID, set := obj.PreviousID()

	require.True(t, set)
	require.Equal(t, prev, oID)
}

func TestObject_SetChildren(t *testing.T) {
	obj := New()

	id1 := randID(t)
	id2 := randID(t)

	obj.SetChildren(id1, id2)

	require.Equal(t, []oid.ID{id1, id2}, obj.Children())
}

func TestObject_SetSplitID(t *testing.T) {
	obj := New()

	require.Nil(t, obj.SplitID())

	splitID := NewSplitID()
	obj.SetSplitID(splitID)

	require.Equal(t, obj.SplitID(), splitID)
}

func TestObject_SetParent(t *testing.T) {
	obj := New()

	require.Nil(t, obj.Parent())

	par := New()
	par.SetID(randID(t))
	par.SetContainerID(cidtest.ID())

	obj.SetParent(par)

	require.Equal(t, par, obj.Parent())
}

func TestObject_ToV2(t *testing.T) {
	objV2 := new(object.Object)
	objV2.SetPayload([]byte{1, 2, 3})

	obj := NewFromV2(objV2)

	require.Equal(t, objV2, obj.ToV2())
}

func TestObject_SetSessionToken(t *testing.T) {
	obj := New()

	tok := sessiontest.ObjectSigned()

	obj.SetSessionToken(tok)

	require.Equal(t, tok, obj.SessionToken())
}

func TestObject_SetType(t *testing.T) {
	obj := New()

	typ := TypeStorageGroup

	obj.SetType(typ)

	require.Equal(t, typ, obj.Type())
}

func TestObject_CutPayload(t *testing.T) {
	o1 := New()

	p1 := []byte{12, 3}
	o1.SetPayload(p1)

	sz := uint64(13)
	o1.SetPayloadSize(sz)

	o2 := o1.CutPayload()

	require.Equal(t, sz, o2.PayloadSize())
	require.Empty(t, o2.Payload())

	sz++
	o1.SetPayloadSize(sz)

	require.Equal(t, sz, o1.PayloadSize())
	require.Equal(t, sz, o2.PayloadSize())

	p2 := []byte{4, 5, 6}
	o2.SetPayload(p2)

	require.Equal(t, p2, o2.Payload())
	require.Equal(t, p1, o1.Payload())
}

func TestObject_SetParentID(t *testing.T) {
	obj := New()

	id := randID(t)
	obj.SetParentID(id)

	oID, set := obj.ParentID()
	require.True(t, set)
	require.Equal(t, id, oID)
}

func TestObject_ResetRelations(t *testing.T) {
	obj := New()

	obj.SetPreviousID(randID(t))

	obj.ResetRelations()

	_, set := obj.PreviousID()
	require.False(t, set)
}

func TestObject_HasParent(t *testing.T) {
	obj := New()

	obj.InitRelations()

	require.True(t, obj.HasParent())

	obj.ResetRelations()

	require.False(t, obj.HasParent())
}

func TestObjectEncoding(t *testing.T) {
	o := New()
	o.SetID(randID(t))
	o.SetContainerID(cidtest.ID())

	t.Run("binary", func(t *testing.T) {
		data, err := o.Marshal()
		require.NoError(t, err)

		o2 := New()
		require.NoError(t, o2.Unmarshal(data))

		require.Equal(t, o, o2)
	})

	t.Run("json", func(t *testing.T) {
		data, err := o.MarshalJSON()
		require.NoError(t, err)

		o2 := New()
		require.NoError(t, o2.UnmarshalJSON(data))

		require.Equal(t, o, o2)
	})
}
