package oid

import (
	"errors"
	"fmt"
	"strings"

	"github.com/TrueCloudLab/frostfs-api-go/v2/refs"
	cid "github.com/TrueCloudLab/frostfs-sdk-go/container/id"
)

// Address represents global object identifier in NeoFS network. Each object
// belongs to exactly one container and is uniquely addressed within the container.
//
// Address is mutually compatible with github.com/TrueCloudLab/frostfs-api-go/v2/refs.Address
// message. See ReadFromV2 / WriteToV2 methods.
//
// Instances can be created using built-in var declaration.
type Address struct {
	cnr cid.ID

	obj ID
}

// ReadFromV2 reads Address from the refs.Address message. Returns an error if
// the message is malformed according to the NeoFS API V2 protocol.
//
// See also WriteToV2.
func (x *Address) ReadFromV2(m refs.Address) error {
	cnr := m.GetContainerID()
	if cnr == nil {
		return errors.New("missing container ID")
	}

	obj := m.GetObjectID()
	if obj == nil {
		return errors.New("missing object ID")
	}

	err := x.cnr.ReadFromV2(*cnr)
	if err != nil {
		return fmt.Errorf("invalid container ID: %w", err)
	}

	err = x.obj.ReadFromV2(*obj)
	if err != nil {
		return fmt.Errorf("invalid object ID: %w", err)
	}

	return nil
}

// WriteToV2 writes Address to the refs.Address message.
// The message must not be nil.
//
// See also ReadFromV2.
func (x Address) WriteToV2(m *refs.Address) {
	var obj refs.ObjectID
	x.obj.WriteToV2(&obj)

	var cnr refs.ContainerID
	x.cnr.WriteToV2(&cnr)

	m.SetObjectID(&obj)
	m.SetContainerID(&cnr)
}

// MarshalJSON encodes Address into a JSON format of the NeoFS API protocol
// (Protocol Buffers JSON).
//
// See also UnmarshalJSON.
func (x Address) MarshalJSON() ([]byte, error) {
	var m refs.Address
	x.WriteToV2(&m)

	return m.MarshalJSON()
}

// UnmarshalJSON decodes NeoFS API protocol JSON format into the Address
// (Protocol Buffers JSON). Returns an error describing a format violation.
//
// See also MarshalJSON.
func (x *Address) UnmarshalJSON(data []byte) error {
	var m refs.Address

	err := m.UnmarshalJSON(data)
	if err != nil {
		return err
	}

	return x.ReadFromV2(m)
}

// Container returns unique identifier of the NeoFS object container.
//
// Zero Address has zero container ID, which is incorrect according to NeoFS
// API protocol.
//
// See also SetContainer.
func (x Address) Container() cid.ID {
	return x.cnr
}

// SetContainer sets unique identifier of the NeoFS object container.
//
// See also Container.
func (x *Address) SetContainer(id cid.ID) {
	x.cnr = id
}

// Object returns unique identifier of the object in the container
// identified by Container().
//
// Zero Address has zero object ID, which is incorrect according to NeoFS
// API protocol.
//
// See also SetObject.
func (x Address) Object() ID {
	return x.obj
}

// SetObject sets unique identifier of the object in the container
// identified by Container().
//
// See also Object.
func (x *Address) SetObject(id ID) {
	x.obj = id
}

// delimiter of container and object IDs in Address protocol string.
const idDelimiter = "/"

// EncodeToString encodes Address into NeoFS API protocol string: concatenation
// of the string-encoded container and object IDs delimited by a slash.
//
// See also DecodeString.
func (x Address) EncodeToString() string {
	return x.cnr.EncodeToString() + "/" + x.obj.EncodeToString()
}

// DecodeString decodes string into Address according to NeoFS API protocol. Returns
// an error if s is malformed.
//
// See also DecodeString.
func (x *Address) DecodeString(s string) error {
	indDelimiter := strings.Index(s, idDelimiter)
	if indDelimiter < 0 {
		return errors.New("missing delimiter")
	}

	err := x.cnr.DecodeString(s[:indDelimiter])
	if err != nil {
		return fmt.Errorf("decode container string: %w", err)
	}

	err = x.obj.DecodeString(s[indDelimiter+1:])
	if err != nil {
		return fmt.Errorf("decode object string: %w", err)
	}

	return nil
}

// String implements fmt.Stringer.
//
// String is designed to be human-readable, and its format MAY differ between
// SDK versions. String MAY return same result as EncodeToString. String MUST NOT
// be used to encode Address into NeoFS protocol string.
func (x Address) String() string {
	return x.EncodeToString()
}
