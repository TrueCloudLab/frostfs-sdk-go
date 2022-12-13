package object

import (
	"testing"

	"github.com/TrueCloudLab/frostfs-api-go/v2/object"
	"github.com/stretchr/testify/require"
)

func TestAttribute(t *testing.T) {
	key, val := "some key", "some value"

	a := NewAttribute()
	a.SetKey(key)
	a.SetValue(val)

	require.Equal(t, key, a.Key())
	require.Equal(t, val, a.Value())

	aV2 := a.ToV2()

	require.Equal(t, key, aV2.GetKey())
	require.Equal(t, val, aV2.GetValue())
}

func TestAttributeEncoding(t *testing.T) {
	a := NewAttribute()
	a.SetKey("key")
	a.SetValue("value")

	t.Run("binary", func(t *testing.T) {
		data, err := a.Marshal()
		require.NoError(t, err)

		a2 := NewAttribute()
		require.NoError(t, a2.Unmarshal(data))

		require.Equal(t, a, a2)
	})

	t.Run("json", func(t *testing.T) {
		data, err := a.MarshalJSON()
		require.NoError(t, err)

		a2 := NewAttribute()
		require.NoError(t, a2.UnmarshalJSON(data))

		require.Equal(t, a, a2)
	})
}

func TestNewAttributeFromV2(t *testing.T) {
	t.Run("from nil", func(t *testing.T) {
		var x *object.Attribute

		require.Nil(t, NewAttributeFromV2(x))
	})
}

func TestAttribute_ToV2(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		var x *Attribute

		require.Nil(t, x.ToV2())
	})
}

func TestNewAttribute(t *testing.T) {
	t.Run("default values", func(t *testing.T) {
		a := NewAttribute()

		// check initial values
		require.Empty(t, a.Key())
		require.Empty(t, a.Value())

		// convert to v2 message
		aV2 := a.ToV2()

		require.Empty(t, aV2.GetKey())
		require.Empty(t, aV2.GetValue())
	})
}
