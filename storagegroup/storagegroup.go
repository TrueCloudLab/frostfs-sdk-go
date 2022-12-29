package storagegroup

import (
	"errors"
	"fmt"
	"strconv"

	objectV2 "github.com/TrueCloudLab/frostfs-api-go/v2/object"
	"github.com/TrueCloudLab/frostfs-api-go/v2/refs"
	"github.com/TrueCloudLab/frostfs-api-go/v2/storagegroup"
	"github.com/TrueCloudLab/frostfs-sdk-go/checksum"
	objectSDK "github.com/TrueCloudLab/frostfs-sdk-go/object"
	oid "github.com/TrueCloudLab/frostfs-sdk-go/object/id"
)

// StorageGroup represents storage group of the FrostFS objects.
//
// StorageGroup is mutually compatible with github.com/TrueCloudLab/frostfs-api-go/v2/storagegroup.StorageGroup
// message. See ReadFromMessageV2 / WriteToMessageV2 methods.
//
// Instances can be created using built-in var declaration.
//
// Note that direct typecast is not safe and may result in loss of compatibility:
//
//	_ = StorageGroup(storagegroup.StorageGroup) // not recommended
type StorageGroup storagegroup.StorageGroup

// reads StorageGroup from the storagegroup.StorageGroup message. If checkFieldPresence is set,
// returns an error on absence of any protocol-required field.
func (sg *StorageGroup) readFromV2(m storagegroup.StorageGroup, checkFieldPresence bool) error {
	var err error

	h := m.GetValidationHash()
	if h != nil {
		err = new(checksum.Checksum).ReadFromV2(*h)
		if err != nil {
			return fmt.Errorf("invalid hash: %w", err)
		}
	} else if checkFieldPresence {
		return errors.New("missing hash")
	}

	members := m.GetMembers()
	if len(members) > 0 {
		var member oid.ID
		mMembers := make(map[oid.ID]struct{}, len(members))
		var exits bool

		for i := range members {
			err = member.ReadFromV2(members[i])
			if err != nil {
				return fmt.Errorf("invalid member: %w", err)
			}

			_, exits = mMembers[member]
			if exits {
				return fmt.Errorf("duplicated member %s", member)
			}

			mMembers[member] = struct{}{}
		}
	} else if checkFieldPresence {
		return errors.New("missing members")
	}

	*sg = StorageGroup(m)

	return nil
}

// ReadFromV2 reads StorageGroup from the storagegroup.StorageGroup message.
// Checks if the message conforms to FrostFS API V2 protocol.
//
// See also WriteToV2.
func (sg *StorageGroup) ReadFromV2(m storagegroup.StorageGroup) error {
	return sg.readFromV2(m, true)
}

// WriteToV2 writes StorageGroup to the storagegroup.StorageGroup message.
// The message must not be nil.
//
// See also ReadFromV2.
func (sg StorageGroup) WriteToV2(m *storagegroup.StorageGroup) {
	*m = (storagegroup.StorageGroup)(sg)
}

// ValidationDataSize returns total size of the payloads
// of objects in the storage group.
//
// Zero StorageGroup has 0 data size.
//
// See also SetValidationDataSize.
func (sg StorageGroup) ValidationDataSize() uint64 {
	v2 := (storagegroup.StorageGroup)(sg)
	return v2.GetValidationDataSize()
}

// SetValidationDataSize sets total size of the payloads
// of objects in the storage group.
//
// See also ValidationDataSize.
func (sg *StorageGroup) SetValidationDataSize(epoch uint64) {
	(*storagegroup.StorageGroup)(sg).SetValidationDataSize(epoch)
}

// ValidationDataHash returns homomorphic hash from the
// concatenation of the payloads of the storage group members
// and bool that indicates checksum presence in the storage
// group.
//
// Zero StorageGroup does not have validation data checksum.
//
// See also SetValidationDataHash.
func (sg StorageGroup) ValidationDataHash() (v checksum.Checksum, isSet bool) {
	v2 := (storagegroup.StorageGroup)(sg)
	if checksumV2 := v2.GetValidationHash(); checksumV2 != nil {
		v.ReadFromV2(*checksumV2) // FIXME(@cthulhu-rider): #226 handle error
		isSet = true
	}

	return
}

// SetValidationDataHash sets homomorphic hash from the
// concatenation of the payloads of the storage group members.
//
// See also ValidationDataHash.
func (sg *StorageGroup) SetValidationDataHash(hash checksum.Checksum) {
	var v2 refs.Checksum
	hash.WriteToV2(&v2)

	(*storagegroup.StorageGroup)(sg).SetValidationHash(&v2)
}

// ExpirationEpoch returns last FrostFS epoch number
// of the storage group lifetime.
//
// Zero StorageGroup has 0 expiration epoch.
//
// See also SetExpirationEpoch.
func (sg StorageGroup) ExpirationEpoch() uint64 {
	v2 := (storagegroup.StorageGroup)(sg)
	return v2.GetExpirationEpoch()
}

// SetExpirationEpoch sets last FrostFS epoch number
// of the storage group lifetime.
//
// See also ExpirationEpoch.
func (sg *StorageGroup) SetExpirationEpoch(epoch uint64) {
	(*storagegroup.StorageGroup)(sg).SetExpirationEpoch(epoch)
}

