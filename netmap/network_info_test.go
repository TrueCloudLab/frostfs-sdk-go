package netmap_test

import (
	"encoding/binary"
	"math"
	"testing"

	"github.com/TrueCloudLab/frostfs-api-go/v2/netmap"
	. "github.com/TrueCloudLab/frostfs-sdk-go/netmap"
	"github.com/stretchr/testify/require"
)

func TestNetworkInfo_CurrentEpoch(t *testing.T) {
	var x NetworkInfo

	require.Zero(t, x.CurrentEpoch())

	const e = 13

	x.SetCurrentEpoch(e)

	require.EqualValues(t, e, x.CurrentEpoch())

	var m netmap.NetworkInfo
	x.WriteToV2(&m)

	require.EqualValues(t, e, m.GetCurrentEpoch())
}

func TestNetworkInfo_MagicNumber(t *testing.T) {
	var x NetworkInfo

	require.Zero(t, x.MagicNumber())

	const magic = 321

	x.SetMagicNumber(magic)

	require.EqualValues(t, magic, x.MagicNumber())

	var m netmap.NetworkInfo
	x.WriteToV2(&m)

	require.EqualValues(t, magic, m.GetMagicNumber())
}

func TestNetworkInfo_MsPerBlock(t *testing.T) {
	var x NetworkInfo

	require.Zero(t, x.MsPerBlock())

	const ms = 789

	x.SetMsPerBlock(ms)

	require.EqualValues(t, ms, x.MsPerBlock())

	var m netmap.NetworkInfo
	x.WriteToV2(&m)

	require.EqualValues(t, ms, m.GetMsPerBlock())
}

func testConfigValue(t *testing.T,
	getter func(x NetworkInfo) interface{},
	setter func(x *NetworkInfo, val interface{}),
	val1, val2 interface{},
	v2Key string, v2Val func(val interface{}) []byte,
) {
	var x NetworkInfo

	require.Zero(t, getter(x))

	checkVal := func(exp interface{}) {
		require.EqualValues(t, exp, getter(x))

		var m netmap.NetworkInfo
		x.WriteToV2(&m)

		require.EqualValues(t, 1, m.GetNetworkConfig().NumberOfParameters())
		found := false
		m.GetNetworkConfig().IterateParameters(func(prm *netmap.NetworkParameter) bool {
			require.False(t, found)
			require.Equal(t, []byte(v2Key), prm.GetKey())
			require.Equal(t, v2Val(exp), prm.GetValue())
			found = true
			return false
		})
		require.True(t, found)
	}

	setter(&x, val1)
	checkVal(val1)

	setter(&x, val2)
	checkVal(val2)
}

func TestNetworkInfo_AuditFee(t *testing.T) {
	testConfigValue(t,
		func(x NetworkInfo) interface{} { return x.AuditFee() },
		func(info *NetworkInfo, val interface{}) { info.SetAuditFee(val.(uint64)) },
		uint64(1), uint64(2),
		"AuditFee", func(val interface{}) []byte {
			data := make([]byte, 8)
			binary.LittleEndian.PutUint64(data, val.(uint64))
			return data
		},
	)
}

func TestNetworkInfo_StoragePrice(t *testing.T) {
	testConfigValue(t,
		func(x NetworkInfo) interface{} { return x.StoragePrice() },
		func(info *NetworkInfo, val interface{}) { info.SetStoragePrice(val.(uint64)) },
		uint64(1), uint64(2),
		"BasicIncomeRate", func(val interface{}) []byte {
			data := make([]byte, 8)
			binary.LittleEndian.PutUint64(data, val.(uint64))
			return data
		},
	)
}

func TestNetworkInfo_ContainerFee(t *testing.T) {
	testConfigValue(t,
		func(x NetworkInfo) interface{} { return x.ContainerFee() },
		func(info *NetworkInfo, val interface{}) { info.SetContainerFee(val.(uint64)) },
		uint64(1), uint64(2),
		"ContainerFee", func(val interface{}) []byte {
			data := make([]byte, 8)
			binary.LittleEndian.PutUint64(data, val.(uint64))
			return data
		},
	)
}

