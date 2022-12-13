package accounting_test

import (
	"testing"

	v2accounting "github.com/TrueCloudLab/frostfs-api-go/v2/accounting"
	"github.com/TrueCloudLab/frostfs-sdk-go/accounting"
	"github.com/stretchr/testify/require"
)

func TestDecimalData(t *testing.T) {
	const v, p = 4, 2

	var d accounting.Decimal

	require.Zero(t, d.Value())
	require.Zero(t, d.Precision())

	d.SetValue(v)
	d.SetPrecision(p)

	require.EqualValues(t, v, d.Value())
	require.EqualValues(t, p, d.Precision())
}

func TestDecimalMessageV2(t *testing.T) {
	var (
		d accounting.Decimal
		m v2accounting.Decimal
	)

	m.SetValue(7)
	m.SetPrecision(8)

	require.NoError(t, d.ReadFromV2(m))

	require.EqualValues(t, m.GetValue(), d.Value())
	require.EqualValues(t, m.GetPrecision(), d.Precision())

	var m2 v2accounting.Decimal

	d.WriteToV2(&m2)

	require.EqualValues(t, d.Value(), m2.GetValue())
	require.EqualValues(t, d.Precision(), m2.GetPrecision())
}
