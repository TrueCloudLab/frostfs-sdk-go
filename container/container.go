package container

import (
	"crypto/ecdsa"
	"crypto/sha256"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/TrueCloudLab/frostfs-api-go/v2/container"
	v2netmap "github.com/TrueCloudLab/frostfs-api-go/v2/netmap"
	"github.com/TrueCloudLab/frostfs-api-go/v2/refs"
	"github.com/TrueCloudLab/frostfs-sdk-go/container/acl"
	cid "github.com/TrueCloudLab/frostfs-sdk-go/container/id"
	frostfscrypto "github.com/TrueCloudLab/frostfs-sdk-go/crypto"
	frostfsecdsa "github.com/TrueCloudLab/frostfs-sdk-go/crypto/ecdsa"
	"github.com/TrueCloudLab/frostfs-sdk-go/netmap"
	subnetid "github.com/TrueCloudLab/frostfs-sdk-go/subnet/id"
	"github.com/TrueCloudLab/frostfs-sdk-go/user"
	"github.com/TrueCloudLab/frostfs-sdk-go/version"
	"github.com/google/uuid"
)

// Container represents descriptor of the NeoFS container. Container logically
// stores NeoFS objects. Container is one of the basic and at the same time
// necessary data storage units in the NeoFS. Container includes data about the
// owner, rules for placing objects and other information necessary for the
// system functioning.
//
// Container type instances can represent different container states in the
// system, depending on the context. To create new container in NeoFS zero
// instance SHOULD be declared, initialized using Init method and filled using
// dedicated methods. Once container is saved in the NeoFS network, it can't be
// changed: containers stored in the system are immutable, and NeoFS is a CAS
// of containers that are identified by a fixed length value (see cid.ID type).
// Instances for existing containers can be initialized using decoding methods
// (e.g Unmarshal).
//
// Container is mutually compatible with github.com/TrueCloudLab/frostfs-api-go/v2/container.Container
// message. See ReadFromV2 / WriteToV2 methods.
type Container struct {
	v2 container.Container
}

const (
	attributeName      = "Name"
	attributeTimestamp = "Timestamp"
)

// reads Container from the container.Container message. If checkFieldPresence is set,
// returns an error on absence of any protocol-required field.
func (x *Container) readFromV2(m container.Container, checkFieldPresence bool) error {
	var err error

	ownerV2 := m.GetOwnerID()
	if ownerV2 != nil {
		var owner user.ID

		err = owner.ReadFromV2(*ownerV2)
		if err != nil {
			return fmt.Errorf("invalid owner: %w", err)
		}
	} else if checkFieldPresence {
		return errors.New("missing owner")
	}

	binNonce := m.GetNonce()
	if len(binNonce) > 0 {
		var nonce uuid.UUID

		err = nonce.UnmarshalBinary(binNonce)
		if err != nil {
			return fmt.Errorf("invalid nonce: %w", err)
		} else if ver := nonce.Version(); ver != 4 {
			return fmt.Errorf("invalid nonce UUID version %d", ver)
		}
	} else if checkFieldPresence {
		return errors.New("missing nonce")
	}

	ver := m.GetVersion()
	if checkFieldPresence && ver == nil {
		return errors.New("missing version")
	}

	policyV2 := m.GetPlacementPolicy()
	if policyV2 != nil {
		var policy netmap.PlacementPolicy

		err = policy.ReadFromV2(*policyV2)
		if err != nil {
			return fmt.Errorf("invalid placement policy: %w", err)
		}
	} else if checkFieldPresence {
		return errors.New("missing placement policy")
	}

	attrs := m.GetAttributes()
	mAttr := make(map[string]struct{}, len(attrs))
	var key, val string
	var was bool

	for i := range attrs {
		key = attrs[i].GetKey()
		if key == "" {
			return errors.New("empty attribute key")
		}

		_, was = mAttr[key]
		if was {
			return fmt.Errorf("duplicated attribute %s", key)
		}

		val = attrs[i].GetValue()
		if val == "" {
			return fmt.Errorf("empty attribute value %s", key)
		}

		switch key {
		case container.SysAttributeSubnet:
			err = new(subnetid.ID).DecodeString(val)
		case attributeTimestamp:
			_, err = strconv.ParseInt(val, 10, 64)
		}

		if err != nil {
			return fmt.Errorf("invalid attribute value %s: %s (%w)", key, val, err)
		}

		mAttr[key] = struct{}{}
	}

	x.v2 = m

	return nil
}

