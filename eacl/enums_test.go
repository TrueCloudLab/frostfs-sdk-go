package eacl_test

import (
	"testing"

	v2acl "github.com/TrueCloudLab/frostfs-api-go/v2/acl"
	"github.com/TrueCloudLab/frostfs-sdk-go/eacl"
	"github.com/stretchr/testify/require"
)

var (
	eqV2Actions = map[eacl.Action]v2acl.Action{
		eacl.ActionUnknown: v2acl.ActionUnknown,
		eacl.ActionAllow:   v2acl.ActionAllow,
		eacl.ActionDeny:    v2acl.ActionDeny,
	}

	eqV2Operations = map[eacl.Operation]v2acl.Operation{
		eacl.OperationUnknown:   v2acl.OperationUnknown,
		eacl.OperationGet:       v2acl.OperationGet,
		eacl.OperationHead:      v2acl.OperationHead,
		eacl.OperationPut:       v2acl.OperationPut,
		eacl.OperationDelete:    v2acl.OperationDelete,
		eacl.OperationSearch:    v2acl.OperationSearch,
		eacl.OperationRange:     v2acl.OperationRange,
		eacl.OperationRangeHash: v2acl.OperationRangeHash,
	}

	eqV2Roles = map[eacl.Role]v2acl.Role{
		eacl.RoleUnknown: v2acl.RoleUnknown,
		eacl.RoleUser:    v2acl.RoleUser,
		eacl.RoleSystem:  v2acl.RoleSystem,
		eacl.RoleOthers:  v2acl.RoleOthers,
	}

	eqV2Matches = map[eacl.Match]v2acl.MatchType{
		eacl.MatchUnknown:        v2acl.MatchTypeUnknown,
		eacl.MatchStringEqual:    v2acl.MatchTypeStringEqual,
		eacl.MatchStringNotEqual: v2acl.MatchTypeStringNotEqual,
	}

	eqV2HeaderTypes = map[eacl.FilterHeaderType]v2acl.HeaderType{
		eacl.HeaderTypeUnknown: v2acl.HeaderTypeUnknown,
		eacl.HeaderFromRequest: v2acl.HeaderTypeRequest,
		eacl.HeaderFromObject:  v2acl.HeaderTypeObject,
		eacl.HeaderFromService: v2acl.HeaderTypeService,
	}
)

func TestAction(t *testing.T) {
	t.Run("known actions", func(t *testing.T) {
		for i := eacl.ActionUnknown; i <= eacl.ActionDeny; i++ {
			require.Equal(t, eqV2Actions[i], i.ToV2())
			require.Equal(t, eacl.ActionFromV2(i.ToV2()), i)
		}
	})

	t.Run("unknown actions", func(t *testing.T) {
		require.Equal(t, (eacl.ActionDeny + 1).ToV2(), v2acl.ActionUnknown)
		require.Equal(t, eacl.ActionFromV2(v2acl.ActionDeny+1), eacl.ActionUnknown)
	})
}

func TestOperation(t *testing.T) {
	t.Run("known operations", func(t *testing.T) {
		for i := eacl.OperationUnknown; i <= eacl.OperationRangeHash; i++ {
			require.Equal(t, eqV2Operations[i], i.ToV2())
			require.Equal(t, eacl.OperationFromV2(i.ToV2()), i)
		}
	})

	t.Run("unknown operations", func(t *testing.T) {
		require.Equal(t, (eacl.OperationRangeHash + 1).ToV2(), v2acl.OperationUnknown)
		require.Equal(t, eacl.OperationFromV2(v2acl.OperationRangeHash+1), eacl.OperationUnknown)
	})
}

func TestRole(t *testing.T) {
	t.Run("known roles", func(t *testing.T) {
		for i := eacl.RoleUnknown; i <= eacl.RoleOthers; i++ {
			require.Equal(t, eqV2Roles[i], i.ToV2())
			require.Equal(t, eacl.RoleFromV2(i.ToV2()), i)
		}
	})

	t.Run("unknown roles", func(t *testing.T) {
		require.Equal(t, (eacl.RoleOthers + 1).ToV2(), v2acl.RoleUnknown)
		require.Equal(t, eacl.RoleFromV2(v2acl.RoleOthers+1), eacl.RoleUnknown)
	})
}

