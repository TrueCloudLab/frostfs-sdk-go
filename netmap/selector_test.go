package netmap

import (
	"fmt"
	"math/rand"
	"sort"
	"strconv"
	"testing"

	"github.com/TrueCloudLab/frostfs-api-go/v2/netmap"
	"github.com/TrueCloudLab/hrw"
	"github.com/stretchr/testify/require"
)

func BenchmarkHRWSort(b *testing.B) {
	const netmapSize = 1000

	vectors := make([]nodes, netmapSize)
	weights := make([]float64, netmapSize)
	for i := range vectors {
		key := make([]byte, 33)
		rand.Read(key)

		var node NodeInfo
		node.SetPrice(1)
		node.SetCapacity(100)
		node.SetPublicKey(key)

		vectors[i] = nodes{node}
		weights[i] = float64(rand.Uint32()%10) / 10.0
	}

	pivot := rand.Uint64()
	b.Run("sort by index, no weight", func(b *testing.B) {
		realNodes := make([]nodes, netmapSize)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			b.StopTimer()
			copy(realNodes, vectors)
			b.StartTimer()

			hrw.SortSliceByIndex(realNodes, pivot)
		}
	})
	b.Run("sort by value, no weight", func(b *testing.B) {
		realNodes := make([]nodes, netmapSize)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			b.StopTimer()
			copy(realNodes, vectors)
			b.StartTimer()

			hrw.SortSliceByValue(realNodes, pivot)
		}
	})
	b.Run("only sort by index", func(b *testing.B) {
		realNodes := make([]nodes, netmapSize)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			b.StopTimer()
			copy(realNodes, vectors)
			b.StartTimer()

			hrw.SortSliceByWeightIndex(realNodes, weights, pivot)
		}
	})
	b.Run("sort by value", func(b *testing.B) {
		realNodes := make([]nodes, netmapSize)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			b.StopTimer()
			copy(realNodes, vectors)
			b.StartTimer()

			hrw.SortSliceByWeightValue(realNodes, weights, pivot)
		}
	})
	b.Run("sort by ID, then by index (deterministic)", func(b *testing.B) {
		realNodes := make([]nodes, netmapSize)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			b.StopTimer()
			copy(realNodes, vectors)
			b.StartTimer()

			sort.Slice(vectors, func(i, j int) bool {
				return less(vectors[i][0], vectors[j][0])
			})
			hrw.SortSliceByWeightIndex(realNodes, weights, pivot)
		}
	})
}

func BenchmarkPolicyHRWType(b *testing.B) {
	const netmapSize = 100

	p := newPlacementPolicy(1,
		[]ReplicaDescriptor{
			newReplica(1, "loc1"),
			newReplica(1, "loc2")},
		[]Selector{
			newSelector("loc1", "Location", 1, "loc1", (*Selector).SelectSame),
			newSelector("loc2", "Location", 1, "loc2", (*Selector).SelectSame)},
		[]Filter{
			newFilter("loc1", "Location", "Shanghai", netmap.EQ),
			newFilter("loc2", "Location", "Shanghai", netmap.NE),
		})

	nodes := make([]NodeInfo, netmapSize)
	for i := range nodes {
		var loc string
		switch i % 20 {
		case 0:
			loc = "Shanghai"
		default:
			loc = strconv.Itoa(i % 20)
		}

		// Having the same price and capacity ensures equal weights for all nodes.
		// This way placement is more dependent on the initial order.
		nodes[i] = nodeInfoFromAttributes("Location", loc, "Price", "1", "Capacity", "10")
		pub := make([]byte, 33)
		pub[0] = byte(i)
		nodes[i].SetPublicKey(pub)
	}

	var nm NetMap
	nm.SetNodes(nodes)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := nm.ContainerNodes(p, []byte{1})
		if err != nil {
			b.Fatal()
		}
	}
}

func TestPlacementPolicy_DeterministicOrder(t *testing.T) {
	const netmapSize = 100

	p := newPlacementPolicy(1,
		[]ReplicaDescriptor{
			newReplica(1, "loc1"),
			newReplica(1, "loc2")},
		[]Selector{
			newSelector("loc1", "Location", 1, "loc1", (*Selector).SelectSame),
			newSelector("loc2", "Location", 1, "loc2", (*Selector).SelectSame)},
		[]Filter{
			newFilter("loc1", "Location", "Shanghai", netmap.EQ),
			newFilter("loc2", "Location", "Shanghai", netmap.NE),
		})

	nodeList := make([]NodeInfo, netmapSize)
	for i := range nodeList {
		var loc string
		switch i % 20 {
		case 0:
			loc = "Shanghai"
		default:
			loc = strconv.Itoa(i % 20)
		}

		// Having the same price and capacity ensures equal weights for all nodes.
		// This way placement is more dependent on the initial order.
		nodeList[i] = nodeInfoFromAttributes("Location", loc, "Price", "1", "Capacity", "10")
		pub := make([]byte, 33)
		pub[0] = byte(i)
		nodeList[i].SetPublicKey(pub)
	}

	var nm NetMap
	nm.SetNodes(nodeList)

	getIndices := func(t *testing.T) (uint64, uint64) {
		v, err := nm.ContainerNodes(p, []byte{1})
		require.NoError(t, err)

		nss := make([]nodes, len(v))
		for i := range v {
			nss[i] = v[i]
		}

		ns := flattenNodes(nss)
		require.Equal(t, 2, len(ns))
		return ns[0].Hash(), ns[1].Hash()
	}

	a, b := getIndices(t)
	for i := 0; i < 10; i++ {
		x, y := getIndices(t)
		require.Equal(t, a, x)
		require.Equal(t, b, y)
	}
}

