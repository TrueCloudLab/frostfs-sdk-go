package apistatus_test

import (
	"errors"
	"testing"

	apistatus "github.com/TrueCloudLab/frostfs-sdk-go/client/status"
	"github.com/stretchr/testify/require"
)

func TestErrors(t *testing.T) {
	t.Run("error source", func(t *testing.T) {
		err := errors.New("some error")

		st := apistatus.ErrToStatus(err)

		success := apistatus.IsSuccessful(st)
		require.False(t, success)

		res := apistatus.ErrFromStatus(st)

		require.Equal(t, err, res)
	})

	t.Run("non-error source", func(t *testing.T) {
		var st apistatus.Status = "any non-error type"

		success := apistatus.IsSuccessful(st)
		require.True(t, success)

		res := apistatus.ErrFromStatus(st)

		require.Nil(t, res)
	})
}
