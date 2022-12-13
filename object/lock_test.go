package object_test

import (
	"testing"

	"github.com/TrueCloudLab/frostfs-sdk-go/object"
	objecttest "github.com/TrueCloudLab/frostfs-sdk-go/object/test"
	"github.com/stretchr/testify/require"
)

func TestLockEncoding(t *testing.T) {
	l := *objecttest.Lock()

	t.Run("binary", func(t *testing.T) {
		data := l.Marshal()

		var l2 object.Lock
		require.NoError(t, l2.Unmarshal(data))

		require.Equal(t, l, l2)
	})
}

func TestWriteLock(t *testing.T) {
	l := *objecttest.Lock()
	var o object.Object

	object.WriteLock(&o, l)

	var l2 object.Lock

	require.NoError(t, object.ReadLock(&l2, o))
	require.Equal(t, l, l2)

	// corrupt payload
	o.Payload()[0]++

	require.Error(t, object.ReadLock(&l2, o))
}
