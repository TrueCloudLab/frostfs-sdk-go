package container

import (
	"errors"
	"fmt"

	"github.com/TrueCloudLab/frostfs-api-go/v2/container"
	"github.com/TrueCloudLab/frostfs-api-go/v2/refs"
	cid "github.com/TrueCloudLab/frostfs-sdk-go/container/id"
)

// SizeEstimation groups information about estimation of the size of the data
// stored in the NeoFS container.
//
// SizeEstimation is mutually compatible with github.com/TrueCloudLab/frostfs-api-go/v2/container.UsedSpaceAnnouncement
// message. See ReadFromV2 / WriteToV2 methods.
type SizeEstimation struct {
	m container.UsedSpaceAnnouncement
}

// ReadFromV2 reads SizeEstimation from the container.UsedSpaceAnnouncement message.
// Checks if the message conforms to NeoFS API V2 protocol.
//
// See also WriteToV2.
func (x *SizeEstimation) ReadFromV2(m container.UsedSpaceAnnouncement) error {
	cnrV2 := m.GetContainerID()
	if cnrV2 == nil {
		return errors.New("missing container")
	}

	var cnr cid.ID

	err := cnr.ReadFromV2(*cnrV2)
	if err != nil {
		return fmt.Errorf("invalid container: %w", err)
	}

	x.m = m

	return nil
}

// WriteToV2 writes SizeEstimation into the container.UsedSpaceAnnouncement message.
// The message MUST NOT be nil.
//
// See also ReadFromV2.
func (x SizeEstimation) WriteToV2(m *container.UsedSpaceAnnouncement) {
	*m = x.m
}

// SetEpoch sets epoch when estimation of the container data size was calculated.
//
// See also Epoch.
func (x *SizeEstimation) SetEpoch(epoch uint64) {
	x.m.SetEpoch(epoch)
}

// Epoch return epoch set using SetEpoch.
//
// Zero SizeEstimation represents estimation in zero epoch.
func (x SizeEstimation) Epoch() uint64 {
	return x.m.GetEpoch()
}

// SetContainer specifies the container for which the amount of data is estimated.
// Required by the NeoFS API protocol.
//
// See also Container.
func (x *SizeEstimation) SetContainer(cnr cid.ID) {
	var cidV2 refs.ContainerID
	cnr.WriteToV2(&cidV2)

	x.m.SetContainerID(&cidV2)
}

// Container returns container set using SetContainer.
//
// Zero SizeEstimation is not bound to any container (returns zero) which is
// incorrect according to NeoFS API protocol.
func (x SizeEstimation) Container() (res cid.ID) {
	m := x.m.GetContainerID()
	if m != nil {
		err := res.ReadFromV2(*m)
		if err != nil {
			panic(fmt.Errorf("unexpected error from cid.ID.ReadFromV2: %w", err))
		}
	}

	return
}

// SetValue sets estimated amount of data (in bytes) in the specified container.
//
// See also Value.
func (x *SizeEstimation) SetValue(value uint64) {
	x.m.SetUsedSpace(value)
}

// Value returns data size estimation set using SetValue.
//
// Zero SizeEstimation has zero value.
func (x SizeEstimation) Value() uint64 {
	return x.m.GetUsedSpace()
}
