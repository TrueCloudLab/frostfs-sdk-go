package netmap_test

import (
	"strings"
	"testing"

	. "github.com/TrueCloudLab/frostfs-sdk-go/netmap"
	netmaptest "github.com/TrueCloudLab/frostfs-sdk-go/netmap/test"
	"github.com/stretchr/testify/require"
)

func TestEncode(t *testing.T) {
	testCases := []string{
		`REP 1 IN X
CBF 1
SELECT 2 IN SAME Location FROM * AS X`,

		`REP 1
SELECT 2 IN City FROM Good
FILTER Country EQ RU AS FromRU
FILTER @FromRU AND Rating GT 7 AS Good`,

		`REP 7 IN SPB
SELECT 1 IN City FROM SPBSSD AS SPB
FILTER City EQ SPB AND SSD EQ true OR City EQ SPB AND Rating GE 5 AS SPBSSD`,
	}

	var p PlacementPolicy

	for _, testCase := range testCases {
		require.NoError(t, p.DecodeString(testCase))

		var b strings.Builder
		require.NoError(t, p.WriteStringTo(&b))

		require.Equal(t, testCase, b.String())
	}

	invalidTestCases := []string{
		`?REP 1`,
		`REP 1 trailing garbage`,
	}

	for i := range invalidTestCases {
		require.Error(t, p.DecodeString(invalidTestCases[i]), "#%d", i)
	}
}

func TestPlacementPolicyEncoding(t *testing.T) {
	v := netmaptest.PlacementPolicy()

	t.Run("binary", func(t *testing.T) {
		var v2 PlacementPolicy
		require.NoError(t, v2.Unmarshal(v.Marshal()))

		require.Equal(t, v, v2)
	})

	t.Run("json", func(t *testing.T) {
		data, err := v.MarshalJSON()
		require.NoError(t, err)

		var v2 PlacementPolicy
		require.NoError(t, v2.UnmarshalJSON(data))

		require.Equal(t, v, v2)
	})
}
