package version

import (
	"testing"

	"github.com/TrueCloudLab/frostfs-api-go/v2/refs"
	"github.com/stretchr/testify/require"
)

func TestNewVersion(t *testing.T) {
	t.Run("default values", func(t *testing.T) {
		var v Version

		// check initial values
		require.Zero(t, v.Major())
		require.Zero(t, v.Minor())

		// convert to v2 message
		var vV2 refs.Version
		v.WriteToV2(&vV2)

		require.Zero(t, vV2.GetMajor())
		require.Zero(t, vV2.GetMinor())
	})

	t.Run("setting values", func(t *testing.T) {
		var v Version

		var mjr, mnr uint32 = 1, 2

		v.SetMajor(mjr)
		v.SetMinor(mnr)
		require.Equal(t, mjr, v.Major())
		require.Equal(t, mnr, v.Minor())

		// convert to v2 message
		var ver refs.Version
		v.WriteToV2(&ver)

		require.Equal(t, mjr, ver.GetMajor())
		require.Equal(t, mnr, ver.GetMinor())
	})
}

func TestSDKVersion(t *testing.T) {
	v := Current()

	require.Equal(t, uint32(sdkMjr), v.Major())
	require.Equal(t, uint32(sdkMnr), v.Minor())
}
