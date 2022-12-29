package bearer

import (
	"crypto/ecdsa"
	"errors"
	"fmt"

	"github.com/TrueCloudLab/frostfs-api-go/v2/acl"
	"github.com/TrueCloudLab/frostfs-api-go/v2/refs"
	cid "github.com/TrueCloudLab/frostfs-sdk-go/container/id"
	frostfscrypto "github.com/TrueCloudLab/frostfs-sdk-go/crypto"
	frostfsecdsa "github.com/TrueCloudLab/frostfs-sdk-go/crypto/ecdsa"
	"github.com/TrueCloudLab/frostfs-sdk-go/eacl"
	"github.com/TrueCloudLab/frostfs-sdk-go/user"
)

// Token represents bearer token for object service operations.
//
// Token is mutually compatible with github.com/TrueCloudLab/frostfs-api-go/v2/acl.BearerToken
// message. See ReadFromV2 / WriteToV2 methods.
//
// Instances can be created using built-in var declaration.
type Token struct {
	targetUserSet bool
	targetUser    user.ID

	eaclTableSet bool
	eaclTable    eacl.Table

	lifetimeSet   bool
	iat, nbf, exp uint64

	sigSet bool
	sig    refs.Signature
}

// reads Token from the acl.BearerToken message. If checkFieldPresence is set,
// returns an error on absence of any protocol-required field.
func (b *Token) readFromV2(m acl.BearerToken, checkFieldPresence bool) error {
	var err error

	body := m.GetBody()
	if checkFieldPresence && body == nil {
		return errors.New("missing token body")
	}

	eaclTable := body.GetEACL()
	if b.eaclTableSet = eaclTable != nil; b.eaclTableSet {
		b.eaclTable = *eacl.NewTableFromV2(eaclTable)
	} else if checkFieldPresence {
		return errors.New("missing eACL table")
	}

	targetUser := body.GetOwnerID()
	if b.targetUserSet = targetUser != nil; b.targetUserSet {
		err = b.targetUser.ReadFromV2(*targetUser)
		if err != nil {
			return fmt.Errorf("invalid target user: %w", err)
		}
	}

	lifetime := body.GetLifetime()
	if b.lifetimeSet = lifetime != nil; b.lifetimeSet {
		b.iat = lifetime.GetIat()
		b.nbf = lifetime.GetNbf()
		b.exp = lifetime.GetExp()
	} else if checkFieldPresence {
		return errors.New("missing token lifetime")
	}

	sig := m.GetSignature()
	if b.sigSet = sig != nil; sig != nil {
		b.sig = *sig
	} else if checkFieldPresence {
		return errors.New("missing body signature")
	}

	return nil
}

// ReadFromV2 reads Token from the acl.BearerToken message.
//
// See also WriteToV2.
func (b *Token) ReadFromV2(m acl.BearerToken) error {
	return b.readFromV2(m, true)
}

func (b Token) fillBody() *acl.BearerTokenBody {
	if !b.eaclTableSet && !b.targetUserSet && !b.lifetimeSet {
		return nil
	}

	var body acl.BearerTokenBody

	if b.eaclTableSet {
		body.SetEACL(b.eaclTable.ToV2())
	}

	if b.targetUserSet {
		var targetUser refs.OwnerID
		b.targetUser.WriteToV2(&targetUser)

		body.SetOwnerID(&targetUser)
	}

	if b.lifetimeSet {
		var lifetime acl.TokenLifetime
		lifetime.SetIat(b.iat)
		lifetime.SetNbf(b.nbf)
		lifetime.SetExp(b.exp)

		body.SetLifetime(&lifetime)
	}

	return &body
}

func (b Token) signedData() []byte {
	return b.fillBody().StableMarshal(nil)
}

// WriteToV2 writes Token to the acl.BearerToken message.
// The message must not be nil.
//
// See also ReadFromV2.
func (b Token) WriteToV2(m *acl.BearerToken) {
	m.SetBody(b.fillBody())

	var sig *refs.Signature

	if b.sigSet {
		sig = &b.sig
	}

	m.SetSignature(sig)
}

// SetExp sets "exp" (expiration time) claim which identifies the
// expiration time (in FrostFS epochs) after which the Token MUST NOT be
// accepted for processing. The processing of the "exp" claim requires
// that the current epoch MUST be before or equal to the expiration epoch
// listed in the "exp" claim.
//
// Naming is inspired by https://datatracker.ietf.org/doc/html/rfc7519#section-4.1.4.
//
// See also InvalidAt.
func (b *Token) SetExp(exp uint64) {
	b.exp = exp
	b.lifetimeSet = true
}

// SetNbf sets "nbf" (not before) claim which identifies the time (in
// FrostFS epochs) before which the Token MUST NOT be accepted for processing. The
// processing of the "nbf" claim requires that the current epoch MUST be
// after or equal to the not-before epoch listed in the "nbf" claim.
//
// Naming is inspired by https://datatracker.ietf.org/doc/html/rfc7519#section-4.1.5.
//
// See also InvalidAt.
func (b *Token) SetNbf(nbf uint64) {
	b.nbf = nbf
	b.lifetimeSet = true
}

// SetIat sets "iat" (issued at) claim which identifies the time (in FrostFS
// epochs) at which the Token was issued. This claim can be used to determine
// the age of the Token.
//
// Naming is inspired by https://datatracker.ietf.org/doc/html/rfc7519#section-4.1.6.
//
// See also InvalidAt.
func (b *Token) SetIat(iat uint64) {
	b.iat = iat
	b.lifetimeSet = true
}

