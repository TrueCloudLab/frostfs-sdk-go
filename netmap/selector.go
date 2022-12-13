package netmap

import (
	"fmt"
	"sort"

	"github.com/TrueCloudLab/frostfs-api-go/v2/netmap"
	subnetid "github.com/TrueCloudLab/frostfs-sdk-go/subnet/id"
	"github.com/TrueCloudLab/hrw"
)

// processSelectors processes selectors and returns error is any of them is invalid.
func (c *context) processSelectors(p PlacementPolicy) error {
	for i := range p.selectors {
		fName := p.selectors[i].GetFilter()
		if fName != mainFilterName {
			_, ok := c.processedFilters[p.selectors[i].GetFilter()]
			if !ok {
				return fmt.Errorf("%w: SELECT FROM '%s'", errFilterNotFound, fName)
			}
		}

		sName := p.selectors[i].GetName()

		c.processedSelectors[sName] = &p.selectors[i]

		result, err := c.getSelection(p, p.selectors[i])
		if err != nil {
			return err
		}

		c.selections[sName] = result
	}

	return nil
}

// calcNodesCount returns number of buckets and minimum number of nodes in every bucket
// for the given selector.
func calcNodesCount(s netmap.Selector) (int, int) {
	switch s.GetClause() {
	case netmap.Same:
		return 1, int(s.GetCount())
	default:
		return int(s.GetCount()), 1
	}
}

// calcBucketWeight computes weight for a node bucket.
func calcBucketWeight(ns nodes, a aggregator, wf weightFunc) float64 {
	for i := range ns {
		a.Add(wf(ns[i]))
	}

	return a.Compute()
}

// getSelection returns nodes grouped by s.attribute.
// Last argument specifies if more buckets can be used to fulfill CBF.
func (c *context) getSelection(p PlacementPolicy, s netmap.Selector) ([]nodes, error) {
	bucketCount, nodesInBucket := calcNodesCount(s)
	buckets := c.getSelectionBase(p.subnet, s)

	if len(buckets) < bucketCount {
		return nil, fmt.Errorf("%w: '%s'", errNotEnoughNodes, s.GetName())
	}

	// We need deterministic output in case there is no pivot.
	// If pivot is set, buckets are sorted by HRW.
	// However, because initial order influences HRW order for buckets with equal weights,
	// we also need to have deterministic input to HRW sorting routine.
	if len(c.hrwSeed) == 0 {
		if s.GetAttribute() == "" {
			sort.Slice(buckets, func(i, j int) bool {
				return less(buckets[i].nodes[0], buckets[j].nodes[0])
			})
		} else {
			sort.Slice(buckets, func(i, j int) bool {
				return buckets[i].attr < buckets[j].attr
			})
		}
	}

	maxNodesInBucket := nodesInBucket * int(c.cbf)
	res := make([]nodes, 0, len(buckets))
	fallback := make([]nodes, 0, len(buckets))

	for i := range buckets {
		ns := buckets[i].nodes
		if len(ns) >= maxNodesInBucket {
			res = append(res, ns[:maxNodesInBucket])
		} else if len(ns) >= nodesInBucket {
			fallback = append(fallback, ns)
		}
	}

	if len(res) < bucketCount {
		// Fallback to using minimum allowed backup factor (1).
		res = append(res, fallback...)
		if len(res) < bucketCount {
			return nil, fmt.Errorf("%w: '%s'", errNotEnoughNodes, s.GetName())
		}
	}

	if len(c.hrwSeed) != 0 {
		weights := make([]float64, len(res))
		for i := range res {
			weights[i] = calcBucketWeight(res[i], newMeanIQRAgg(), c.weightFunc)
		}

		hrw.SortSliceByWeightValue(res, weights, c.hrwSeedHash)
	}

	if s.GetAttribute() == "" {
		res, fallback = res[:bucketCount], res[bucketCount:]
		for i := range fallback {
			index := i % bucketCount
			if len(res[index]) >= maxNodesInBucket {
				break
			}
			res[index] = append(res[index], fallback[i]...)
		}
	}

	return res[:bucketCount], nil
}

type nodeAttrPair struct {
	attr  string
	nodes nodes
}

// getSelectionBase returns nodes grouped by selector attribute.
// It it guaranteed that each pair will contain at least one node.
func (c *context) getSelectionBase(subnetID subnetid.ID, s netmap.Selector) []nodeAttrPair {
	fName := s.GetFilter()
	f := c.processedFilters[fName]
	isMain := fName == mainFilterName
	result := []nodeAttrPair{}
	nodeMap := map[string][]NodeInfo{}
	attr := s.GetAttribute()

	for i := range c.netMap.nodes {
		// TODO(fyrchik): make `BelongsToSubnet` to accept pointer
		if !BelongsToSubnet(c.netMap.nodes[i], subnetID) {
			continue
		}
		if isMain || c.match(f, c.netMap.nodes[i]) {
			if attr == "" {
				// Default attribute is transparent identifier which is different for every node.
				result = append(result, nodeAttrPair{attr: "", nodes: nodes{c.netMap.nodes[i]}})
			} else {
				v := c.netMap.nodes[i].Attribute(attr)
				nodeMap[v] = append(nodeMap[v], c.netMap.nodes[i])
			}
		}
	}

	if attr != "" {
		for k, ns := range nodeMap {
			result = append(result, nodeAttrPair{attr: k, nodes: ns})
		}
	}

	if len(c.hrwSeed) != 0 {
		for i := range result {
			hrw.SortSliceByWeightValue(result[i].nodes, result[i].nodes.weights(c.weightFunc), c.hrwSeedHash)
		}
	}

	return result
}
