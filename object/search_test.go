package object_test

import (
	"crypto/sha256"
	"math/rand"
	"testing"

	v2object "github.com/TrueCloudLab/frostfs-api-go/v2/object"
	"github.com/TrueCloudLab/frostfs-sdk-go/object"
	oid "github.com/TrueCloudLab/frostfs-sdk-go/object/id"
	"github.com/stretchr/testify/require"
)

var eqV2Matches = map[object.SearchMatchType]v2object.MatchType{
	object.MatchUnknown:        v2object.MatchUnknown,
	object.MatchStringEqual:    v2object.MatchStringEqual,
	object.MatchStringNotEqual: v2object.MatchStringNotEqual,
	object.MatchNotPresent:     v2object.MatchNotPresent,
	object.MatchCommonPrefix:   v2object.MatchCommonPrefix,
}

func TestMatch(t *testing.T) {
	t.Run("known matches", func(t *testing.T) {
		for matchType, matchTypeV2 := range eqV2Matches {
			require.Equal(t, matchTypeV2, matchType.ToV2())
			require.Equal(t, object.SearchMatchFromV2(matchTypeV2), matchType)
		}
	})

	t.Run("unknown matches", func(t *testing.T) {
		var unknownMatchType object.SearchMatchType

		for matchType := range eqV2Matches {
			unknownMatchType += matchType
		}

		unknownMatchType++

		require.Equal(t, unknownMatchType.ToV2(), v2object.MatchUnknown)

		var unknownMatchTypeV2 v2object.MatchType

		for _, matchTypeV2 := range eqV2Matches {
			unknownMatchTypeV2 += matchTypeV2
		}

		unknownMatchTypeV2++

		require.Equal(t, object.SearchMatchFromV2(unknownMatchTypeV2), object.MatchUnknown)
	})
}

func TestFilter(t *testing.T) {
	inputs := [][]string{
		{"user-header", "user-value"},
	}

	filters := object.NewSearchFilters()
	for i := range inputs {
		filters.AddFilter(inputs[i][0], inputs[i][1], object.MatchStringEqual)
	}

	require.Len(t, filters, len(inputs))
	for i := range inputs {
		require.Equal(t, inputs[i][0], filters[i].Header())
		require.Equal(t, inputs[i][1], filters[i].Value())
		require.Equal(t, object.MatchStringEqual, filters[i].Operation())
	}

	v2 := filters.ToV2()
	newFilters := object.NewSearchFiltersFromV2(v2)
	require.Equal(t, filters, newFilters)
}

func TestSearchFilters_AddRootFilter(t *testing.T) {
	fs := new(object.SearchFilters)

	fs.AddRootFilter()

	require.Len(t, *fs, 1)

	f := (*fs)[0]

	require.Equal(t, object.MatchUnknown, f.Operation())
	require.Equal(t, v2object.FilterPropertyRoot, f.Header())
	require.Equal(t, "", f.Value())
}

func TestSearchFilters_AddPhyFilter(t *testing.T) {
	fs := new(object.SearchFilters)

	fs.AddPhyFilter()

	require.Len(t, *fs, 1)

	f := (*fs)[0]

	require.Equal(t, object.MatchUnknown, f.Operation())
	require.Equal(t, v2object.FilterPropertyPhy, f.Header())
	require.Equal(t, "", f.Value())
}

func testOID() oid.ID {
	cs := [sha256.Size]byte{}

	rand.Read(cs[:])

	var id oid.ID
	id.SetSHA256(cs)

	return id
}

func TestSearchFilters_AddParentIDFilter(t *testing.T) {
	par := testOID()

	fs := object.SearchFilters{}
	fs.AddParentIDFilter(object.MatchStringEqual, par)

	fsV2 := fs.ToV2()

	require.Len(t, fsV2, 1)

	require.Equal(t, v2object.FilterHeaderParent, fsV2[0].GetKey())
	require.Equal(t, par.EncodeToString(), fsV2[0].GetValue())
	require.Equal(t, v2object.MatchStringEqual, fsV2[0].GetMatchType())
}

func TestSearchFilters_AddObjectIDFilter(t *testing.T) {
	id := testOID()

	fs := new(object.SearchFilters)
	fs.AddObjectIDFilter(object.MatchStringEqual, id)

	t.Run("v2", func(t *testing.T) {
		fsV2 := fs.ToV2()

		require.Len(t, fsV2, 1)

		require.Equal(t, v2object.FilterHeaderObjectID, fsV2[0].GetKey())
		require.Equal(t, id.EncodeToString(), fsV2[0].GetValue())
		require.Equal(t, v2object.MatchStringEqual, fsV2[0].GetMatchType())
	})
}

func TestSearchFilters_AddSplitIDFilter(t *testing.T) {
	id := object.NewSplitID()

	fs := new(object.SearchFilters)
	fs.AddSplitIDFilter(object.MatchStringEqual, id)

	t.Run("v2", func(t *testing.T) {
		fsV2 := fs.ToV2()

		require.Len(t, fsV2, 1)

		require.Equal(t, v2object.FilterHeaderSplitID, fsV2[0].GetKey())
		require.Equal(t, id.String(), fsV2[0].GetValue())
		require.Equal(t, v2object.MatchStringEqual, fsV2[0].GetMatchType())
	})
}

func TestSearchFilters_AddTypeFilter(t *testing.T) {
	typ := object.TypeTombstone

	fs := new(object.SearchFilters)
	fs.AddTypeFilter(object.MatchStringEqual, typ)

	t.Run("v2", func(t *testing.T) {
		fsV2 := fs.ToV2()

		require.Len(t, fsV2, 1)

		require.Equal(t, v2object.FilterHeaderObjectType, fsV2[0].GetKey())
		require.Equal(t, typ.String(), fsV2[0].GetValue())
		require.Equal(t, v2object.MatchStringEqual, fsV2[0].GetMatchType())
	})
}

func TestSearchFiltersEncoding(t *testing.T) {
	fs := object.NewSearchFilters()
	fs.AddFilter("key 1", "value 2", object.MatchStringEqual)
	fs.AddFilter("key 2", "value 2", object.MatchStringNotEqual)
	fs.AddFilter("key 2", "value 2", object.MatchCommonPrefix)

	t.Run("json", func(t *testing.T) {
		data, err := fs.MarshalJSON()
		require.NoError(t, err)

		fs2 := object.NewSearchFilters()
		require.NoError(t, fs2.UnmarshalJSON(data))

		require.Equal(t, fs, fs2)
	})
}

func TestSearchMatchType_String(t *testing.T) {
	toPtr := func(v object.SearchMatchType) *object.SearchMatchType {
		return &v
	}

	testEnumStrings(t, new(object.SearchMatchType), []enumStringItem{
		{val: toPtr(object.MatchCommonPrefix), str: "COMMON_PREFIX"},
		{val: toPtr(object.MatchStringEqual), str: "STRING_EQUAL"},
		{val: toPtr(object.MatchStringNotEqual), str: "STRING_NOT_EQUAL"},
		{val: toPtr(object.MatchNotPresent), str: "NOT_PRESENT"},
		{val: toPtr(object.MatchUnknown), str: "MATCH_TYPE_UNSPECIFIED"},
	})
}
