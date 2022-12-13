package netmap

import (
	"github.com/TrueCloudLab/frostfs-api-go/v2/netmap"
)

func newFilter(name string, k, v string, op netmap.Operation, fs ...Filter) (f Filter) {
	f.SetName(name)
	f.m.SetKey(k)
	f.m.SetOp(op)
	f.m.SetValue(v)
	inner := make([]netmap.Filter, len(fs))
	for i := range fs {
		inner[i] = fs[i].m
	}
	f.m.SetFilters(inner)
	return f
}

func newSelector(name string, attr string, count uint32, filter string, clause func(*Selector)) (s Selector) {
	s.SetName(name)
	s.SelectByBucketAttribute(attr)
	s.SetNumberOfNodes(count)
	clause(&s)
	s.SetFilterName(filter)
	return s
}

func newPlacementPolicy(bf uint32, rs []ReplicaDescriptor, ss []Selector, fs []Filter) (p PlacementPolicy) {
	p.SetContainerBackupFactor(bf)
	p.AddReplicas(rs...)
	p.AddSelectors(ss...)
	p.AddFilters(fs...)
	return p
}

func newReplica(c uint32, s string) (r ReplicaDescriptor) {
	r.SetNumberOfObjects(c)
	r.SetSelectorName(s)
	return r
}

func nodeInfoFromAttributes(props ...string) (n NodeInfo) {
	for i := 0; i < len(props); i += 2 {
		n.SetAttribute(props[i], props[i+1])
	}

	return
}
