package session

import (
	"crypto/ecdsa"
	"errors"
	"fmt"

	"github.com/TrueCloudLab/frostfs-api-go/v2/refs"
	"github.com/TrueCloudLab/frostfs-api-go/v2/session"
	cid "github.com/TrueCloudLab/frostfs-sdk-go/container/id"
	frostfscrypto "github.com/TrueCloudLab/frostfs-sdk-go/crypto"
	"github.com/TrueCloudLab/frostfs-sdk-go/user"
)

// Container represents token of the NeoFS Container session. A session is opened
// between any two sides of the system, and implements a mechanism for transferring
// the power of attorney of actions to another network member. The session has a
// limited validity period, and applies to a strictly defined set of operations.
// See methods for details.
//
// Container is mutually compatible with github.com/TrueCloudLab/frostfs-api-go/v2/session.Token
// message. See ReadFromV2 / WriteToV2 methods.
//
// Instances can be created using built-in var declaration.
type Container struct {
	commonData

	verb ContainerVerb

	cnrSet bool
	cnr    cid.ID
}

// readContext is a contextReader needed for commonData methods.
func (x *Container) readContext(c session.TokenContext, checkFieldPresence bool) error {
	cCnr, ok := c.(*session.ContainerSessionContext)
	if !ok || cCnr == nil {
		return fmt.Errorf("invalid context %T", c)
	}

	x.cnrSet = !cCnr.Wildcard()
	cnr := cCnr.ContainerID()

	if x.cnrSet {
		if cnr != nil {
			err := x.cnr.ReadFromV2(*cnr)
			if err != nil {
				return fmt.Errorf("invalid container ID: %w", err)
			}
		} else if checkFieldPresence {
			return errors.New("missing container or wildcard flag")
		}
	} else if cnr != nil {
		return errors.New("container conflicts with wildcard flag")
	}

	x.verb = ContainerVerb(cCnr.Verb())

	return nil
}

func (x *Container) readFromV2(m session.Token, checkFieldPresence bool) error {
	return x.commonData.readFromV2(m, checkFieldPresence, x.readContext)
}

// ReadFromV2 reads Container from the session.Token message. Checks if the
// message conforms to NeoFS API V2 protocol.
//
// See also WriteToV2.
func (x *Container) ReadFromV2(m session.Token) error {
	return x.readFromV2(m, true)
}

func (x Container) writeContext() session.TokenContext {
	var c session.ContainerSessionContext
	c.SetWildcard(!x.cnrSet)
	c.SetVerb(session.ContainerSessionVerb(x.verb))

	if x.cnrSet {
		var cnr refs.ContainerID
		x.cnr.WriteToV2(&cnr)

		c.SetContainerID(&cnr)
	}

	return &c
}

// WriteToV2 writes Container to the session.Token message.
// The message must not be nil.
//
// See also ReadFromV2.
func (x Container) WriteToV2(m *session.Token) {
	x.writeToV2(m, x.writeContext)
}

// Marshal encodes Container into a binary format of the NeoFS API protocol
// (Protocol Buffers with direct field order).
//
// See also Unmarshal.
func (x Container) Marshal() []byte {
	return x.marshal(x.writeContext)
}

// Unmarshal decodes NeoFS API protocol binary format into the Container
// (Protocol Buffers with direct field order). Returns an error describing
// a format violation.
//
// See also Marshal.
func (x *Container) Unmarshal(data []byte) error {
	return x.unmarshal(data, x.readContext)
}

// MarshalJSON encodes Container into a JSON format of the NeoFS API protocol
// (Protocol Buffers JSON).
//
// See also UnmarshalJSON.
func (x Container) MarshalJSON() ([]byte, error) {
	return x.marshalJSON(x.writeContext)
}

// UnmarshalJSON decodes NeoFS API protocol JSON format into the Container
// (Protocol Buffers JSON). Returns an error describing a format violation.
//
// See also MarshalJSON.
func (x *Container) UnmarshalJSON(data []byte) error {
	return x.unmarshalJSON(data, x.readContext)
}

// Sign calculates and writes signature of the Container data.
// Returns signature calculation errors.
//
// Zero Container is unsigned.
//
// Note that any Container mutation is likely to break the signature, so it is
// expected to be calculated as a final stage of Container formation.
//
// See also VerifySignature.
func (x *Container) Sign(key ecdsa.PrivateKey) error {
	return x.sign(key, x.writeContext)
}

// VerifySignature checks if Container signature is presented and valid.
//
// Zero Container fails the check.
//
// See also Sign.
func (x Container) VerifySignature() bool {
	return x.verifySignature(x.writeContext)
}

// ApplyOnlyTo limits session scope to a given author container.
//
// See also AppliedTo.
func (x *Container) ApplyOnlyTo(cnr cid.ID) {
	x.cnr = cnr
	x.cnrSet = true
}

// AppliedTo checks if the session is propagated to the given container.
//
// Zero Container is applied to all author's containers.
//
// See also ApplyOnlyTo.
func (x Container) AppliedTo(cnr cid.ID) bool {
	return !x.cnrSet || x.cnr.Equals(cnr)
}

// ContainerVerb enumerates container operations.
type ContainerVerb int8

const (
	_ ContainerVerb = iota

	VerbContainerPut     // Put rpc
	VerbContainerDelete  // Delete rpc
	VerbContainerSetEACL // SetExtendedACL rpc
)

// ForVerb specifies the container operation of the session scope. Each
// Container is related to the single operation.
//
// See also AssertVerb.
func (x *Container) ForVerb(verb ContainerVerb) {
	x.verb = verb
}

// AssertVerb checks if Container relates to the given container operation.
//
// Zero Container relates to zero (unspecified) verb.
//
// See also ForVerb.
func (x Container) AssertVerb(verb ContainerVerb) bool {
	return x.verb == verb
}

// IssuedBy checks if Container session is issued by the given user.
//
// See also Container.Issuer.
func IssuedBy(cnr Container, id user.ID) bool {
	return cnr.Issuer().Equals(id)
}

// VerifySessionDataSignature verifies signature of the session data. In practice,
// the method is used to authenticate an operation with session data.
func (x Container) VerifySessionDataSignature(data, signature []byte) bool {
	var sigV2 refs.Signature
	sigV2.SetKey(x.authKey)
	sigV2.SetScheme(refs.ECDSA_RFC6979_SHA256)
	sigV2.SetSign(signature)

	var sig frostfscrypto.Signature

	return sig.ReadFromV2(sigV2) == nil && sig.Verify(data)
}
