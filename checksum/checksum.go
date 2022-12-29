package checksum

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/TrueCloudLab/frostfs-api-go/v2/refs"
	"github.com/TrueCloudLab/tzhash/tz"
)

// Checksum represents checksum of some digital data.
//
// Checksum is mutually compatible with github.com/TrueCloudLab/frostfs-api-go/v2/refs.Checksum
// message. See ReadFromV2 / WriteToV2 methods.
//
// Instances can be created using built-in var declaration.
//
// Note that direct typecast is not safe and may result in loss of compatibility:
//
//	_ = Checksum(refs.Checksum{}) // not recommended
type Checksum refs.Checksum

// Type represents the enumeration
// of checksum types.
type Type uint8

const (
	// Unknown is an undefined checksum type.
	Unknown Type = iota

	// SHA256 is a SHA256 checksum type.
	SHA256

	// TZ is a Tillich-Zémor checksum type.
	TZ
)

// ReadFromV2 reads Checksum from the refs.Checksum message. Checks if the
// message conforms to FrostFS API V2 protocol.
//
// See also WriteToV2.
func (c *Checksum) ReadFromV2(m refs.Checksum) error {
	if len(m.GetSum()) == 0 {
		return errors.New("missing value")
	}

	switch m.GetType() {
	default:
		return fmt.Errorf("unsupported type %v", m.GetType())
	case refs.SHA256, refs.TillichZemor:
	}

	*c = Checksum(m)

	return nil
}

// WriteToV2 writes Checksum to the refs.Checksum message.
// The message must not be nil.
//
// See also ReadFromV2.
func (c Checksum) WriteToV2(m *refs.Checksum) {
	*m = (refs.Checksum)(c)
}

// Type returns checksum type.
//
// Zero Checksum has Unknown checksum type.
//
// See also SetTillichZemor and SetSHA256.
func (c Checksum) Type() Type {
	v2 := (refs.Checksum)(c)
	switch v2.GetType() {
	case refs.SHA256:
		return SHA256
	case refs.TillichZemor:
		return TZ
	default:
		return Unknown
	}
}

// Value returns checksum bytes. Return value
// MUST NOT be mutated.
//
// Zero Checksum has nil sum.
//
// See also SetTillichZemor and SetSHA256.
func (c Checksum) Value() []byte {
	v2 := (refs.Checksum)(c)
	return v2.GetSum()
}

// SetSHA256 sets checksum to SHA256 hash.
//
// See also Calculate.
func (c *Checksum) SetSHA256(v [sha256.Size]byte) {
	v2 := (*refs.Checksum)(c)

	v2.SetType(refs.SHA256)
	v2.SetSum(v[:])
}

// Calculate calculates checksum and sets it
// to the passed checksum. Checksum must not be nil.
//
// Does nothing if the passed type is not one of the:
//   - SHA256;
//   - TZ.
//
// Does not mutate the passed value.
//
// See also SetSHA256, SetTillichZemor.
func Calculate(c *Checksum, t Type, v []byte) {
	switch t {
	case SHA256:
		c.SetSHA256(sha256.Sum256(v))
	case TZ:
		c.SetTillichZemor(tz.Sum(v))
	default:
	}
}

// SetTillichZemor sets checksum to Tillich-Zémor hash.
//
// See also Calculate.
func (c *Checksum) SetTillichZemor(v [tz.Size]byte) {
	v2 := (*refs.Checksum)(c)

	v2.SetType(refs.TillichZemor)
	v2.SetSum(v[:])
}

// String implements fmt.Stringer.
//
// String is designed to be human-readable, and its format MAY differ between
// SDK versions.
func (c Checksum) String() string {
	v2 := (refs.Checksum)(c)
	return fmt.Sprintf("%s:%s", c.Type(), hex.EncodeToString(v2.GetSum()))
}

// String implements fmt.Stringer.
//
// String is designed to be human-readable, and its format MAY differ between
// SDK versions.
func (m Type) String() string {
	var m2 refs.ChecksumType

	switch m {
	default:
		m2 = refs.UnknownChecksum
	case TZ:
		m2 = refs.TillichZemor
	case SHA256:
		m2 = refs.SHA256
	}

	return m2.String()
}