// Members returns strictly ordered list of
// storage group member objects.
//
// Zero StorageGroup has nil members value.
//
// See also SetMembers.
func (sg StorageGroup) Members() []oid.ID {
	v2 := (storagegroup.StorageGroup)(sg)
	mV2 := v2.GetMembers()

	if mV2 == nil {
		return nil
	}

	m := make([]oid.ID, len(mV2))

	for i := range mV2 {
		_ = m[i].ReadFromV2(mV2[i])
	}

	return m
}

// SetMembers sets strictly ordered list of
// storage group member objects.
//
// See also Members.
func (sg *StorageGroup) SetMembers(members []oid.ID) {
	mV2 := (*storagegroup.StorageGroup)(sg).GetMembers()

	if members == nil {
		mV2 = nil
	} else {
		ln := len(members)

		if cap(mV2) >= ln {
			mV2 = mV2[:0]
		} else {
			mV2 = make([]refs.ObjectID, 0, ln)
		}

		var oidV2 refs.ObjectID

		for i := 0; i < ln; i++ {
			members[i].WriteToV2(&oidV2)
			mV2 = append(mV2, oidV2)
		}
	}

	(*storagegroup.StorageGroup)(sg).SetMembers(mV2)
}

// Marshal marshals StorageGroup into a protobuf binary form.
//
// See also Unmarshal.
func (sg StorageGroup) Marshal() ([]byte, error) {
	return (*storagegroup.StorageGroup)(&sg).StableMarshal(nil), nil
}

// Unmarshal unmarshals protobuf binary representation of StorageGroup.
//
// See also Marshal.
func (sg *StorageGroup) Unmarshal(data []byte) error {
	v2 := (*storagegroup.StorageGroup)(sg)
	err := v2.Unmarshal(data)
	if err != nil {
		return err
	}

	return sg.readFromV2(*v2, false)
}

// MarshalJSON encodes StorageGroup to protobuf JSON format.
//
// See also UnmarshalJSON.
func (sg StorageGroup) MarshalJSON() ([]byte, error) {
	v2 := (storagegroup.StorageGroup)(sg)
	return v2.MarshalJSON()
}

// UnmarshalJSON decodes StorageGroup from protobuf JSON format.
//
// See also MarshalJSON.
func (sg *StorageGroup) UnmarshalJSON(data []byte) error {
	v2 := (*storagegroup.StorageGroup)(sg)
	err := v2.UnmarshalJSON(data)
	if err != nil {
		return err
	}

	return sg.readFromV2(*v2, false)
}

// ReadFromObject assemble StorageGroup from a regular
// Object structure. Object must contain unambiguous information
// about its expiration epoch, otherwise behaviour is undefined.
//
// Returns any error appeared during storage group parsing; returns
// error if object is not of TypeStorageGroup type.
func ReadFromObject(sg *StorageGroup, o objectSDK.Object) error {
	if typ := o.Type(); typ != objectSDK.TypeStorageGroup {
		return fmt.Errorf("object is not of StorageGroup type: %s", typ)
	}

	err := sg.Unmarshal(o.Payload())
	if err != nil {
		return fmt.Errorf("could not unmarshal object: %w", err)
	}

	var expObj uint64

	for _, attr := range o.Attributes() {
		if attr.Key() == objectV2.SysAttributeExpEpoch {
			expObj, err = strconv.ParseUint(attr.Value(), 10, 64)
			if err != nil {
				return fmt.Errorf("could not get expiration from object: %w", err)
			}

			break
		}
	}

	// Supporting deprecated functionality.
	// See https://github.com/nspcc-dev/neofs-api/pull/205.
	if expSG := sg.ExpirationEpoch(); expObj != expSG {
		return fmt.Errorf(
			"expiration does not match: from object: %d, from payload: %d",
			expObj, expSG)
	}

	return nil
}

// WriteToObject writes StorageGroup to a regular
// Object structure. Object must not contain ambiguous
// information about its expiration epoch or must not
// have it at all.
//
// Written information:
//   - expiration epoch;
//   - object type (TypeStorageGroup);
//   - raw payload.
func WriteToObject(sg StorageGroup, o *objectSDK.Object) {
	sgRaw, err := sg.Marshal()
	if err != nil {
		// Marshal() does not return errors
		// in the next API release
		panic(fmt.Errorf("could not marshal storage group: %w", err))
	}

	o.SetPayload(sgRaw)
	o.SetType(objectSDK.TypeStorageGroup)

	attrs := o.Attributes()
	var expAttrFound bool

	for i := range attrs {
		if attrs[i].Key() == objectV2.SysAttributeExpEpoch {
			expAttrFound = true
			attrs[i].SetValue(strconv.FormatUint(sg.ExpirationEpoch(), 10))

			break
		}
	}

	if !expAttrFound {
		var attr objectSDK.Attribute

		attr.SetKey(objectV2.SysAttributeExpEpoch)
		attr.SetValue(strconv.FormatUint(sg.ExpirationEpoch(), 10))

		attrs = append(attrs, attr)
	}

	o.SetAttributes(attrs...)
}
