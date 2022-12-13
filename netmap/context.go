package netmap

import (
	"errors"

	"github.com/TrueCloudLab/frostfs-api-go/v2/netmap"
	"github.com/TrueCloudLab/hrw"
)

// context of a placement build process.
type context struct {
	// network map to operate on
	netMap NetMap

	// cache of processed filters
	processedFilters map[string]*netmap.Filter

	// cache of processed selectors
	processedSelectors map[string]*netmap.Selector

	// stores results of selector processing
	selections map[string][]nodes

	// cache of parsed numeric values
	numCache map[string]uint64

	hrwSeed []byte

	// hrw.Hash of hrwSeed
	hrwSeedHash uint64

	// weightFunc is a weighting function for determining node priority
	// which combines low price and high performance
	weightFunc weightFunc

	// container backup factor
	cbf uint32
}

// Various validation errors.
var (
	errInvalidFilterName = errors.New("filter name is invalid")
	errInvalidNumber     = errors.New("invalid number")
	errInvalidFilterOp   = errors.New("invalid filter operation")
	errFilterNotFound    = errors.New("filter not found")
	errNonEmptyFilters   = errors.New("simple filter contains sub-filters")
	errNotEnoughNodes    = errors.New("not enough nodes to SELECT from")
	errUnnamedTopFilter  = errors.New("unnamed top-level filter")
)

// newContext returns initialized context.
func newContext(nm NetMap) *context {
	return &context{
		netMap:             nm,
		processedFilters:   make(map[string]*netmap.Filter),
		processedSelectors: make(map[string]*netmap.Selector),
		selections:         make(map[string][]nodes),

		numCache:   make(map[string]uint64),
		weightFunc: defaultWeightFunc(nm.nodes),
	}
}

func (c *context) setPivot(pivot []byte) {
	if len(pivot) != 0 {
		c.hrwSeed = pivot
		c.hrwSeedHash = hrw.Hash(pivot)
	}
}

func (c *context) setCBF(cbf uint32) {
	if cbf == 0 {
		c.cbf = 3
	} else {
		c.cbf = cbf
	}
}

func defaultWeightFunc(ns nodes) weightFunc {
	mean := newMeanAgg()
	min := newMinAgg()

	for i := range ns {
		mean.Add(float64(ns[i].capacity()))
		min.Add(float64(ns[i].Price()))
	}

	return newWeightFunc(
		newSigmoidNorm(mean.Compute()),
		newReverseMinNorm(min.Compute()))
}
