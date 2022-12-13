package netmap

import (
	"errors"
	"testing"

	"github.com/TrueCloudLab/frostfs-api-go/v2/netmap"
	"github.com/stretchr/testify/require"
)

func TestContext_ProcessFilters(t *testing.T) {
	fs := []Filter{
		newFilter("StorageSSD", "Storage", "SSD", netmap.EQ),
		newFilter("GoodRating", "Rating", "4", netmap.GE),
		newFilter("Main", "", "", netmap.AND,
			newFilter("StorageSSD", "", "", 0),
			newFilter("", "IntField", "123", netmap.LT),
			newFilter("GoodRating", "", "", 0)),
	}

	c := newContext(NetMap{})
	p := newPlacementPolicy(1, nil, nil, fs)
	require.NoError(t, c.processFilters(p))
	require.Equal(t, 3, len(c.processedFilters))
	for _, f := range fs {
		require.Equal(t, f.m, *c.processedFilters[f.m.GetName()])
	}

	require.Equal(t, uint64(4), c.numCache[fs[1].m.GetValue()])
	require.Equal(t, uint64(123), c.numCache[fs[2].m.GetFilters()[1].GetValue()])
}

func TestContext_ProcessFiltersInvalid(t *testing.T) {
	errTestCases := []struct {
		name   string
		filter Filter
		err    error
	}{
		{
			"UnnamedTop",
			newFilter("", "Storage", "SSD", netmap.EQ),
			errUnnamedTopFilter,
		},
		{
			"InvalidReference",
			newFilter("Main", "", "", netmap.AND,
				newFilter("StorageSSD", "", "", 0)),
			errFilterNotFound,
		},
		{
			"NonEmptyKeyed",
			newFilter("Main", "Storage", "SSD", netmap.EQ,
				newFilter("StorageSSD", "", "", 0)),
			errNonEmptyFilters,
		},
		{
			"InvalidNumber",
			newFilter("Main", "Rating", "three", netmap.GE),
			errInvalidNumber,
		},
		{
			"InvalidOp",
			newFilter("Main", "Rating", "3", 0),
			errInvalidFilterOp,
		},
		{
			"InvalidName",
			newFilter("*", "Rating", "3", netmap.GE),
			errInvalidFilterName,
		},
	}
	for _, tc := range errTestCases {
		t.Run(tc.name, func(t *testing.T) {
			c := newContext(NetMap{})
			p := newPlacementPolicy(1, nil, nil, []Filter{tc.filter})
			err := c.processFilters(p)
			require.True(t, errors.Is(err, tc.err), "got: %v", err)
		})
	}
}

func TestFilter_MatchSimple_InvalidOp(t *testing.T) {
	var b NodeInfo
	b.SetAttribute("Rating", "4")
	b.SetAttribute("Country", "Germany")

	f := newFilter("Main", "Rating", "5", netmap.EQ)
	c := newContext(NetMap{})
	p := newPlacementPolicy(1, nil, nil, []Filter{f})
	require.NoError(t, c.processFilters(p))

	// just for the coverage
	f.m.SetOp(0)
	require.False(t, c.match(&f.m, b))
}
