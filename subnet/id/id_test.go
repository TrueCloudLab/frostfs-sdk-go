package subnetid_test

import (
	"testing"

	"github.com/TrueCloudLab/frostfs-api-go/v2/refs"
	subnetid "github.com/TrueCloudLab/frostfs-sdk-go/subnet/id"
	subnetidtest "github.com/TrueCloudLab/frostfs-sdk-go/subnet/id/test"
	"github.com/stretchr/testify/require"
)

func TestIsZero(t *testing.T) {
	var id subnetid.ID

	require.True(t, subnetid.IsZero(id))

	id.SetNumeric(13)
	require.False(t, subnetid.IsZero(id))

	id.SetNumeric(0)
	require.True(t, subnetid.IsZero(id))
}

func TestID_ReadFromV2(t *testing.T) {
	const num = 13

	var id1 subnetid.ID
	id1.SetNumeric(num)

	var idv2 refs.SubnetID
	idv2.SetValue(num)

	var id2 subnetid.ID
	require.NoError(t, id2.ReadFromV2(idv2))

	require.True(t, id1.Equals(id2))
}

func TestID_WriteToV2(t *testing.T) {
	const num = 13

	var (
		id   subnetid.ID
		idv2 refs.SubnetID
	)

	id.WriteToV2(&idv2)
	require.Zero(t, idv2.GetValue())

	id.SetNumeric(num)

	id.WriteToV2(&idv2)
	require.EqualValues(t, num, idv2.GetValue())
}

func TestID_Equals(t *testing.T) {
	const num = 13

	var id1, id2, idOther, id0 subnetid.ID

	id0.Equals(subnetid.ID{})

	id1.SetNumeric(num)
	id2.SetNumeric(num)
	idOther.SetNumeric(num + 1)

	require.True(t, id1.Equals(id2))
	require.False(t, id1.Equals(idOther))
	require.False(t, id2.Equals(idOther))
}

func TestSubnetIDEncoding(t *testing.T) {
	id := subnetidtest.ID()

	t.Run("binary", func(t *testing.T) {
		var id2 subnetid.ID
		require.NoError(t, id2.Unmarshal(id.Marshal()))

		require.True(t, id2.Equals(id))
	})

	t.Run("text", func(t *testing.T) {
		var id2 subnetid.ID
		require.NoError(t, id2.DecodeString(id.EncodeToString()))

		require.True(t, id2.Equals(id))
	})
}

func TestMakeZero(t *testing.T) {
	var id subnetid.ID
	id.SetNumeric(13)

	require.False(t, subnetid.IsZero(id))

	subnetid.MakeZero(&id)

	require.True(t, subnetid.IsZero(id))
	require.Equal(t, subnetid.ID{}, id)
}
