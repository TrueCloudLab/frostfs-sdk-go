package container_test

import (
	"crypto/sha256"
	"testing"

	v2container "github.com/TrueCloudLab/frostfs-api-go/v2/container"
	"github.com/TrueCloudLab/frostfs-api-go/v2/refs"
	"github.com/TrueCloudLab/frostfs-sdk-go/container"
	cid "github.com/TrueCloudLab/frostfs-sdk-go/container/id"
	cidtest "github.com/TrueCloudLab/frostfs-sdk-go/container/id/test"
	"github.com/stretchr/testify/require"
)

func TestSizeEstimation_Epoch(t *testing.T) {
	var val container.SizeEstimation

	require.Zero(t, val.Epoch())

	const epoch = 123

	val.SetEpoch(epoch)
	require.EqualValues(t, epoch, val.Epoch())

	var msg v2container.UsedSpaceAnnouncement
	val.WriteToV2(&msg)

	require.EqualValues(t, epoch, msg.GetEpoch())
}

func TestSizeEstimation_Container(t *testing.T) {
	var val container.SizeEstimation

	require.Zero(t, val.Container())

	cnr := cidtest.ID()

	val.SetContainer(cnr)
	require.True(t, val.Container().Equals(cnr))

	var msg v2container.UsedSpaceAnnouncement
	val.WriteToV2(&msg)

	var msgCnr refs.ContainerID
	cnr.WriteToV2(&msgCnr)

	require.Equal(t, &msgCnr, msg.GetContainerID())
}

func TestSizeEstimation_Value(t *testing.T) {
	var val container.SizeEstimation

	require.Zero(t, val.Value())

	const value = 876

	val.SetValue(value)
	require.EqualValues(t, value, val.Value())

	var msg v2container.UsedSpaceAnnouncement
	val.WriteToV2(&msg)

	require.EqualValues(t, value, msg.GetUsedSpace())
}

func TestSizeEstimation_ReadFromV2(t *testing.T) {
	const epoch = 654
	const value = 903
	var cnrMsg refs.ContainerID

	var msg v2container.UsedSpaceAnnouncement

	var val container.SizeEstimation

	require.Error(t, val.ReadFromV2(msg))

	msg.SetContainerID(&cnrMsg)

	require.Error(t, val.ReadFromV2(msg))

	cnrMsg.SetValue(make([]byte, sha256.Size))

	var cnr cid.ID
	require.NoError(t, cnr.ReadFromV2(cnrMsg))

	msg.SetEpoch(epoch)
	msg.SetUsedSpace(value)

	require.NoError(t, val.ReadFromV2(msg))

	require.EqualValues(t, epoch, val.Epoch())
	require.EqualValues(t, value, val.Value())
	require.EqualValues(t, cnr, val.Container())
}