func TestMatch(t *testing.T) {
	t.Run("known matches", func(t *testing.T) {
		for i := eacl.MatchUnknown; i <= eacl.MatchStringNotEqual; i++ {
			require.Equal(t, eqV2Matches[i], i.ToV2())
			require.Equal(t, eacl.MatchFromV2(i.ToV2()), i)
		}
	})

	t.Run("unknown matches", func(t *testing.T) {
		require.Equal(t, (eacl.MatchStringNotEqual + 1).ToV2(), v2acl.MatchTypeUnknown)
		require.Equal(t, eacl.MatchFromV2(v2acl.MatchTypeStringNotEqual+1), eacl.MatchUnknown)
	})
}

func TestFilterHeaderType(t *testing.T) {
	t.Run("known header types", func(t *testing.T) {
		for i := eacl.HeaderTypeUnknown; i <= eacl.HeaderFromService; i++ {
			require.Equal(t, eqV2HeaderTypes[i], i.ToV2())
			require.Equal(t, eacl.FilterHeaderTypeFromV2(i.ToV2()), i)
		}
	})

	t.Run("unknown header types", func(t *testing.T) {
		require.Equal(t, (eacl.HeaderFromService + 1).ToV2(), v2acl.HeaderTypeUnknown)
		require.Equal(t, eacl.FilterHeaderTypeFromV2(v2acl.HeaderTypeService+1), eacl.HeaderTypeUnknown)
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
		"UNSPECIFIED",
	} {
		require.False(t, e.FromString(str))
	}
}

func TestAction_String(t *testing.T) {
	toPtr := func(v eacl.Action) *eacl.Action {
		return &v
	}

	testEnumStrings(t, new(eacl.Action), []enumStringItem{
		{val: toPtr(eacl.ActionAllow), str: "ALLOW"},
		{val: toPtr(eacl.ActionDeny), str: "DENY"},
		{val: toPtr(eacl.ActionUnknown), str: "ACTION_UNSPECIFIED"},
	})
}

func TestRole_String(t *testing.T) {
	toPtr := func(v eacl.Role) *eacl.Role {
		return &v
	}

	testEnumStrings(t, new(eacl.Role), []enumStringItem{
		{val: toPtr(eacl.RoleUser), str: "USER"},
		{val: toPtr(eacl.RoleSystem), str: "SYSTEM"},
		{val: toPtr(eacl.RoleOthers), str: "OTHERS"},
		{val: toPtr(eacl.RoleUnknown), str: "ROLE_UNSPECIFIED"},
	})
}

func TestOperation_String(t *testing.T) {
	toPtr := func(v eacl.Operation) *eacl.Operation {
		return &v
	}

	testEnumStrings(t, new(eacl.Operation), []enumStringItem{
		{val: toPtr(eacl.OperationGet), str: "GET"},
		{val: toPtr(eacl.OperationPut), str: "PUT"},
		{val: toPtr(eacl.OperationHead), str: "HEAD"},
		{val: toPtr(eacl.OperationDelete), str: "DELETE"},
		{val: toPtr(eacl.OperationSearch), str: "SEARCH"},
		{val: toPtr(eacl.OperationRange), str: "GETRANGE"},
		{val: toPtr(eacl.OperationRangeHash), str: "GETRANGEHASH"},
		{val: toPtr(eacl.OperationUnknown), str: "OPERATION_UNSPECIFIED"},
	})
}

func TestMatch_String(t *testing.T) {
	toPtr := func(v eacl.Match) *eacl.Match {
		return &v
	}

	testEnumStrings(t, new(eacl.Match), []enumStringItem{
		{val: toPtr(eacl.MatchStringEqual), str: "STRING_EQUAL"},
		{val: toPtr(eacl.MatchStringNotEqual), str: "STRING_NOT_EQUAL"},
		{val: toPtr(eacl.MatchUnknown), str: "MATCH_TYPE_UNSPECIFIED"},
	})
}

func TestFilterHeaderType_String(t *testing.T) {
	toPtr := func(v eacl.FilterHeaderType) *eacl.FilterHeaderType {
		return &v
	}

	testEnumStrings(t, new(eacl.FilterHeaderType), []enumStringItem{
		{val: toPtr(eacl.HeaderFromRequest), str: "REQUEST"},
		{val: toPtr(eacl.HeaderFromObject), str: "OBJECT"},
		{val: toPtr(eacl.HeaderTypeUnknown), str: "HEADER_UNSPECIFIED"},
	})
}
