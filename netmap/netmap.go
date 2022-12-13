package netmap

import (
	"fmt"

	"github.com/TrueCloudLab/frostfs-api-go/v2/netmap"
	"github.com/TrueCloudLab/hrw"
)

// NetMap represents NeoFS network map. It includes information about all
// storage nodes registered in NeoFS the network.
//
// NetMap is mutually compatible with github.com/TrueCloudLab/frostfs-api-go/v2/netmap.NetMap
// message. See ReadFromV2 / WriteToV2 methods.
//
// Instances can be created using built-in var declaration.
type NetMap struct {
	epoch uint64

	nodes []NodeInfo
}

// ReadFromV2 reads NetMap from the netmap.NetMap message. Checks if the
// message conforms to NeoFS API V2 protocol.
//
// See also WriteToV2.
func (m *NetMap) ReadFromV2(msg netmap.NetMap) error {
	var err error
	nodes := msg.Nodes()

	if nodes == nil {
		m.nodes = nil
	} else {
		m.nodes = make([]NodeInfo, len(nodes))

		for i := range nodes {
			err = m.nodes[i].ReadFromV2(nodes[i])
			if err != nil {
				return fmt.Errorf("invalid node info: %w", err)
			}
		}
	}

	m.epoch = msg.Epoch()

	return nil
}

// WriteToV2 writes NetMap to the netmap.NetMap message. The message
// MUST NOT be nil.
//
// See also ReadFromV2.
func (m NetMap) WriteToV2(msg *netmap.NetMap) {
	var nodes []netmap.NodeInfo

	if m.nodes != nil {
		nodes = make([]netmap.NodeInfo, len(m.nodes))

		for i := range m.nodes {
			m.nodes[i].WriteToV2(&nodes[i])
		}

		msg.SetNodes(nodes)
	}

	msg.SetEpoch(m.epoch)
}

// SetNodes sets information list about all storage nodes from the NeoFS network.
//
// Argument MUST NOT be mutated, make a copy first.
//
// See also Nodes.
func (m *NetMap) SetNodes(nodes []NodeInfo) {
	m.nodes = nodes
}

// Nodes returns nodes set using SetNodes.
//
// Return value MUST not be mutated, make a copy first.
func (m NetMap) Nodes() []NodeInfo {
	return m.nodes
}

// SetEpoch specifies revision number of the NetMap.
//
// See also Epoch.
func (m *NetMap) SetEpoch(epoch uint64) {
	m.epoch = epoch
}

// Epoch returns epoch set using SetEpoch.
//
// Zero NetMap has zero revision.
func (m NetMap) Epoch() uint64 {
	return m.epoch
}

// nodes is a slice of NodeInfo instances needed for HRW sorting.
type nodes []NodeInfo

// assert nodes type provides hrw.Hasher required for HRW sorting.
var _ hrw.Hasher = nodes{}

// Hash is a function from hrw.Hasher interface. It is implemented
// to support weighted hrw sorting of buckets. Each bucket is already sorted by hrw,
// thus giving us needed "randomness".
func (n nodes) Hash() uint64 {
	if len(n) > 0 {
		return n[0].Hash()
	}

	return 0
}

// weights returns slice of nodes weights W.
func (n nodes) weights(wf weightFunc) []float64 {
	w := make([]float64, 0, len(n))
	for i := range n {
		w = append(w, wf(n[i]))
	}

	return w
}

func flattenNodes(ns []nodes) nodes {
	var sz, i int

	for i = range ns {
		sz += len(ns[i])
	}

	result := make(nodes, 0, sz)

	for i := range ns {
		result = append(result, ns[i]...)
	}

	return result
}

// PlacementVectors sorts container nodes returned by ContainerNodes method
// and returns placement vectors for the entity identified by the given pivot.
// For example, in order to build node list to store the object, binary-encoded
// object identifier can be used as pivot. Result is deterministic for
// the fixed NetMap and parameters.
func (m NetMap) PlacementVectors(vectors [][]NodeInfo, pivot []byte) ([][]NodeInfo, error) {
	h := hrw.Hash(pivot)
	wf := defaultWeightFunc(m.nodes)
	result := make([][]NodeInfo, len(vectors))

	for i := range vectors {
		result[i] = make([]NodeInfo, len(vectors[i]))
		copy(result[i], vectors[i])
		hrw.SortSliceByWeightValue(result[i], nodes(result[i]).weights(wf), h)
	}

	return result, nil
}

// ContainerNodes returns two-dimensional list of nodes as a result of applying
// given PlacementPolicy to the NetMap. Each line of the list corresponds to a
// replica descriptor. Line order corresponds to order of ReplicaDescriptor list
// in the policy. Nodes are pre-filtered according to the Filter list from
// the policy, and then selected by Selector list. Result is deterministic for
// the fixed NetMap and parameters.
//
// Result can be used in PlacementVectors.
func (m NetMap) ContainerNodes(p PlacementPolicy, pivot []byte) ([][]NodeInfo, error) {
	c := newContext(m)
	c.setPivot(pivot)
	c.setCBF(p.backupFactor)

	if err := c.processFilters(p); err != nil {
		return nil, err
	}

	if err := c.processSelectors(p); err != nil {
		return nil, err
	}

	result := make([][]NodeInfo, len(p.replicas))

	for i := range p.replicas {
		sName := p.replicas[i].GetSelector()
		if sName == "" {
			if len(p.selectors) == 0 {
				var s netmap.Selector
				s.SetCount(p.replicas[i].GetCount())
				s.SetFilter(mainFilterName)

				nodes, err := c.getSelection(p, s)
				if err != nil {
					return nil, err
				}

				result[i] = flattenNodes(nodes)
			}

			for i := range p.selectors {
				result[i] = append(result[i], flattenNodes(c.selections[p.selectors[i].GetName()])...)
			}

			continue
		}

		nodes, ok := c.selections[sName]
		if !ok {
			return nil, fmt.Errorf("selector not found: REPLICA '%s'", sName)
		}

		result[i] = append(result[i], flattenNodes(nodes)...)
	}

	return result, nil
}