// ReadFromV2 reads Container from the container.Container message. Checks if the
// message conforms to NeoFS API V2 protocol.
//
// See also WriteToV2.
func (x *Container) ReadFromV2(m container.Container) error {
	return x.readFromV2(m, true)
}

// WriteToV2 writes Container into the container.Container message.
// The message MUST NOT be nil.
//
// See also ReadFromV2.
func (x Container) WriteToV2(m *container.Container) {
	*m = x.v2
}

// Marshal encodes Container into a binary format of the NeoFS API protocol
// (Protocol Buffers with direct field order).
//
// See also Unmarshal.
func (x Container) Marshal() []byte {
	return x.v2.StableMarshal(nil)
}

// Unmarshal decodes NeoFS API protocol binary format into the Container
// (Protocol Buffers with direct field order). Returns an error describing
// a format violation.
//
// See also Marshal.
func (x *Container) Unmarshal(data []byte) error {
	var m container.Container

	err := m.Unmarshal(data)
	if err != nil {
		return err
	}

	return x.readFromV2(m, false)
}

// MarshalJSON encodes Container into a JSON format of the NeoFS API protocol
// (Protocol Buffers JSON).
//
// See also UnmarshalJSON.
func (x Container) MarshalJSON() ([]byte, error) {
	return x.v2.MarshalJSON()
}

// UnmarshalJSON decodes NeoFS API protocol JSON format into the Container
// (Protocol Buffers JSON). Returns an error describing a format violation.
//
// See also MarshalJSON.
func (x *Container) UnmarshalJSON(data []byte) error {
	return x.v2.UnmarshalJSON(data)
}

// Init initializes all internal data of the Container required by NeoFS API
// protocol. Init MUST be called when creating a new container. Init SHOULD NOT
// be called multiple times. Init SHOULD NOT be called if the Container instance
// is used for decoding only.
func (x *Container) Init() {
	var ver refs.Version
	version.Current().WriteToV2(&ver)

	x.v2.SetVersion(&ver)

	nonce, err := uuid.New().MarshalBinary()
	if err != nil {
		panic(fmt.Sprintf("unexpected error from UUID.MarshalBinary: %v", err))
	}

	x.v2.SetNonce(nonce)
}

// SetOwner specifies the owner of the Container. Each Container has exactly
// one owner, so SetOwner MUST be called for instances to be saved in the
// NeoFS.
//
// See also Owner.
func (x *Container) SetOwner(owner user.ID) {
	var m refs.OwnerID
	owner.WriteToV2(&m)

	x.v2.SetOwnerID(&m)
}

// Owner returns owner of the Container set using SetOwner.
//
// Zero Container has no owner which is incorrect according to NeoFS API
// protocol.
func (x Container) Owner() (res user.ID) {
	m := x.v2.GetOwnerID()
	if m != nil {
		err := res.ReadFromV2(*m)
		if err != nil {
			panic(fmt.Sprintf("unexpected error from user.ID.ReadFromV2: %v", err))
		}
	}

	return
}

// SetBasicACL specifies basic part of the Container ACL. Basic ACL is used
// to control access inside container storage.
//
// See also BasicACL.
func (x *Container) SetBasicACL(basicACL acl.Basic) {
	x.v2.SetBasicACL(basicACL.Bits())
}

// BasicACL returns basic ACL set using SetBasicACL.
//
// Zero Container has zero basic ACL which structurally correct but doesn't
// make sense since it denies any access to any party.
func (x Container) BasicACL() (res acl.Basic) {
	res.FromBits(x.v2.GetBasicACL())
	return
}

