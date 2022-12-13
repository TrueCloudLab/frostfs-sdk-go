package eacl

import (
	"bytes"
	"crypto/ecdsa"

	v2acl "github.com/TrueCloudLab/frostfs-api-go/v2/acl"
	"github.com/nspcc-dev/neo-go/pkg/crypto/keys"
)

// Target is a group of request senders to match ContainerEACL. Defined by role enum
// and set of public keys.
//
// Target is compatible with v2 acl.EACLRecord.Target message.
type Target struct {
	role Role
	keys [][]byte
}

func ecdsaKeysToPtrs(keys []ecdsa.PublicKey) []*ecdsa.PublicKey {
	keysPtr := make([]*ecdsa.PublicKey, len(keys))

	for i := range keys {
		keysPtr[i] = &keys[i]
	}

	return keysPtr
}

// BinaryKeys returns list of public keys to identify
// target subject in a binary format.
func (t *Target) BinaryKeys() [][]byte {
	return t.keys
}

// SetBinaryKeys sets list of binary public keys to identify
// target subject.
func (t *Target) SetBinaryKeys(keys [][]byte) {
	t.keys = keys
}

// SetTargetECDSAKeys converts ECDSA public keys to a binary
// format and stores them in Target.
func SetTargetECDSAKeys(t *Target, pubs ...*ecdsa.PublicKey) {
	binKeys := t.BinaryKeys()
	ln := len(pubs)

	if cap(binKeys) >= ln {
		binKeys = binKeys[:0]
	} else {
		binKeys = make([][]byte, 0, ln)
	}

	for i := 0; i < ln; i++ {
		binKeys = append(binKeys, (*keys.PublicKey)(pubs[i]).Bytes())
	}

	t.SetBinaryKeys(binKeys)
}

// TargetECDSAKeys interprets binary public keys of Target
// as ECDSA public keys. If any key has a different format,
// the corresponding element will be nil.
func TargetECDSAKeys(t *Target) []*ecdsa.PublicKey {
	binKeys := t.BinaryKeys()
	ln := len(binKeys)

	pubs := make([]*ecdsa.PublicKey, ln)

	for i := 0; i < ln; i++ {
		p := new(keys.PublicKey)
		if p.DecodeBytes(binKeys[i]) == nil {
			pubs[i] = (*ecdsa.PublicKey)(p)
		}
	}

	return pubs
}

// SetRole sets target subject's role class.
func (t *Target) SetRole(r Role) {
	t.role = r
}

// Role returns target subject's role class.
func (t Target) Role() Role {
	return t.role
}

// ToV2 converts Target to v2 acl.EACLRecord.Target message.
//
// Nil Target converts to nil.
func (t *Target) ToV2() *v2acl.Target {
	if t == nil {
		return nil
	}

	target := new(v2acl.Target)
	target.SetRole(t.role.ToV2())
	target.SetKeys(t.keys)

	return target
}

// NewTarget creates, initializes and returns blank Target instance.
//
// Defaults:
//   - role: RoleUnknown;
//   - keys: nil.
func NewTarget() *Target {
	return NewTargetFromV2(new(v2acl.Target))
}

// NewTargetFromV2 converts v2 acl.EACLRecord.Target message to Target.
func NewTargetFromV2(target *v2acl.Target) *Target {
	if target == nil {
		return new(Target)
	}

	return &Target{
		role: RoleFromV2(target.GetRole()),
		keys: target.GetKeys(),
	}
}

// Marshal marshals Target into a protobuf binary form.
func (t *Target) Marshal() ([]byte, error) {
	return t.ToV2().StableMarshal(nil), nil
}

// Unmarshal unmarshals protobuf binary representation of Target.
func (t *Target) Unmarshal(data []byte) error {
	fV2 := new(v2acl.Target)
	if err := fV2.Unmarshal(data); err != nil {
		return err
	}

	*t = *NewTargetFromV2(fV2)

	return nil
}

// MarshalJSON encodes Target to protobuf JSON format.
func (t *Target) MarshalJSON() ([]byte, error) {
	return t.ToV2().MarshalJSON()
}

// UnmarshalJSON decodes Target from protobuf JSON format.
func (t *Target) UnmarshalJSON(data []byte) error {
	tV2 := new(v2acl.Target)
	if err := tV2.UnmarshalJSON(data); err != nil {
		return err
	}

	*t = *NewTargetFromV2(tV2)

	return nil
}

// equalTargets compares Target with each other.
func equalTargets(t1, t2 Target) bool {
	if t1.Role() != t2.Role() {
		return false
	}

	keys1, keys2 := t1.BinaryKeys(), t2.BinaryKeys()

	if len(keys1) != len(keys2) {
		return false
	}

	for i := 0; i < len(keys1); i++ {
		if !bytes.Equal(keys1[i], keys2[i]) {
			return false
		}
	}

	return true
}