// InvalidAt asserts "exp", "nbf" and "iat" claims for the given epoch.
//
// Zero Container is invalid in any epoch.
//
// See also SetExp, SetNbf, SetIat.
func (b Token) InvalidAt(epoch uint64) bool {
	return !b.lifetimeSet || b.nbf > epoch || b.iat > epoch || b.exp < epoch
}

// SetEACLTable sets eacl.Table that replaces the one from the issuer's
// container. If table has specified container, bearer token can be used only
// for operations within this specific container. Otherwise, Token can be used
// within any issuer's container.
//
// SetEACLTable MUST be called if Token is going to be transmitted over
// FrostFS API V2 protocol.
//
// See also EACLTable, AssertContainer.
func (b *Token) SetEACLTable(table eacl.Table) {
	b.eaclTable = table
	b.eaclTableSet = true
}

// EACLTable returns extended ACL table set by SetEACLTable.
//
// Zero Token has zero eacl.Table.
func (b Token) EACLTable() eacl.Table {
	if b.eaclTableSet {
		return b.eaclTable
	}

	return eacl.Table{}
}

// AssertContainer checks if the token is valid within the given container.
//
// Note: cnr is assumed to refer to the issuer's container, otherwise the check
// is meaningless.
//
// Zero Token is valid in any container.
//
// See also SetEACLTable.
func (b Token) AssertContainer(cnr cid.ID) bool {
	if !b.eaclTableSet {
		return true
	}

	cnrTable, set := b.eaclTable.CID()
	return !set || cnrTable.Equals(cnr)
}

// ForUser specifies ID of the user who can use the Token for the operations
// within issuer's container(s).
//
// Optional: by default, any user has access to Token usage.
//
// See also AssertUser.
func (b *Token) ForUser(id user.ID) {
	b.targetUser = id
	b.targetUserSet = true
}

// AssertUser checks if the Token is issued to the given user.
//
// Zero Token is available to any user.
//
// See also ForUser.
func (b Token) AssertUser(id user.ID) bool {
	return !b.targetUserSet || b.targetUser.Equals(id)
}

// Sign calculates and writes signature of the Token data using issuer's secret.
// Returns signature calculation errors.
//
// Sign MUST be called if Token is going to be transmitted over
// FrostFS API V2 protocol.
//
// Note that any Token mutation is likely to break the signature, so it is
// expected to be calculated as a final stage of Token formation.
//
// See also VerifySignature, Issuer.
func (b *Token) Sign(key ecdsa.PrivateKey) error {
	var sig frostfscrypto.Signature

	err := sig.Calculate(frostfsecdsa.Signer(key), b.signedData())
	if err != nil {
		return err
	}

	sig.WriteToV2(&b.sig)
	b.sigSet = true

	return nil
}

// VerifySignature checks if Token signature is presented and valid.
//
// Zero Token fails the check.
//
// See also Sign.
func (b Token) VerifySignature() bool {
	if !b.sigSet {
		return false
	}

	var sig frostfscrypto.Signature

	// TODO: (#233) check owner<->key relation
	return sig.ReadFromV2(b.sig) == nil && sig.Verify(b.signedData())
}

// Marshal encodes Token into a binary format of the FrostFS API protocol
// (Protocol Buffers V3 with direct field order).
//
// See also Unmarshal.
func (b Token) Marshal() []byte {
	var m acl.BearerToken
	b.WriteToV2(&m)

	return m.StableMarshal(nil)
}

// Unmarshal decodes FrostFS API protocol binary data into the Token
// (Protocol Buffers V3 with direct field order). Returns an error describing
// a format violation.
//
// See also Marshal.
func (b *Token) Unmarshal(data []byte) error {
	var m acl.BearerToken

	err := m.Unmarshal(data)
	if err != nil {
		return err
	}

	return b.readFromV2(m, false)
}

// MarshalJSON encodes Token into a JSON format of the FrostFS API protocol
// (Protocol Buffers V3 JSON).
//
// See also UnmarshalJSON.
func (b Token) MarshalJSON() ([]byte, error) {
	var m acl.BearerToken
	b.WriteToV2(&m)

	return m.MarshalJSON()
}

// UnmarshalJSON decodes FrostFS API protocol JSON data into the Token
// (Protocol Buffers V3 JSON). Returns an error describing a format violation.
//
// See also MarshalJSON.
func (b *Token) UnmarshalJSON(data []byte) error {
	var m acl.BearerToken

	err := m.UnmarshalJSON(data)
	if err != nil {
		return err
	}

	return b.readFromV2(m, false)
}

// SigningKeyBytes returns issuer's public key in a binary format of
// FrostFS API protocol.
//
// Unsigned Token has empty key.
//
// See also ResolveIssuer.
func (b Token) SigningKeyBytes() []byte {
	if b.sigSet {
		return b.sig.GetKey()
	}

	return nil
}

// ResolveIssuer resolves issuer's user.ID from the key used for Token signing.
// Returns zero user.ID if Token is unsigned or key has incorrect format.
//
// See also SigningKeyBytes.
func ResolveIssuer(b Token) (usr user.ID) {
	binKey := b.SigningKeyBytes()

	if len(binKey) != 0 {
		var key frostfsecdsa.PublicKey
		if key.Decode(binKey) == nil {
			user.IDFromKey(&usr, ecdsa.PublicKey(key))
		}
	}

	return
}
