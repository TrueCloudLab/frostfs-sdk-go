package object

import (
	"testing"

	"github.com/TrueCloudLab/frostfs-api-go/v2/object"
	"github.com/stretchr/testify/require"
)

func TestRange_SetOffset(t *testing.T) {
	r := NewRange()

	off := uint64(13)
	r.SetOffset(off)

	require.Equal(t, off, r.GetOffset())
}

func TestRange_SetLength(t *testing.T) {
	r := NewRange()

	ln := uint64(7)
	r.SetLength(ln)

	require.Equal(t, ln, r.GetLength())
}

func TestNewRangeFromV2(t *testing.T) {
	t.Run("from nil", func(t *testing.T) {
		var x *object.Range

		require.Nil(t, NewRangeFromV2(x))
	})
}

func TestRange_ToV2(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		var x *Range

		require.Nil(t, x.ToV2())
	})
}

func TestNewRange(t *testing.T) {
	t.Run("default values", func(t *testing.T) {
		r := NewRange()

		// check initial values
		require.Zero(t, r.GetLength())
		require.Zero(t, r.GetOffset())

		// convert to v2 message
		rV2 := r.ToV2()

		require.Zero(t, rV2.GetLength())
		require.Zero(t, rV2.GetOffset())
	})
}
