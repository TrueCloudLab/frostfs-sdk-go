package storagegroup

import (
	"github.com/nspcc-dev/neofs-api-go/v2/refs"
	"github.com/nspcc-dev/neofs-api-go/v2/storagegroup"
	"github.com/nspcc-dev/neofs-sdk-go/checksum"
	oid "github.com/nspcc-dev/neofs-sdk-go/object/id"
)

// StorageGroup represents v2-compatible storage group.
type StorageGroup storagegroup.StorageGroup

// NewFromV2 wraps v2 StorageGroup message to StorageGroup.
//
// Nil storagegroup.StorageGroup converts to nil.
func NewFromV2(aV2 *storagegroup.StorageGroup) *StorageGroup {
	return (*StorageGroup)(aV2)
}

// New creates and initializes blank StorageGroup.
//
// Defaults:
//  - size: 0;
//  - exp: 0;
//  - members: nil;
//  - hash: nil.
func New() *StorageGroup {
	return NewFromV2(new(storagegroup.StorageGroup))
}

// ValidationDataSize returns total size of the payloads
// of objects in the storage group.
func (sg *StorageGroup) ValidationDataSize() uint64 {
	return (*storagegroup.StorageGroup)(sg).GetValidationDataSize()
}

// SetValidationDataSize sets total size of the payloads
// of objects in the storage group.
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
		v.ReadFromV2(*checksumV2)
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

// ExpirationEpoch returns last NeoFS epoch number
// of the storage group lifetime.
func (sg *StorageGroup) ExpirationEpoch() uint64 {
	return (*storagegroup.StorageGroup)(sg).GetExpirationEpoch()
}

// SetExpirationEpoch sets last NeoFS epoch number
// of the storage group lifetime.
func (sg *StorageGroup) SetExpirationEpoch(epoch uint64) {
	(*storagegroup.StorageGroup)(sg).SetExpirationEpoch(epoch)
}

// Members returns strictly ordered list of
// storage group member objects.
func (sg *StorageGroup) Members() []oid.ID {
	mV2 := (*storagegroup.StorageGroup)(sg).GetMembers()

	if mV2 == nil {
		return nil
	}

	m := make([]oid.ID, len(mV2))

	for i := range mV2 {
		m[i] = *oid.NewIDFromV2(&mV2[i])
	}

	return m
}

// SetMembers sets strictly ordered list of
// storage group member objects.
func (sg *StorageGroup) SetMembers(members []oid.ID) {
	mV2 := (*storagegroup.StorageGroup)(sg).GetMembers()

	if members == nil {
		mV2 = nil
	} else {
		ln := len(members)

		if cap(mV2) >= ln {
			mV2 = mV2[:0]
		} else {
			mV2 = make([]refs.ObjectID, ln)
		}

		for i := 0; i < ln; i++ {
			mV2[i] = *members[i].ToV2()
		}
	}

	(*storagegroup.StorageGroup)(sg).SetMembers(mV2)
}

// ToV2 converts StorageGroup to v2 StorageGroup message.
//
// Nil StorageGroup converts to nil.
func (sg *StorageGroup) ToV2() *storagegroup.StorageGroup {
	return (*storagegroup.StorageGroup)(sg)
}

// Marshal marshals StorageGroup into a protobuf binary form.
func (sg *StorageGroup) Marshal() ([]byte, error) {
	return (*storagegroup.StorageGroup)(sg).StableMarshal(nil)
}

// Unmarshal unmarshals protobuf binary representation of StorageGroup.
func (sg *StorageGroup) Unmarshal(data []byte) error {
	return (*storagegroup.StorageGroup)(sg).Unmarshal(data)
}

// MarshalJSON encodes StorageGroup to protobuf JSON format.
func (sg *StorageGroup) MarshalJSON() ([]byte, error) {
	return (*storagegroup.StorageGroup)(sg).MarshalJSON()
}

// UnmarshalJSON decodes StorageGroup from protobuf JSON format.
func (sg *StorageGroup) UnmarshalJSON(data []byte) error {
	return (*storagegroup.StorageGroup)(sg).UnmarshalJSON(data)
}
