package object_test

import (
	"testing"

	v2object "github.com/TrueCloudLab/frostfs-api-go/v2/object"
	"github.com/TrueCloudLab/frostfs-sdk-go/object"
	"github.com/stretchr/testify/require"
)

func TestType_ToV2(t *testing.T) {
	typs := []struct {
		t  object.Type
		t2 v2object.Type
	}{
		{
			t:  object.TypeRegular,
			t2: v2object.TypeRegular,
		},
		{
			t:  object.TypeTombstone,
			t2: v2object.TypeTombstone,
		},
		{
			t:  object.TypeStorageGroup,
			t2: v2object.TypeStorageGroup,
		},
		{
			t:  object.TypeLock,
			t2: v2object.TypeLock,
		},
	}

	for _, item := range typs {
		t2 := item.t.ToV2()

		require.Equal(t, item.t2, t2)

		require.Equal(t, item.t, object.TypeFromV2(item.t2))
	}
}

func TestType_String(t *testing.T) {
	toPtr := func(v object.Type) *object.Type {
		return &v
	}

	testEnumStrings(t, new(object.Type), []enumStringItem{
		{val: toPtr(object.TypeTombstone), str: "TOMBSTONE"},
		{val: toPtr(object.TypeStorageGroup), str: "STORAGE_GROUP"},
		{val: toPtr(object.TypeRegular), str: "REGULAR"},
		{val: toPtr(object.TypeLock), str: "LOCK"},
	})
}

type enumIface interface {
	FromString(string) bool
	String() string
}

type enumStringItem struct {
	val enumIface
	str string
}

func testEnumStrings(t *testing.T, e enumIface, items []enumStringItem) {
	for _, item := range items {
		require.Equal(t, item.str, item.val.String())

		s := item.val.String()

		require.True(t, e.FromString(s), s)

		require.EqualValues(t, item.val, e, item.val)
	}

	// incorrect strings
	for _, str := range []string{
		"some string",
		"undefined",
	} {
		require.False(t, e.FromString(str))
	}
}
