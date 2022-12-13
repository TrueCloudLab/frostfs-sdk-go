package container_test

import (
	"crypto/sha256"
	"strconv"
	"testing"
	"time"

	v2container "github.com/TrueCloudLab/frostfs-api-go/v2/container"
	v2netmap "github.com/TrueCloudLab/frostfs-api-go/v2/netmap"
	"github.com/TrueCloudLab/frostfs-api-go/v2/refs"
	"github.com/TrueCloudLab/frostfs-sdk-go/container"
	cid "github.com/TrueCloudLab/frostfs-sdk-go/container/id"
	cidtest "github.com/TrueCloudLab/frostfs-sdk-go/container/id/test"
	containertest "github.com/TrueCloudLab/frostfs-sdk-go/container/test"
	frostfscrypto "github.com/TrueCloudLab/frostfs-sdk-go/crypto"
	netmaptest "github.com/TrueCloudLab/frostfs-sdk-go/netmap/test"
	subnetid "github.com/TrueCloudLab/frostfs-sdk-go/subnet/id"
	subnetidtest "github.com/TrueCloudLab/frostfs-sdk-go/subnet/id/test"
	usertest "github.com/TrueCloudLab/frostfs-sdk-go/user/test"
	"github.com/TrueCloudLab/frostfs-sdk-go/version"
	"github.com/google/uuid"
	"github.com/nspcc-dev/neo-go/pkg/crypto/keys"
	"github.com/stretchr/testify/require"
)

func TestPlacementPolicyEncoding(t *testing.T) {
	v := containertest.Container()

	t.Run("binary", func(t *testing.T) {
		var v2 container.Container
		require.NoError(t, v2.Unmarshal(v.Marshal()))

		require.Equal(t, v, v2)
	})

	t.Run("json", func(t *testing.T) {
		data, err := v.MarshalJSON()
		require.NoError(t, err)

		var v2 container.Container
		require.NoError(t, v2.UnmarshalJSON(data))

		require.Equal(t, v, v2)
	})
}

func TestContainer_Init(t *testing.T) {
	val := containertest.Container()

	val.Init()

	var msg v2container.Container
	val.WriteToV2(&msg)

	binNonce := msg.GetNonce()

	var nonce uuid.UUID
	require.NoError(t, nonce.UnmarshalBinary(binNonce))
	require.EqualValues(t, 4, nonce.Version())

	verV2 := msg.GetVersion()
	require.NotNil(t, verV2)

	var ver version.Version
	require.NoError(t, ver.ReadFromV2(*verV2))

	require.Equal(t, version.Current(), ver)

	var val2 container.Container
	require.NoError(t, val2.ReadFromV2(msg))

	require.Equal(t, val, val2)
}

func TestContainer_Owner(t *testing.T) {
	var val container.Container

	require.Zero(t, val.Owner())

	val = containertest.Container()

	owner := *usertest.ID()

	val.SetOwner(owner)

	var msg v2container.Container
	val.WriteToV2(&msg)

	var msgOwner refs.OwnerID
	owner.WriteToV2(&msgOwner)

	require.Equal(t, &msgOwner, msg.GetOwnerID())

	var val2 container.Container
	require.NoError(t, val2.ReadFromV2(msg))

	require.True(t, val2.Owner().Equals(owner))
}

func TestContainer_BasicACL(t *testing.T) {
	var val container.Container

	require.Zero(t, val.BasicACL())

	val = containertest.Container()

	basicACL := containertest.BasicACL()
	val.SetBasicACL(basicACL)

	var msg v2container.Container
	val.WriteToV2(&msg)

	require.EqualValues(t, basicACL.Bits(), msg.GetBasicACL())

	var val2 container.Container
	require.NoError(t, val2.ReadFromV2(msg))

	require.Equal(t, basicACL, val2.BasicACL())
}

func TestContainer_PlacementPolicy(t *testing.T) {
	var val container.Container

	require.Zero(t, val.PlacementPolicy())

	val = containertest.Container()

	pp := netmaptest.PlacementPolicy()
	val.SetPlacementPolicy(pp)

	var msgPolicy v2netmap.PlacementPolicy
	pp.WriteToV2(&msgPolicy)

	var msg v2container.Container
	val.WriteToV2(&msg)

	require.Equal(t, &msgPolicy, msg.GetPlacementPolicy())

	var val2 container.Container
	require.NoError(t, val2.ReadFromV2(msg))

	require.Equal(t, pp, val2.PlacementPolicy())
}

func assertContainsAttribute(t *testing.T, m v2container.Container, key, val string) {
	var msgAttr v2container.Attribute

	msgAttr.SetKey(key)
	msgAttr.SetValue(val)
	require.Contains(t, m.GetAttributes(), msgAttr)
}

