package netmap_test

import (
	"testing"

	v2netmap "github.com/TrueCloudLab/frostfs-api-go/v2/netmap"
	"github.com/TrueCloudLab/frostfs-sdk-go/netmap"
	netmaptest "github.com/TrueCloudLab/frostfs-sdk-go/netmap/test"
	"github.com/stretchr/testify/require"
)

func TestNetMapNodes(t *testing.T) {
	var nm netmap.NetMap

	require.Empty(t, nm.Nodes())

	nodes := []netmap.NodeInfo{netmaptest.NodeInfo(), netmaptest.NodeInfo()}

	nm.SetNodes(nodes)
	require.ElementsMatch(t, nodes, nm.Nodes())

	nodesV2 := make([]v2netmap.NodeInfo, len(nodes))
	for i := range nodes {
		nodes[i].WriteToV2(&nodesV2[i])
	}

	var m v2netmap.NetMap
	nm.WriteToV2(&m)

	require.ElementsMatch(t, nodesV2, m.Nodes())
}

func TestNetMap_SetEpoch(t *testing.T) {
	var nm netmap.NetMap

	require.Zero(t, nm.Epoch())

	const e = 158

	nm.SetEpoch(e)
	require.EqualValues(t, e, nm.Epoch())

	var m v2netmap.NetMap
	nm.WriteToV2(&m)

	require.EqualValues(t, e, m.Epoch())
}
