package reputation

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/TrueCloudLab/frostfs-api-go/v2/reputation"
	"github.com/mr-tron/base58"
)

// PeerID represents unique identifier of the peer participating in the FrostFS
// reputation system.
//
// ID is mutually compatible with github.com/TrueCloudLab/frostfs-api-go/v2/reputation.PeerID
// message. See ReadFromV2 / WriteToV2 methods.
//
// Instances can be created using built-in var declaration.
type PeerID struct {
	m reputation.PeerID
}

// ReadFromV2 reads PeerID from the reputation.PeerID message. Returns an
// error if the message is malformed according to the FrostFS API V2 protocol.
//
// See also WriteToV2.
func (x *PeerID) ReadFromV2(m reputation.PeerID) error {
	val := m.GetPublicKey()
	if len(val) == 0 {
		return errors.New("missing ID bytes")
	}

	x.m = m

	return nil
}

// WriteToV2 writes PeerID to the reputation.PeerID message.
// The message must not be nil.
//
// See also ReadFromV2.
func (x PeerID) WriteToV2(m *reputation.PeerID) {
	*m = x.m
}

// SetPublicKey sets PeerID as a binary-encoded public key which authenticates
// the participant of the FrostFS reputation system.
//
// Argument MUST NOT be mutated, make a copy first.
//
// See also CompareKey.
func (x *PeerID) SetPublicKey(key []byte) {
	x.m.SetPublicKey(key)
}

// PublicKey return public key set using SetPublicKey.
//
// Zero PeerID has zero key which is incorrect according to FrostFS API
// protocol.
//
// Return value MUST NOT be mutated, make a copy first.
func (x PeerID) PublicKey() []byte {
	return x.m.GetPublicKey()
}

// ComparePeerKey checks if the given PeerID corresponds to the party
// authenticated by the given binary public key.
func ComparePeerKey(peer PeerID, key []byte) bool {
	return bytes.Equal(peer.PublicKey(), key)
}

// EncodeToString encodes ID into FrostFS API protocol string.
//
// Zero PeerID is base58 encoding of PeerIDSize zeros.
//
// See also DecodeString.
func (x PeerID) EncodeToString() string {
	return base58.Encode(x.m.GetPublicKey())
}

// DecodeString decodes string into PeerID according to FrostFS API protocol.
// Returns an error if s is malformed.
//
// See also DecodeString.
func (x *PeerID) DecodeString(s string) error {
	data, err := base58.Decode(s)
	if err != nil {
		return fmt.Errorf("decode base58: %w", err)
	}

	x.m.SetPublicKey(data)

	return nil
}

// String implements fmt.Stringer.
//
// String is designed to be human-readable, and its format MAY differ between
// SDK versions. String MAY return same result as EncodeToString. String MUST NOT
// be used to encode ID into FrostFS protocol string.
func (x PeerID) String() string {
	return x.EncodeToString()
}
