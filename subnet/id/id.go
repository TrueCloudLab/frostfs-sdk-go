package subnetid

import (
	"fmt"
	"strconv"

	"github.com/TrueCloudLab/frostfs-api-go/v2/refs"
)

// ID represents unique identifier of the subnet in the NeoFS network.
//
// ID is mutually compatible with github.com/TrueCloudLab/frostfs-api-go/v2/refs.SubnetID
// message. See ReadFromV2 / WriteToV2 methods.
//
// Instances can be created using built-in var declaration. Zero value is
// equivalent to identifier of the zero subnet (whole NeoFS network).
type ID struct {
	m refs.SubnetID
}

// ReadFromV2 reads ID from the refs.SubnetID message. Checks if the
// message conforms to NeoFS API V2 protocol.
//
// See also WriteToV2.
func (x *ID) ReadFromV2(msg refs.SubnetID) error {
	x.m = msg
	return nil
}

// WriteToV2 writes ID to refs.SubnetID message structure. The message MUST NOT
// be nil.
//
// See also ReadFromV2.
func (x ID) WriteToV2(msg *refs.SubnetID) {
	*msg = x.m
}

// Equals defines a comparison relation between two ID instances.
//
// Note that comparison using '==' operator is not recommended since it MAY result
// in loss of compatibility.
func (x ID) Equals(x2 ID) bool {
	return x.m.GetValue() == x2.m.GetValue()
}

// EncodeToString encodes ID into NeoFS API protocol string (base10 encoding).
//
// See also DecodeString.
func (x ID) EncodeToString() string {
	return strconv.FormatUint(uint64(x.m.GetValue()), 10)
}

// DecodeString decodes string calculated using EncodeToString. Returns
// an error if s is malformed.
func (x *ID) DecodeString(s string) error {
	num, err := strconv.ParseUint(s, 10, 32)
	if err != nil {
		return fmt.Errorf("invalid numeric value: %w", err)
	}

	x.m.SetValue(uint32(num))

	return nil
}

// String implements fmt.Stringer.
//
// String is designed to be human-readable, and its format MAY differ between
// SDK versions. String MAY return same result as EncodeToString. String MUST NOT
// be used to encode ID into NeoFS protocol string.
func (x ID) String() string {
	return "#" + strconv.FormatUint(uint64(x.m.GetValue()), 10)
}

// Marshal encodes ID into a binary format of the NeoFS API protocol
// (Protocol Buffers with direct field order).
//
// See also Unmarshal.
func (x ID) Marshal() []byte {
	return x.m.StableMarshal(nil)
}

// Unmarshal decodes binary ID calculated using Marshal. Returns an error
// describing a format violation.
func (x *ID) Unmarshal(data []byte) error {
	return x.m.Unmarshal(data)
}

// SetNumeric sets ID value in numeric format. By default, number is 0 which
// refers to the zero subnet.
func (x *ID) SetNumeric(num uint32) {
	x.m.SetValue(num)
}

// IsZero compares id with zero subnet ID.
func IsZero(id ID) bool {
	return id.Equals(ID{})
}

// MakeZero makes ID to refer to zero subnet.
//
// Makes no sense to call on zero value (e.g. declared using var).
func MakeZero(id *ID) {
	id.SetNumeric(0)
}
