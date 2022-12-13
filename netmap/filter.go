package netmap

import (
	"fmt"
	"strconv"

	"github.com/TrueCloudLab/frostfs-api-go/v2/netmap"
)

// mainFilterName is a name of the filter
// which points to the whole netmap.
const mainFilterName = "*"

// processFilters processes filters and returns error is any of them is invalid.
func (c *context) processFilters(p PlacementPolicy) error {
	for i := range p.filters {
		if err := c.processFilter(p.filters[i], true); err != nil {
			return fmt.Errorf("process filter #%d (%s): %w", i, p.filters[i].GetName(), err)
		}
	}

	return nil
}

func (c *context) processFilter(f netmap.Filter, top bool) error {
	fName := f.GetName()
	if fName == mainFilterName {
		return fmt.Errorf("%w: '%s' is reserved", errInvalidFilterName, mainFilterName)
	}

	if top && fName == "" {
		return errUnnamedTopFilter
	}

	if !top && fName != "" && c.processedFilters[fName] == nil {
		return errFilterNotFound
	}

	inner := f.GetFilters()

	switch op := f.GetOp(); op {
	case netmap.AND, netmap.OR:
		for i := range inner {
			if err := c.processFilter(inner[i], false); err != nil {
				return fmt.Errorf("process inner filter #%d: %w", i, err)
			}
		}
	default:
		if len(inner) != 0 {
			return errNonEmptyFilters
		} else if !top && fName != "" { // named reference
			return nil
		}

		switch op {
		case netmap.EQ, netmap.NE:
		case netmap.GT, netmap.GE, netmap.LT, netmap.LE:
			val := f.GetValue()
			n, err := strconv.ParseUint(val, 10, 64)
			if err != nil {
				return fmt.Errorf("%w: '%s'", errInvalidNumber, f.GetValue())
			}

			c.numCache[val] = n
		default:
			return fmt.Errorf("%w: %s", errInvalidFilterOp, op)
		}
	}

	if top {
		c.processedFilters[fName] = &f
	}

	return nil
}

// match matches f against b. It returns no errors because
// filter should have been parsed during context creation
// and missing node properties are considered as a regular fail.
func (c *context) match(f *netmap.Filter, b NodeInfo) bool {
	switch f.GetOp() {
	case netmap.AND, netmap.OR:
		inner := f.GetFilters()
		for i := range inner {
			fSub := &inner[i]
			if name := inner[i].GetName(); name != "" {
				fSub = c.processedFilters[name]
			}

			ok := c.match(fSub, b)
			if ok == (f.GetOp() == netmap.OR) {
				return ok
			}
		}

		return f.GetOp() == netmap.AND
	default:
		return c.matchKeyValue(f, b)
	}
}

func (c *context) matchKeyValue(f *netmap.Filter, b NodeInfo) bool {
	switch op := f.GetOp(); op {
	case netmap.EQ:
		return b.Attribute(f.GetKey()) == f.GetValue()
	case netmap.NE:
		return b.Attribute(f.GetKey()) != f.GetValue()
	default:
		var attr uint64

		switch f.GetKey() {
		case attrPrice:
			attr = b.Price()
		case attrCapacity:
			attr = b.capacity()
		default:
			var err error

			attr, err = strconv.ParseUint(b.Attribute(f.GetKey()), 10, 64)
			if err != nil {
				// Note: because filters are somewhat independent from nodes attributes,
				// We don't report an error here, and fail filter instead.
				return false
			}
		}

		switch op {
		case netmap.GT:
			return attr > c.numCache[f.GetValue()]
		case netmap.GE:
			return attr >= c.numCache[f.GetValue()]
		case netmap.LT:
			return attr < c.numCache[f.GetValue()]
		case netmap.LE:
			return attr <= c.numCache[f.GetValue()]
		default:
			// do nothing and return false
		}
	}
	// will not happen if context was created from f (maybe panic?)
	return false
}
