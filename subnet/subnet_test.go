package subnet_test

import (
	"testing"

	. "github.com/TrueCloudLab/frostfs-sdk-go/subnet"
	subnetid "github.com/TrueCloudLab/frostfs-sdk-go/subnet/id"
	subnetidtest "github.com/TrueCloudLab/frostfs-sdk-go/subnet/id/test"
	subnettest "github.com/TrueCloudLab/frostfs-sdk-go/subnet/test"
	usertest "github.com/TrueCloudLab/frostfs-sdk-go/user/test"
	"github.com/stretchr/testify/require"
)

func TestInfoZero(t *testing.T) {
	var info Info

	require.Zero(t, info.ID())
	require.True(t, subnetid.IsZero(info.ID()))
}

func TestInfo_SetID(t *testing.T) {
	id := subnetidtest.ID()

	var info Info
	info.SetID(id)

	require.Equal(t, id, info.ID())
	require.True(t, AssertReference(info, id))
}

func TestInfo_SetOwner(t *testing.T) {
	id := *usertest.ID()

	var info Info
	info.SetOwner(id)

	require.Equal(t, id, info.Owner())
	require.True(t, AssertOwnership(info, id))
}

func TestInfo_Marshal(t *testing.T) {
	info := subnettest.Info()

	var info2 Info
	require.NoError(t, info2.Unmarshal(info.Marshal()))

	require.Equal(t, info, info2)
}