func TestPlacementPolicy_ProcessSelectors(t *testing.T) {
	p := newPlacementPolicy(2, nil,
		[]Selector{
			newSelector("SameRU", "City", 2, "FromRU", (*Selector).SelectSame),
			newSelector("DistinctRU", "City", 2, "FromRU", (*Selector).SelectDistinct),
			newSelector("Good", "Country", 2, "Good", (*Selector).SelectDistinct),
			newSelector("Main", "Country", 3, "*", (*Selector).SelectDistinct),
		},
		[]Filter{
			newFilter("FromRU", "Country", "Russia", netmap.EQ),
			newFilter("Good", "Rating", "4", netmap.GE),
		})
	nodes := []NodeInfo{
		nodeInfoFromAttributes("Country", "Russia", "Rating", "1", "City", "SPB"),
		nodeInfoFromAttributes("Country", "Germany", "Rating", "5", "City", "Berlin"),
		nodeInfoFromAttributes("Country", "Russia", "Rating", "6", "City", "Moscow"),
		nodeInfoFromAttributes("Country", "France", "Rating", "4", "City", "Paris"),
		nodeInfoFromAttributes("Country", "France", "Rating", "1", "City", "Lyon"),
		nodeInfoFromAttributes("Country", "Russia", "Rating", "5", "City", "SPB"),
		nodeInfoFromAttributes("Country", "Russia", "Rating", "7", "City", "Moscow"),
		nodeInfoFromAttributes("Country", "Germany", "Rating", "3", "City", "Darmstadt"),
		nodeInfoFromAttributes("Country", "Germany", "Rating", "7", "City", "Frankfurt"),
		nodeInfoFromAttributes("Country", "Russia", "Rating", "9", "City", "SPB"),
		nodeInfoFromAttributes("Country", "Russia", "Rating", "9", "City", "SPB"),
	}

	var nm NetMap
	nm.SetNodes(nodes)
	c := newContext(nm)
	c.setCBF(p.backupFactor)
	require.NoError(t, c.processFilters(p))
	require.NoError(t, c.processSelectors(p))

	for _, s := range p.selectors {
		sel := c.selections[s.GetName()]
		s := c.processedSelectors[s.GetName()]
		bucketCount, nodesInBucket := calcNodesCount(*s)
		nodesInBucket *= int(c.cbf)
		targ := fmt.Sprintf("selector '%s'", s.GetName())
		require.Equal(t, bucketCount, len(sel), targ)
		fName := s.GetFilter()
		for _, res := range sel {
			require.Equal(t, nodesInBucket, len(res), targ)
			for j := range res {
				require.True(t, fName == mainFilterName || c.match(c.processedFilters[fName], res[j]), targ)
			}
		}
	}
}

func TestSelector_SetName(t *testing.T) {
	const name = "some name"
	var s Selector

	require.Zero(t, s.m.GetName())

	s.SetName(name)
	require.Equal(t, name, s.m.GetName())
}

func TestSelector_SetNumberOfNodes(t *testing.T) {
	const num = 3
	var s Selector

	require.Zero(t, s.m.GetCount())

	s.SetNumberOfNodes(num)

	require.EqualValues(t, num, s.m.GetCount())
}

func TestSelectorClauses(t *testing.T) {
	var s Selector

	require.Equal(t, netmap.UnspecifiedClause, s.m.GetClause())

	s.SelectDistinct()
	require.Equal(t, netmap.Distinct, s.m.GetClause())

	s.SelectSame()
	require.Equal(t, netmap.Same, s.m.GetClause())
}

func TestSelector_SelectByBucketAttribute(t *testing.T) {
	const attr = "some attribute"
	var s Selector

	require.Zero(t, s.m.GetAttribute())

	s.SelectByBucketAttribute(attr)
	require.Equal(t, attr, s.m.GetAttribute())
}

func TestSelector_SetFilterName(t *testing.T) {
	const fName = "some filter"
	var s Selector

	require.Zero(t, s.m.GetFilter())

	s.SetFilterName(fName)
	require.Equal(t, fName, s.m.GetFilter())
}
