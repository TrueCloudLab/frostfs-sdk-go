package storagegroup_test

import (
	"crypto/sha256"
	"strconv"
	"testing"

	objectV2 "github.com/TrueCloudLab/frostfs-api-go/v2/object"
	"github.com/TrueCloudLab/frostfs-api-go/v2/refs"
	storagegroupV2 "github.com/TrueCloudLab/frostfs-api-go/v2/storagegroup"
	storagegroupV2test "github.com/TrueCloudLab/frostfs-api-go/v2/storagegroup/test"
	"github.com/TrueCloudLab/frostfs-sdk-go/checksum"
	checksumtest "github.com/TrueCloudLab/frostfs-sdk-go/checksum/test"
	objectSDK "github.com/TrueCloudLab/frostfs-sdk-go/object"
	oid "github.com/TrueCloudLab/frostfs-sdk-go/object/id"
	oidtest "github.com/TrueCloudLab/frostfs-sdk-go/object/id/test"
	"github.com/TrueCloudLab/frostfs-sdk-go/storagegroup"
	storagegrouptest "github.com/TrueCloudLab/frostfs-sdk-go/storagegroup/test"
	"github.com/stretchr/testify/require"
)

func TestStorageGroup(t *testing.T) {
	var sg storagegroup.StorageGroup

	sz := uint64(13)
	sg.SetValidationDataSize(sz)
	require.Equal(t, sz, sg.ValidationDataSize())

	cs := checksumtest.Checksum()
	sg.SetValidationDataHash(cs)
	cs2, set := sg.ValidationDataHash()

	require.True(t, set)
	require.Equal(t, cs, cs2)

	exp := uint64(33)
	sg.SetExpirationEpoch(exp)
	require.Equal(t, exp, sg.ExpirationEpoch())

	members := []oid.ID{oidtest.ID(), oidtest.ID()}
	sg.SetMembers(members)
	require.Equal(t, members, sg.Members())
}

func TestStorageGroup_ReadFromV2(t *testing.T) {
	t.Run("from zero", func(t *testing.T) {
		var (
			x  storagegroup.StorageGroup
			v2 storagegroupV2.StorageGroup
		)

		require.Error(t, x.ReadFromV2(v2))
	})

	t.Run("from non-zero", func(t *testing.T) {
		var (
			x  storagegroup.StorageGroup
			v2 = storagegroupV2test.GenerateStorageGroup(false)
		)

		// https://github.com/nspcc-dev/neofs-api-go/issues/394
		v2.SetMembers(generateOIDList())

		size := v2.GetValidationDataSize()
		epoch := v2.GetExpirationEpoch()
		mm := v2.GetMembers()
		hashV2 := v2.GetValidationHash()

		require.NoError(t, x.ReadFromV2(*v2))

		require.Equal(t, epoch, x.ExpirationEpoch())
		require.Equal(t, size, x.ValidationDataSize())

		var hash checksum.Checksum
		require.NoError(t, hash.ReadFromV2(*hashV2))
		h, set := x.ValidationDataHash()
		require.True(t, set)
		require.Equal(t, hash, h)

		var oidV2 refs.ObjectID

		for i, m := range mm {
			x.Members()[i].WriteToV2(&oidV2)
			require.Equal(t, m, oidV2)
		}
	})
}

func TestStorageGroupEncoding(t *testing.T) {
	sg := storagegrouptest.StorageGroup()

	t.Run("binary", func(t *testing.T) {
		data, err := sg.Marshal()
		require.NoError(t, err)

		var sg2 storagegroup.StorageGroup
		require.NoError(t, sg2.Unmarshal(data))

		require.Equal(t, sg, sg2)
	})

	t.Run("json", func(t *testing.T) {
		data, err := sg.MarshalJSON()
		require.NoError(t, err)

		var sg2 storagegroup.StorageGroup
		require.NoError(t, sg2.UnmarshalJSON(data))

		require.Equal(t, sg, sg2)
	})
}

func TestStorageGroup_WriteToV2(t *testing.T) {
	t.Run("zero to v2", func(t *testing.T) {
		var (
			x  storagegroup.StorageGroup
			v2 storagegroupV2.StorageGroup
		)

		x.WriteToV2(&v2)

		require.Nil(t, v2.GetValidationHash())
		require.Nil(t, v2.GetMembers())
		require.Zero(t, v2.GetValidationDataSize())
		require.Zero(t, v2.GetExpirationEpoch())
	})

	t.Run("non-zero to v2", func(t *testing.T) {
		var (
			x  = storagegrouptest.StorageGroup()
			v2 storagegroupV2.StorageGroup
		)

		x.WriteToV2(&v2)

		require.Equal(t, x.ExpirationEpoch(), v2.GetExpirationEpoch())
		require.Equal(t, x.ValidationDataSize(), v2.GetValidationDataSize())

		var hash checksum.Checksum
		require.NoError(t, hash.ReadFromV2(*v2.GetValidationHash()))

		h, set := x.ValidationDataHash()
		require.True(t, set)
		require.Equal(t, h, hash)

		var oidV2 refs.ObjectID

		for i, m := range x.Members() {
			m.WriteToV2(&oidV2)
			require.Equal(t, oidV2, v2.GetMembers()[i])
		}
	})
}

