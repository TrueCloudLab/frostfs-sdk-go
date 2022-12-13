package object

import (
	v2object "github.com/TrueCloudLab/frostfs-api-go/v2/object"
	"github.com/TrueCloudLab/frostfs-api-go/v2/refs"
	oid "github.com/TrueCloudLab/frostfs-sdk-go/object/id"
)

// Lock represents record with locked objects. It is compatible with
// NeoFS API V2 protocol.
//
// Lock instance can be written to the Object, see WriteLock/ReadLock.
type Lock v2object.Lock

// WriteLock writes Lock to the Object, and sets its type to TypeLock.
// The object must not be nil.
//
// See also ReadLock.
func WriteLock(obj *Object, l Lock) {
	obj.SetType(TypeLock)
	obj.SetPayload(l.Marshal())
}

// ReadLock reads Lock from the Object. The lock must not be nil.
// Returns an error describing incorrect format. Makes sense only
// if object has TypeLock type.
//
// See also WriteLock.
func ReadLock(l *Lock, obj Object) error {
	return l.Unmarshal(obj.Payload())
}

// NumberOfMembers returns number of members in lock list.
func (x Lock) NumberOfMembers() int {
	return (*v2object.Lock)(&x).NumberOfMembers()
}

// ReadMembers reads list of locked members.
//
// Buffer length must not be less than NumberOfMembers.
func (x Lock) ReadMembers(buf []oid.ID) {
	var i int

	(*v2object.Lock)(&x).IterateMembers(func(idV2 refs.ObjectID) {
		_ = buf[i].ReadFromV2(idV2)
		i++
	})
}

// WriteMembers writes list of locked members.
func (x *Lock) WriteMembers(ids []oid.ID) {
	var members []refs.ObjectID

	if ids != nil {
		members = make([]refs.ObjectID, len(ids))

		for i := range ids {
			ids[i].WriteToV2(&members[i])
		}
	}

	(*v2object.Lock)(x).SetMembers(members)
}

// Marshal encodes the Lock into a NeoFS protocol binary format.
func (x Lock) Marshal() []byte {
	return (*v2object.Lock)(&x).StableMarshal(nil)
}

// Unmarshal decodes the Lock from its NeoFS protocol binary representation.
func (x *Lock) Unmarshal(data []byte) error {
	return (*v2object.Lock)(x).Unmarshal(data)
}
