package object_test

import (
	"errors"
	"testing"

	"github.com/TrueCloudLab/frostfs-sdk-go/object"
	"github.com/stretchr/testify/require"
)

func TestNewSplitInfoError(t *testing.T) {
	var (
		si = generateSplitInfo()

		err         error = object.NewSplitInfoError(si)
		expectedErr *object.SplitInfoError
	)

	require.True(t, errors.As(err, &expectedErr))

	siErr, ok := err.(*object.SplitInfoError)
	require.True(t, ok)
	require.Equal(t, si, siErr.SplitInfo())
}

func generateSplitInfo() *object.SplitInfo {
	si := object.NewSplitInfo()
	si.SetSplitID(object.NewSplitID())
	si.SetLastPart(generateID())
	si.SetLink(generateID())

	return si
}
