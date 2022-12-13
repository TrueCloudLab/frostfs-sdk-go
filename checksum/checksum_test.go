package checksum

import (
	"crypto/rand"
	"crypto/sha256"
	"testing"

	"github.com/TrueCloudLab/frostfs-api-go/v2/refs"
	"github.com/TrueCloudLab/tzhash/tz"
	"github.com/stretchr/testify/require"
)

func TestChecksum(t *testing.T) {
	var c Checksum

	cSHA256 := [sha256.Size]byte{}
	_, _ = rand.Read(cSHA256[:])

	c.SetSHA256(cSHA256)

	require.Equal(t, SHA256, c.Type())
	require.Equal(t, cSHA256[:], c.Value())

	var cV2 refs.Checksum
	c.WriteToV2(&cV2)

	require.Equal(t, refs.SHA256, cV2.GetType())
	require.Equal(t, cSHA256[:], cV2.GetSum())

	cTZ := [tz.Size]byte{}
	_, _ = rand.Read(cSHA256[:])

	c.SetTillichZemor(cTZ)

	require.Equal(t, TZ, c.Type())
	require.Equal(t, cTZ[:], c.Value())

	c.WriteToV2(&cV2)

	require.Equal(t, refs.TillichZemor, cV2.GetType())
	require.Equal(t, cTZ[:], cV2.GetSum())
}

func TestNewChecksum(t *testing.T) {
	t.Run("default values", func(t *testing.T) {
		var chs Checksum

		// check initial values
		require.Equal(t, Unknown, chs.Type())
		require.Nil(t, chs.Value())

		// convert to v2 message
		var chsV2 refs.Checksum
		chs.WriteToV2(&chsV2)

		require.Equal(t, refs.UnknownChecksum, chsV2.GetType())
		require.Nil(t, chsV2.GetSum())
	})
}

func TestCalculation(t *testing.T) {
	var c Checksum
	payload := []byte{0, 1, 2, 3, 4, 5}

	t.Run("SHA256", func(t *testing.T) {
		orig := sha256.Sum256(payload)

		Calculate(&c, SHA256, payload)

		require.Equal(t, orig[:], c.Value())
	})

	t.Run("TZ", func(t *testing.T) {
		orig := tz.Sum(payload)

		Calculate(&c, TZ, payload)

		require.Equal(t, orig[:], c.Value())
	})
}
