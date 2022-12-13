package eacl

import (
	"strconv"
	"testing"

	"github.com/TrueCloudLab/frostfs-api-go/v2/acl"
	v2acl "github.com/TrueCloudLab/frostfs-api-go/v2/acl"
	"github.com/stretchr/testify/require"
)

func newObjectFilter(match Match, key, val string) *Filter {
	return &Filter{
		from: HeaderFromObject,
		key: filterKey{
			str: key,
		},
		matcher: match,
		value:   staticStringer(val),
	}
}

func TestFilter(t *testing.T) {
	filter := newObjectFilter(MatchStringEqual, "some name", "200")

	v2 := filter.ToV2()
	require.NotNil(t, v2)
	require.Equal(t, v2acl.HeaderTypeObject, v2.GetHeaderType())
	require.EqualValues(t, v2acl.MatchTypeStringEqual, v2.GetMatchType())
	require.Equal(t, filter.Key(), v2.GetKey())
	require.Equal(t, filter.Value(), v2.GetValue())

	newFilter := NewFilterFromV2(v2)
	require.Equal(t, filter, newFilter)

	t.Run("from nil v2 filter", func(t *testing.T) {
		require.Equal(t, new(Filter), NewFilterFromV2(nil))
	})
}

func TestFilterEncoding(t *testing.T) {
	f := newObjectFilter(MatchStringEqual, "key", "value")

	t.Run("binary", func(t *testing.T) {
		data, err := f.Marshal()
		require.NoError(t, err)

		f2 := NewFilter()
		require.NoError(t, f2.Unmarshal(data))

		require.Equal(t, f, f2)
	})

	t.Run("json", func(t *testing.T) {
		data, err := f.MarshalJSON()
		require.NoError(t, err)

		d2 := NewFilter()
		require.NoError(t, d2.UnmarshalJSON(data))

		require.Equal(t, f, d2)
	})
}

func TestFilter_ToV2(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		var x *Filter

		require.Nil(t, x.ToV2())
	})

	t.Run("default values", func(t *testing.T) {
		filter := NewFilter()

		// check initial values
		require.Empty(t, filter.Key())
		require.Empty(t, filter.Value())
		require.Equal(t, HeaderTypeUnknown, filter.From())
		require.Equal(t, MatchUnknown, filter.Matcher())

		// convert to v2 message
		filterV2 := filter.ToV2()

		require.Empty(t, filterV2.GetKey())
		require.Empty(t, filterV2.GetValue())
		require.Equal(t, acl.HeaderTypeUnknown, filterV2.GetHeaderType())
		require.Equal(t, acl.MatchTypeUnknown, filterV2.GetMatchType())
	})

	t.Run("reserved types", func(t *testing.T) {
		r := NewRecord()
		for i := filterKeyType(1); i < fKeyObjLast; i++ {
			r.addObjectReservedFilter(MatchStringEqual, i, staticStringer(strconv.FormatUint(uint64(i), 16)))
		}

		for i := range r.filters {
			fv2 := r.filters[i].ToV2()
			actual := NewFilterFromV2(fv2)
			require.Equal(t, actual, &r.filters[i])
		}
	})
}
