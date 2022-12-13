package user_test

import (
	"math/rand"
	"testing"

	"github.com/TrueCloudLab/frostfs-api-go/v2/refs"
	. "github.com/TrueCloudLab/frostfs-sdk-go/user"
	usertest "github.com/TrueCloudLab/frostfs-sdk-go/user/test"
	"github.com/mr-tron/base58"
	"github.com/nspcc-dev/neo-go/pkg/util"
	"github.com/nspcc-dev/neo-go/pkg/util/slice"
	"github.com/stretchr/testify/require"
)

func TestID_WalletBytes(t *testing.T) {
	var scriptHash util.Uint160
	rand.Read(scriptHash[:])

	var id ID
	id.SetScriptHash(scriptHash)

	w := id.WalletBytes()

	var m refs.OwnerID
	m.SetValue(w)

	err := id.ReadFromV2(m)
	require.NoError(t, err)
}

func TestID_SetScriptHash(t *testing.T) {
	var scriptHash util.Uint160
	rand.Read(scriptHash[:])

	var id ID
	id.SetScriptHash(scriptHash)

	var m refs.OwnerID
	id.WriteToV2(&m)

	var id2 ID

	err := id2.ReadFromV2(m)
	require.NoError(t, err)

	require.True(t, id2.Equals(id))
}

func TestV2_ID(t *testing.T) {
	id := *usertest.ID()
	var m refs.OwnerID
	var id2 ID

	t.Run("OK", func(t *testing.T) {
		id.WriteToV2(&m)

		err := id2.ReadFromV2(m)
		require.NoError(t, err)
		require.True(t, id2.Equals(id))
	})

	val := m.GetValue()

	t.Run("invalid size", func(t *testing.T) {
		m.SetValue(val[:24])

		err := id2.ReadFromV2(m)
		require.Error(t, err)
	})

	t.Run("invalid prefix", func(t *testing.T) {
		val := slice.Copy(val)
		val[0]++

		m.SetValue(val)

		err := id2.ReadFromV2(m)
		require.Error(t, err)
	})

	t.Run("invalid checksum", func(t *testing.T) {
		val := slice.Copy(val)
		val[21]++

		m.SetValue(val)

		err := id2.ReadFromV2(m)
		require.Error(t, err)
	})
}

func TestID_EncodeToString(t *testing.T) {
	id := *usertest.ID()

	s := id.EncodeToString()

	_, err := base58.Decode(s)
	require.NoError(t, err)

	var id2 ID

	err = id2.DecodeString(s)
	require.NoError(t, err)

	require.Equal(t, id, id2)

	err = id2.DecodeString("_") // any invalid bas58 string
	require.Error(t, err)
}

func TestID_Equal(t *testing.T) {
	id1 := *usertest.ID()
	id2 := *usertest.ID()
	id3 := id1

	require.True(t, id1.Equals(id1)) // self-equality
	require.True(t, id1.Equals(id3))
	require.True(t, id3.Equals(id1)) // commutativity
	require.False(t, id1.Equals(id2))
}