func TestNew(t *testing.T) {
	t.Run("default values", func(t *testing.T) {
		var sg storagegroup.StorageGroup

		// check initial values
		require.Nil(t, sg.Members())
		_, set := sg.ValidationDataHash()
		require.False(t, set)
		require.Zero(t, sg.ExpirationEpoch())
		require.Zero(t, sg.ValidationDataSize())
	})
}

func generateOIDList() []refs.ObjectID {
	const size = 3

	mmV2 := make([]refs.ObjectID, size)
	for i := 0; i < size; i++ {
		oidV2 := make([]byte, sha256.Size)
		oidV2[i] = byte(i)

		mmV2[i].SetValue(oidV2)
	}

	return mmV2
}

func TestStorageGroup_SetMembers_DoubleSetting(t *testing.T) {
	var sg storagegroup.StorageGroup

	mm := []oid.ID{oidtest.ID(), oidtest.ID(), oidtest.ID()} // cap is 3 at least
	sg.SetMembers(mm)

	// the previous cap is more that a new length;
	// slicing should not lead to `out of range`
	// and apply update correctly
	sg.SetMembers(mm[:1])
}

func TestStorageGroupFromObject(t *testing.T) {
	sg := storagegrouptest.StorageGroup()

	var o objectSDK.Object

	var expAttr objectSDK.Attribute
	expAttr.SetKey(objectV2.SysAttributeExpEpoch)
	expAttr.SetValue(strconv.FormatUint(sg.ExpirationEpoch(), 10))

	sgRaw, err := sg.Marshal()
	require.NoError(t, err)

	o.SetPayload(sgRaw)
	o.SetType(objectSDK.TypeStorageGroup)

	t.Run("correct object", func(t *testing.T) {
		o.SetAttributes(objectSDK.Attribute{}, expAttr, objectSDK.Attribute{})

		var sg2 storagegroup.StorageGroup
		require.NoError(t, storagegroup.ReadFromObject(&sg2, o))
		require.Equal(t, sg, sg2)
	})

	t.Run("incorrect exp attr", func(t *testing.T) {
		var sg2 storagegroup.StorageGroup

		expAttr.SetValue(strconv.FormatUint(sg.ExpirationEpoch()+1, 10))
		o.SetAttributes(expAttr)

		require.Error(t, storagegroup.ReadFromObject(&sg2, o))
	})

	t.Run("incorrect object type", func(t *testing.T) {
		var sg2 storagegroup.StorageGroup

		o.SetType(objectSDK.TypeTombstone)
		require.Error(t, storagegroup.ReadFromObject(&sg2, o))
	})
}

func TestStorageGroupToObject(t *testing.T) {
	sg := storagegrouptest.StorageGroup()

	sgRaw, err := sg.Marshal()
	require.NoError(t, err)

	t.Run("empty object", func(t *testing.T) {
		var o objectSDK.Object
		storagegroup.WriteToObject(sg, &o)

		exp, found := expFromObj(t, o)
		require.True(t, found)

		require.Equal(t, sgRaw, o.Payload())
		require.Equal(t, sg.ExpirationEpoch(), exp)
		require.Equal(t, objectSDK.TypeStorageGroup, o.Type())
	})

	t.Run("obj already has exp attr", func(t *testing.T) {
		var o objectSDK.Object

		var attr objectSDK.Attribute
		attr.SetKey(objectV2.SysAttributeExpEpoch)
		attr.SetValue(strconv.FormatUint(sg.ExpirationEpoch()+1, 10))

		o.SetAttributes(objectSDK.Attribute{}, attr, objectSDK.Attribute{})

		storagegroup.WriteToObject(sg, &o)

		exp, found := expFromObj(t, o)
		require.True(t, found)

		require.Equal(t, sgRaw, o.Payload())
		require.Equal(t, sg.ExpirationEpoch(), exp)
		require.Equal(t, objectSDK.TypeStorageGroup, o.Type())
	})
}

func expFromObj(t *testing.T, o objectSDK.Object) (uint64, bool) {
	for _, attr := range o.Attributes() {
		if attr.Key() == objectV2.SysAttributeExpEpoch {
			exp, err := strconv.ParseUint(attr.Value(), 10, 64)
			require.NoError(t, err)

			return exp, true
		}
	}

	return 0, false
}