// SetPlacementPolicy sets placement policy for the objects within the Container.
// NeoFS storage layer strives to follow the specified policy.
//
// See also PlacementPolicy.
func (x *Container) SetPlacementPolicy(policy netmap.PlacementPolicy) {
	var m v2netmap.PlacementPolicy
	policy.WriteToV2(&m)

	x.v2.SetPlacementPolicy(&m)
}

// PlacementPolicy returns placement policy set using SetPlacementPolicy.
//
// Zero Container has no placement policy which is incorrect according to
// NeoFS API protocol.
func (x Container) PlacementPolicy() (res netmap.PlacementPolicy) {
	m := x.v2.GetPlacementPolicy()
	if m != nil {
		err := res.ReadFromV2(*m)
		if err != nil {
			panic(fmt.Sprintf("unexpected error from PlacementPolicy.ReadFromV2: %v", err))
		}
	}

	return
}

// SetAttribute sets Container attribute value by key. Both key and value
// MUST NOT be empty. Attributes set by the creator (owner) are most commonly
// ignored by the NeoFS system and used for application layer. Some attributes
// are so-called system or well-known attributes: they are reserved for system
// needs. System attributes SHOULD NOT be modified using SetAttribute, use
// corresponding methods/functions. List of the reserved keys is documented
// in the particular protocol version.
//
// SetAttribute overwrites existing attribute value.
//
// See also Attribute, IterateAttributes.
func (x *Container) SetAttribute(key, value string) {
	if key == "" {
		panic("empty attribute key")
	} else if value == "" {
		panic("empty attribute value")
	}

	attrs := x.v2.GetAttributes()
	ln := len(attrs)

	for i := 0; i < ln; i++ {
		if attrs[i].GetKey() == key {
			attrs[i].SetValue(value)
			return
		}
	}

	attrs = append(attrs, container.Attribute{})
	attrs[ln].SetKey(key)
	attrs[ln].SetValue(value)

	x.v2.SetAttributes(attrs)
}

// Attribute reads value of the Container attribute by key. Empty result means
// attribute absence.
//
// See also SetAttribute, IterateAttributes.
func (x Container) Attribute(key string) string {
	attrs := x.v2.GetAttributes()
	for i := range attrs {
		if attrs[i].GetKey() == key {
			return attrs[i].GetValue()
		}
	}

	return ""
}

// IterateAttributes iterates over all Container attributes and passes them
// into f. The handler MUST NOT be nil.
//
// See also SetAttribute, Attribute.
func (x Container) IterateAttributes(f func(key, val string)) {
	attrs := x.v2.GetAttributes()
	for i := range attrs {
		f(attrs[i].GetKey(), attrs[i].GetValue())
	}
}

// SetName sets human-readable name of the Container. Name MUST NOT be empty.
//
// See also Name.
func SetName(cnr *Container, name string) {
	cnr.SetAttribute(attributeName, name)
}

// Name returns container name set using SetName.
//
// Zero Container has no name.
func Name(cnr Container) string {
	return cnr.Attribute(attributeName)
}

// SetCreationTime writes container's creation time in Unix Timestamp format.
//
// See also CreatedAt.
func SetCreationTime(cnr *Container, t time.Time) {
	cnr.SetAttribute(attributeTimestamp, strconv.FormatInt(t.Unix(), 10))
}

// CreatedAt returns container's creation time set using SetCreationTime.
//
// Zero Container has zero timestamp (in seconds).
func CreatedAt(cnr Container) time.Time {
	var sec int64

	attr := cnr.Attribute(attributeTimestamp)
	if attr != "" {
		var err error

		sec, err = strconv.ParseInt(cnr.Attribute(attributeTimestamp), 10, 64)
		if err != nil {
			panic(fmt.Sprintf("parse container timestamp: %v", err))
		}
	}

	return time.Unix(sec, 0)
}

// SetSubnet places the Container on the specified NeoFS subnet. If called,
// container nodes will only be selected from the given subnet, otherwise from
// the entire network.
func SetSubnet(cnr *Container, subNet subnetid.ID) {
	cnr.SetAttribute(container.SysAttributeSubnet, subNet.EncodeToString())
}

