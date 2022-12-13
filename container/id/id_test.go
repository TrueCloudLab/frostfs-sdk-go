package cid_test

import (
	"crypto/sha256"
	"math/rand"
	"testing"

	"github.com/TrueCloudLab/frostfs-api-go/v2/refs"
	cid "github.com/TrueCloudLab/frostfs-sdk-go/container/id"
	cidtest "github.com/TrueCloudLab/frostfs-sdk-go/container/id/test"
	"github.com/mr-tron/base58"
	"github.com/stretchr/testify/require"
)

func randSHA256Checksum() (cs [sha256.Size]byte) {
	rand.Read(cs[:])
	return
}

const emptyID = "11111111111111111111111111111111"

func TestID_ToV2(t *testing.T) {
	t.Run("non-zero", func(t *testing.T) {
		checksum := randSHA256Checksum()

		id := cidtest.IDWithChecksum(checksum)

		var idV2 refs.ContainerID
		id.WriteToV2(&idV2)

		var newID cid.ID
		require.NoError(t, newID.ReadFromV2(idV2))

		require.Equal(t, id, newID)
		require.Equal(t, checksum[:], idV2.GetValue())
	})

	t.Run("zero", func(t *testing.T) {
		var (
			x  cid.ID
			v2 refs.ContainerID
		)

		x.WriteToV2(&v2)
		require.Equal(t, emptyID, base58.Encode(v2.GetValue()))
	})
}

func TestID_Equal(t *testing.T) {
	cs := randSHA256Checksum()

	id1 := cidtest.IDWithChecksum(cs)
	id2 := cidtest.IDWithChecksum(cs)

	require.True(t, id1.Equals(id2))

	id3 := cidtest.ID()

	require.False(t, id1.Equals(id3))
}

func TestID_String(t *testing.T) {
	t.Run("DecodeString/EncodeToString", func(t *testing.T) {
		id := cidtest.ID()
		var id2 cid.ID

		require.NoError(t, id2.DecodeString(id.EncodeToString()))
		require.Equal(t, id, id2)
	})

	t.Run("zero", func(t *testing.T) {
		var id cid.ID

		require.Equal(t, emptyID, id.EncodeToString())
	})
}

func TestNewFromV2(t *testing.T) {
	t.Run("from zero", func(t *testing.T) {
		var (
			x  cid.ID
			v2 refs.ContainerID
		)

		require.Error(t, x.ReadFromV2(v2))
	})
}

func TestID_Encode(t *testing.T) {
	var id cid.ID

	t.Run("panic", func(t *testing.T) {
		dst := make([]byte, sha256.Size-1)

		require.Panics(t, func() {
			id.Encode(dst)
		})
	})

	t.Run("correct", func(t *testing.T) {
		dst := make([]byte, sha256.Size)

		require.NotPanics(t, func() {
			id.Encode(dst)
		})
		require.Equal(t, emptyID, id.EncodeToString())
	})
}
