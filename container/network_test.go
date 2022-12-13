package container_test

import (
	"testing"

	"github.com/TrueCloudLab/frostfs-sdk-go/container"
	containertest "github.com/TrueCloudLab/frostfs-sdk-go/container/test"
	netmaptest "github.com/TrueCloudLab/frostfs-sdk-go/netmap/test"
	"github.com/stretchr/testify/require"
)

func TestContainer_NetworkConfig(t *testing.T) {
	c := containertest.Container()
	nc := netmaptest.NetworkInfo()

	t.Run("default", func(t *testing.T) {
		require.False(t, container.IsHomomorphicHashingDisabled(c))

		res := container.AssertNetworkConfig(c, nc)

		require.True(t, res)
	})

	nc.DisableHomomorphicHashing()

	t.Run("apply", func(t *testing.T) {
		require.False(t, container.IsHomomorphicHashingDisabled(c))

		container.ApplyNetworkConfig(&c, nc)

		require.True(t, container.IsHomomorphicHashingDisabled(c))
	})
}