func TestContainer_Attribute(t *testing.T) {
	const attrKey1, attrKey2 = "key1", "key2"
	const attrVal1, attrVal2 = "val1", "val2"

	val := containertest.Container()

	val.SetAttribute(attrKey1, attrVal1)
	val.SetAttribute(attrKey2, attrVal2)

	var msg v2container.Container
	val.WriteToV2(&msg)

	require.GreaterOrEqual(t, len(msg.GetAttributes()), 2)
	assertContainsAttribute(t, msg, attrKey1, attrVal1)
	assertContainsAttribute(t, msg, attrKey2, attrVal2)

	var val2 container.Container
	require.NoError(t, val2.ReadFromV2(msg))

	require.Equal(t, attrVal1, val2.Attribute(attrKey1))
	require.Equal(t, attrVal2, val2.Attribute(attrKey2))

	m := map[string]string{}

	val2.IterateAttributes(func(key, val string) {
		m[key] = val
	})

	require.GreaterOrEqual(t, len(m), 2)
	require.Equal(t, attrVal1, m[attrKey1])
	require.Equal(t, attrVal2, m[attrKey2])

	val2.SetAttribute(attrKey1, attrVal1+"_")
	require.Equal(t, attrVal1+"_", val2.Attribute(attrKey1))
}

func TestSetName(t *testing.T) {
	var val container.Container

	require.Panics(t, func() {
		container.SetName(&val, "")
	})

	val = containertest.Container()

	const name = "some name"

	container.SetName(&val, name)

	var msg v2container.Container
	val.WriteToV2(&msg)

	assertContainsAttribute(t, msg, "Name", name)

	var val2 container.Container
	require.NoError(t, val2.ReadFromV2(msg))

	require.Equal(t, name, container.Name(val2))
}

func TestSetCreationTime(t *testing.T) {
	var val container.Container

	require.Zero(t, container.CreatedAt(val).Unix())

	val = containertest.Container()

	creat := time.Now()

	container.SetCreationTime(&val, creat)

	var msg v2container.Container
	val.WriteToV2(&msg)

	assertContainsAttribute(t, msg, "Timestamp", strconv.FormatInt(creat.Unix(), 10))

	var val2 container.Container
	require.NoError(t, val2.ReadFromV2(msg))

	require.Equal(t, creat.Unix(), container.CreatedAt(val2).Unix())
}

func TestSetSubnet(t *testing.T) {
	var val container.Container

	require.True(t, subnetid.IsZero(container.Subnet(val)))

	val = containertest.Container()

	sub := subnetidtest.ID()

	container.SetSubnet(&val, sub)

	var msg v2container.Container
	val.WriteToV2(&msg)

	assertContainsAttribute(t, msg, v2container.SysAttributeSubnet, sub.EncodeToString())

	var val2 container.Container
	require.NoError(t, val2.ReadFromV2(msg))

	require.Equal(t, sub, container.Subnet(val))
}

func TestDisableHomomorphicHashing(t *testing.T) {
	var val container.Container

	require.False(t, container.IsHomomorphicHashingDisabled(val))

	val = containertest.Container()

	container.DisableHomomorphicHashing(&val)

	var msg v2container.Container
	val.WriteToV2(&msg)

	assertContainsAttribute(t, msg, v2container.SysAttributePrefix+"DISABLE_HOMOMORPHIC_HASHING", "true")

	var val2 container.Container
	require.NoError(t, val2.ReadFromV2(msg))

	require.True(t, container.IsHomomorphicHashingDisabled(val2))
}

func TestWriteDomain(t *testing.T) {
	var val container.Container

	require.Zero(t, container.ReadDomain(val).Name())

	val = containertest.Container()

	const name = "domain name"

	var d container.Domain
	d.SetName(name)

	container.WriteDomain(&val, d)

	var msg v2container.Container
	val.WriteToV2(&msg)

	assertContainsAttribute(t, msg, v2container.SysAttributeName, name)
	assertContainsAttribute(t, msg, v2container.SysAttributeZone, "container")

	const zone = "domain zone"

	d.SetZone(zone)

	container.WriteDomain(&val, d)

	val.WriteToV2(&msg)

	assertContainsAttribute(t, msg, v2container.SysAttributeZone, zone)

	var val2 container.Container
	require.NoError(t, val2.ReadFromV2(msg))

	require.Equal(t, d, container.ReadDomain(val2))
}

func TestCalculateID(t *testing.T) {
	val := containertest.Container()

	require.False(t, container.AssertID(cidtest.ID(), val))

	var id cid.ID
	container.CalculateID(&id, val)

	var msg refs.ContainerID
	id.WriteToV2(&msg)

	h := sha256.Sum256(val.Marshal())
	require.Equal(t, h[:], msg.GetValue())

	var id2 cid.ID
	require.NoError(t, id2.ReadFromV2(msg))

	require.True(t, container.AssertID(id2, val))
}

func TestCalculateSignature(t *testing.T) {
	key, err := keys.NewPrivateKey()
	require.NoError(t, err)

	val := containertest.Container()

	var sig frostfscrypto.Signature

	require.NoError(t, container.CalculateSignature(&sig, val, key.PrivateKey))

	var msg refs.Signature
	sig.WriteToV2(&msg)

	var sig2 frostfscrypto.Signature
	require.NoError(t, sig2.ReadFromV2(msg))

	require.True(t, container.VerifySignature(sig2, val))
}
