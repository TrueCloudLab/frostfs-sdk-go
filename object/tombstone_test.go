package object

import (
	"crypto/sha256"
	"math/rand"
	"testing"

	"github.com/TrueCloudLab/frostfs-api-go/v2/tombstone"
	oid "github.com/TrueCloudLab/frostfs-sdk-go/object/id"
	"github.com/stretchr/testify/require"
)

func generateIDList(sz int) []oid.ID {
	res := make([]oid.ID, sz)
	cs := [sha256.Size]byte{}

	for i := 0; i < sz; i++ {
		var oID oid.ID

		res[i] = oID
		rand.Read(cs[:])
		res[i].SetSHA256(cs)
	}

	return res
}

func TestTombstone(t *testing.T) {
	ts := NewTombstone()

	exp := uint64(13)
	ts.SetExpirationEpoch(exp)
	require.Equal(t, exp, ts.ExpirationEpoch())

	splitID := NewSplitID()
	ts.SetSplitID(splitID)
	require.Equal(t, splitID, ts.SplitID())

	members := generateIDList(3)
	ts.SetMembers(members)
	require.Equal(t, members, ts.Members())
}

func TestTombstoneEncoding(t *testing.T) {
	ts := NewTombstone()
	ts.SetExpirationEpoch(13)
	ts.SetSplitID(NewSplitID())
	ts.SetMembers(generateIDList(5))

	t.Run("binary", func(t *testing.T) {
		data, err := ts.Marshal()
		require.NoError(t, err)

		ts2 := NewTombstone()
		require.NoError(t, ts2.Unmarshal(data))

		require.Equal(t, ts, ts2)
	})

	t.Run("json", func(t *testing.T) {
		data, err := ts.MarshalJSON()
		require.NoError(t, err)

		ts2 := NewTombstone()
		require.NoError(t, ts2.UnmarshalJSON(data))

		require.Equal(t, ts, ts2)
	})
}

func TestNewTombstoneFromV2(t *testing.T) {
	t.Run("from nil", func(t *testing.T) {
		var x *tombstone.Tombstone

		require.Nil(t, NewTombstoneFromV2(x))
	})
}

func TestNewTombstone(t *testing.T) {
	t.Run("default values", func(t *testing.T) {
		ts := NewTombstone()

		// check initial values
		require.Nil(t, ts.SplitID())
		require.Nil(t, ts.Members())
		require.Zero(t, ts.ExpirationEpoch())

		// convert to v2 message
		tsV2 := ts.ToV2()

		require.Nil(t, tsV2.GetSplitID())
		require.Nil(t, tsV2.GetMembers())
		require.Zero(t, tsV2.GetExpirationEpoch())
	})
}
