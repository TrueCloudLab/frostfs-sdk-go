package subnet

import (
	"fmt"

	"github.com/TrueCloudLab/frostfs-api-go/v2/refs"
	"github.com/TrueCloudLab/frostfs-api-go/v2/subnet"
	subnetid "github.com/TrueCloudLab/frostfs-sdk-go/subnet/id"
	"github.com/TrueCloudLab/frostfs-sdk-go/user"
)

// Info represents information about NeoFS subnet.
//
// Instances can be created using built-in var declaration.
type Info struct {
	id subnetid.ID

	owner user.ID
}

// Marshal encodes Info into a binary format of the NeoFS API protocol
// (Protocol Buffers with direct field order).
//
// See also Unmarshal.
func (x Info) Marshal() []byte {
	var id refs.SubnetID
	x.id.WriteToV2(&id)

	var owner refs.OwnerID
	x.owner.WriteToV2(&owner)

	var m subnet.Info
	m.SetID(&id)
	m.SetOwner(&owner)

	return m.StableMarshal(nil)
}

// Unmarshal decodes binary Info calculated using Marshal. Returns an error
// describing a format violation.
func (x *Info) Unmarshal(data []byte) error {
	var m subnet.Info

	err := m.Unmarshal(data)
	if err != nil {
		return err
	}

	id := m.ID()
	if id != nil {
		err = x.id.ReadFromV2(*id)
		if err != nil {
			return fmt.Errorf("invalid ID: %w", err)
		}
	} else {
		subnetid.MakeZero(&x.id)
	}

	owner := m.Owner()
	if owner != nil {
		err = x.owner.ReadFromV2(*owner)
		if err != nil {
			return fmt.Errorf("invalid owner: %w", err)
		}
	} else {
		x.owner = user.ID{}
	}

	return nil
}

// SetID sets the identifier of the subnet that Info describes.
//
// See also ID.
func (x *Info) SetID(id subnetid.ID) {
	x.id = id
}

// ID returns subnet identifier set using SetID.
//
// Zero Info refers to the zero subnet.
func (x Info) ID() subnetid.ID {
	return x.id
}

// SetOwner sets identifier of the subnet owner.
func (x *Info) SetOwner(id user.ID) {
	x.owner = id
}

// Owner returns subnet owner set using SetOwner.
//
// Zero Info has no owner which is incorrect according to the
// NeoFS API protocol.
func (x Info) Owner() user.ID {
	return x.owner
}

// AssertOwnership checks if the given info describes the subnet owned by the
// given user.
func AssertOwnership(info Info, id user.ID) bool {
	return id.Equals(info.Owner())
}

// AssertReference checks if the given info describes the subnet referenced by
// the given id.
func AssertReference(info Info, id subnetid.ID) bool {
	return id.Equals(info.ID())
}
