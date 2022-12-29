package cid

import (
	"crypto/sha256"
	"fmt"

	"github.com/TrueCloudLab/frostfs-api-go/v2/refs"
	"github.com/mr-tron/base58"
)

// ID represents FrostFS container identifier.
//
// ID is mutually compatible with github.com/TrueCloudLab/frostfs-api-go/v2/refs.ContainerID
// message. See ReadFromV2 / WriteToV2 methods.
//
// Instances can be created using built-in var declaration.
//
// Note that direct typecast is not safe and may result in loss of compatibility:
//
//	_ = ID([32]byte) // not recommended
type ID [sha256.Size]byte

// ReadFromV2 reads ID from the refs.ContainerID message.
// Returns an error if the message is malformed according
// to the FrostFS API V2 protocol.
//
// See also WriteToV2.
func (id *ID) ReadFromV2(m refs.ContainerID) error {
	return id.Decode(m.GetValue())
}

// WriteToV2 writes ID to the refs.ContainerID message.
// The message must not be nil.
//
// See also ReadFromV2.
func (id ID) WriteToV2(m *refs.ContainerID) {
	m.SetValue(id[:])
}

// Encode encodes ID into 32 bytes of dst. Panics if
// dst length is less than 32.
//
// Zero ID is all zeros.
//
// See also Decode.
func (id ID) Encode(dst []byte) {
	if l := len(dst); l < sha256.Size {
		panic(fmt.Sprintf("destination length is less than %d bytes: %d", sha256.Size, l))
	}

	copy(dst, id[:])
}

// Decode decodes src bytes into ID.
//
// Decode expects that src has 32 bytes length. If the input is malformed,
// Decode returns an error describing format violation. In this case ID
// remains unchanged.
//
// Decode doesn't mutate src.
//
// See also Encode.
func (id *ID) Decode(src []byte) error {
	if len(src) != sha256.Size {
		return fmt.Errorf("invalid length %d", len(src))
	}

	copy(id[:], src)

	return nil
}

// SetSHA256 sets container identifier value to SHA256 checksum of container structure.
func (id *ID) SetSHA256(v [sha256.Size]byte) {
	copy(id[:], v[:])
}

// Equals defines a comparison relation between two ID instances.
//
// Note that comparison using '==' operator is not recommended since it MAY result
// in loss of compatibility.
func (id ID) Equals(id2 ID) bool {
	return id == id2
}

// EncodeToString encodes ID into FrostFS API protocol string.
//
// Zero ID is base58 encoding of 32 zeros.
//
// See also DecodeString.
func (id ID) EncodeToString() string {
	return base58.Encode(id[:])
}

// DecodeString decodes string into ID according to FrostFS API protocol. Returns
// an error if s is malformed.
//
// See also DecodeString.
func (id *ID) DecodeString(s string) error {
	data, err := base58.Decode(s)
	if err != nil {
		return fmt.Errorf("decode base58: %w", err)
	}

	return id.Decode(data)
}

// String implements fmt.Stringer.
//
// String is designed to be human-readable, and its format MAY differ between
// SDK versions. String MAY return same result as EncodeToString. String MUST NOT
// be used to encode ID into FrostFS protocol string.
func (id ID) String() string {
	return id.EncodeToString()
}
