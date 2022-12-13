package oid_test

import (
	"testing"

	"github.com/TrueCloudLab/frostfs-api-go/v2/refs"
	cidtest "github.com/TrueCloudLab/frostfs-sdk-go/container/id/test"
	oid "github.com/TrueCloudLab/frostfs-sdk-go/object/id"
	oidtest "github.com/TrueCloudLab/frostfs-sdk-go/object/id/test"
	"github.com/stretchr/testify/require"
)

func TestAddress_SetContainer(t *testing.T) {
	var x oid.Address

	require.Zero(t, x.Container())

	cnr := cidtest.ID()

	x.SetContainer(cnr)
	require.Equal(t, cnr, x.Container())
}

func TestAddress_SetObject(t *testing.T) {
	var x oid.Address

	require.Zero(t, x.Object())

	obj := oidtest.ID()

	x.SetObject(obj)
	require.Equal(t, obj, x.Object())
}

func TestAddress_ReadFromV2(t *testing.T) {
	var x oid.Address
	var xV2 refs.Address

	require.Error(t, x.ReadFromV2(xV2))

	var cnrV2 refs.ContainerID
	xV2.SetContainerID(&cnrV2)

	require.Error(t, x.ReadFromV2(xV2))

	cnr := cidtest.ID()
	cnr.WriteToV2(&cnrV2)

	require.Error(t, x.ReadFromV2(xV2))

	var objV2 refs.ObjectID
	xV2.SetObjectID(&objV2)

	require.Error(t, x.ReadFromV2(xV2))

	obj := oidtest.ID()
	obj.WriteToV2(&objV2)

	require.NoError(t, x.ReadFromV2(xV2))
	require.Equal(t, cnr, x.Container())
	require.Equal(t, obj, x.Object())

	var xV2To refs.Address
	x.WriteToV2(&xV2To)

	require.Equal(t, xV2, xV2To)
}

func TestAddress_DecodeString(t *testing.T) {
	var x, x2 oid.Address

	require.NoError(t, x2.DecodeString(x.EncodeToString()))
	require.Equal(t, x, x2)

	cnr := cidtest.ID()
	obj := oidtest.ID()

	x.SetContainer(cnr)
	x.SetObject(obj)

	require.NoError(t, x2.DecodeString(x.EncodeToString()))
	require.Equal(t, x, x2)

	strCnr := cnr.EncodeToString()
	strObj := obj.EncodeToString()

	require.Error(t, x2.DecodeString(""))
	require.Error(t, x2.DecodeString("/"))
	require.Error(t, x2.DecodeString(strCnr))
	require.Error(t, x2.DecodeString(strCnr+"/"))
	require.Error(t, x2.DecodeString("/"+strCnr))
	require.Error(t, x2.DecodeString(strCnr+strObj))
	require.Error(t, x2.DecodeString(strCnr+"\\"+strObj))
	require.NoError(t, x2.DecodeString(strCnr+"/"+strObj))
}

func TestAddressEncoding(t *testing.T) {
	v := oidtest.Address()

	t.Run("json", func(t *testing.T) {
		data, err := v.MarshalJSON()
		require.NoError(t, err)

		var v2 oid.Address
		require.NoError(t, v2.UnmarshalJSON(data))

		require.Equal(t, v, v2)
	})
}