// Subnet return container subnet set using SetSubnet.
//
// Zero Container is bound to zero subnet.
func Subnet(cnr Container) (res subnetid.ID) {
	val := cnr.Attribute(container.SysAttributeSubnet)
	if val != "" {
		err := res.DecodeString(val)
		if err != nil {
			panic(fmt.Sprintf("invalid subnet attribute: %s (%v)", val, err))
		}
	}

	return
}

const attributeHomoHashEnabled = "true"

// DisableHomomorphicHashing sets flag to disable homomorphic hashing of the
// Container data.
//
// See also IsHomomorphicHashingDisabled.
func DisableHomomorphicHashing(cnr *Container) {
	cnr.SetAttribute(container.SysAttributeHomomorphicHashing, attributeHomoHashEnabled)
}

// IsHomomorphicHashingDisabled checks if DisableHomomorphicHashing was called.
//
// Zero Container has enabled hashing.
func IsHomomorphicHashingDisabled(cnr Container) bool {
	return cnr.Attribute(container.SysAttributeHomomorphicHashing) == attributeHomoHashEnabled
}

// Domain represents information about container domain registered in the NNS
// contract deployed in the NeoFS network.
type Domain struct {
	name, zone string
}

// SetName sets human-friendly container domain name.
func (x *Domain) SetName(name string) {
	x.name = name
}

// Name returns name set using SetName.
//
// Zero Domain has zero name.
func (x Domain) Name() string {
	return x.name
}

// SetZone sets zone which is used as a TLD of a domain name in NNS contract.
func (x *Domain) SetZone(zone string) {
	x.zone = zone
}

// Zone returns domain zone set using SetZone.
//
// Zero Domain has "container" zone.
func (x Domain) Zone() string {
	if x.zone != "" {
		return x.zone
	}

	return "container"
}

// WriteDomain writes Domain into the Container. Name MUST NOT be empty.
func WriteDomain(cnr *Container, domain Domain) {
	cnr.SetAttribute(container.SysAttributeName, domain.Name())
	cnr.SetAttribute(container.SysAttributeZone, domain.Zone())
}

// ReadDomain reads Domain from the Container. Returns value with empty name
// if domain is not specified.
func ReadDomain(cnr Container) (res Domain) {
	name := cnr.Attribute(container.SysAttributeName)
	if name != "" {
		res.SetName(name)
		res.SetZone(cnr.Attribute(container.SysAttributeZone))
	}

	return
}

// CalculateSignature calculates signature of the Container using provided signer
// and writes it into dst. Signature instance MUST NOT be nil. CalculateSignature
// is expected to be called after all the Container data is filled and before
// saving the Container in the NeoFS network. Note that мany subsequent change
// will most likely break the signature.
//
// See also VerifySignature.
func CalculateSignature(dst *frostfscrypto.Signature, cnr Container, signer ecdsa.PrivateKey) error {
	return dst.Calculate(frostfsecdsa.SignerRFC6979(signer), cnr.Marshal())
}

// VerifySignature verifies Container signature calculated using CalculateSignature.
// Result means signature correctness.
func VerifySignature(sig frostfscrypto.Signature, cnr Container) bool {
	return sig.Verify(cnr.Marshal())
}

// CalculateIDFromBinary calculates identifier of the binary-encoded container
// in CAS of the NeoFS containers and writes it into dst. ID instance MUST NOT
// be nil.
//
// See also CalculateID, AssertID.
func CalculateIDFromBinary(dst *cid.ID, cnr []byte) {
	dst.SetSHA256(sha256.Sum256(cnr))
}

// CalculateID encodes the given Container and passes the result into
// CalculateIDFromBinary.
//
// See also Container.Marshal, AssertID.
func CalculateID(dst *cid.ID, cnr Container) {
	CalculateIDFromBinary(dst, cnr.Marshal())
}

// AssertID checks if the given Container matches its identifier in CAS of the
// NeoFS containers.
//
// See also CalculateID.
func AssertID(id cid.ID, cnr Container) bool {
	var id2 cid.ID
	CalculateID(&id2, cnr)

	return id2.Equals(id)
}