func TestNetworkInfo_NamedContainerFee(t *testing.T) {
	testConfigValue(t,
		func(x NetworkInfo) interface{} { return x.NamedContainerFee() },
		func(info *NetworkInfo, val interface{}) { info.SetNamedContainerFee(val.(uint64)) },
		uint64(1), uint64(2),
		"ContainerAliasFee", func(val interface{}) []byte {
			data := make([]byte, 8)
			binary.LittleEndian.PutUint64(data, val.(uint64))
			return data
		},
	)
}

func TestNetworkInfo_EigenTrustAlpha(t *testing.T) {
	testConfigValue(t,
		func(x NetworkInfo) interface{} { return x.EigenTrustAlpha() },
		func(info *NetworkInfo, val interface{}) { info.SetEigenTrustAlpha(val.(float64)) },
		0.1, 0.2,
		"EigenTrustAlpha", func(val interface{}) []byte {
			data := make([]byte, 8)
			binary.LittleEndian.PutUint64(data, math.Float64bits(val.(float64)))
			return data
		},
	)
}

func TestNetworkInfo_NumberOfEigenTrustIterations(t *testing.T) {
	testConfigValue(t,
		func(x NetworkInfo) interface{} { return x.NumberOfEigenTrustIterations() },
		func(info *NetworkInfo, val interface{}) { info.SetNumberOfEigenTrustIterations(val.(uint64)) },
		uint64(1), uint64(2),
		"EigenTrustIterations", func(val interface{}) []byte {
			data := make([]byte, 8)
			binary.LittleEndian.PutUint64(data, val.(uint64))
			return data
		},
	)
}

func TestNetworkInfo_IRCandidateFee(t *testing.T) {
	testConfigValue(t,
		func(x NetworkInfo) interface{} { return x.IRCandidateFee() },
		func(info *NetworkInfo, val interface{}) { info.SetIRCandidateFee(val.(uint64)) },
		uint64(1), uint64(2),
		"InnerRingCandidateFee", func(val interface{}) []byte {
			data := make([]byte, 8)
			binary.LittleEndian.PutUint64(data, val.(uint64))
			return data
		},
	)
}

func TestNetworkInfo_MaxObjectSize(t *testing.T) {
	testConfigValue(t,
		func(x NetworkInfo) interface{} { return x.MaxObjectSize() },
		func(info *NetworkInfo, val interface{}) { info.SetMaxObjectSize(val.(uint64)) },
		uint64(1), uint64(2),
		"MaxObjectSize", func(val interface{}) []byte {
			data := make([]byte, 8)
			binary.LittleEndian.PutUint64(data, val.(uint64))
			return data
		},
	)
}

func TestNetworkInfo_WithdrawalFee(t *testing.T) {
	testConfigValue(t,
		func(x NetworkInfo) interface{} { return x.WithdrawalFee() },
		func(info *NetworkInfo, val interface{}) { info.SetWithdrawalFee(val.(uint64)) },
		uint64(1), uint64(2),
		"WithdrawFee", func(val interface{}) []byte {
			data := make([]byte, 8)
			binary.LittleEndian.PutUint64(data, val.(uint64))
			return data
		},
	)
}

func TestNetworkInfo_HomomorphicHashingDisabled(t *testing.T) {
	testConfigValue(t,
		func(x NetworkInfo) interface{} { return x.HomomorphicHashingDisabled() },
		func(info *NetworkInfo, val interface{}) {
			if val.(bool) {
				info.DisableHomomorphicHashing()
			}
		},
		true, true, // it is impossible to enable hashing
		"HomomorphicHashingDisabled", func(val interface{}) []byte {
			data := make([]byte, 1)

			if val.(bool) {
				data[0] = 1
			}

			return data
		},
	)
}

func TestNetworkInfo_MaintenanceModeAllowed(t *testing.T) {
	testConfigValue(t,
		func(x NetworkInfo) interface{} { return x.MaintenanceModeAllowed() },
		func(info *NetworkInfo, val interface{}) {
			if val.(bool) {
				info.AllowMaintenanceMode()
			}
		},
		true, true,
		"MaintenanceModeAllowed", func(val interface{}) []byte {
			if val.(bool) {
				return []byte{1}
			}
			return []byte{0}
		},
	)
}
