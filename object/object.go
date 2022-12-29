package object

import (
	"errors"
	"fmt"

	"github.com/TrueCloudLab/frostfs-api-go/v2/object"
	"github.com/TrueCloudLab/frostfs-api-go/v2/refs"
	v2session "github.com/TrueCloudLab/frostfs-api-go/v2/session"
	"github.com/TrueCloudLab/frostfs-sdk-go/checksum"
	cid "github.com/TrueCloudLab/frostfs-sdk-go/container/id"
	frostfscrypto "github.com/TrueCloudLab/frostfs-sdk-go/crypto"
	oid "github.com/TrueCloudLab/frostfs-sdk-go/object/id"
	"github.com/TrueCloudLab/frostfs-sdk-go/session"
	"github.com/TrueCloudLab/frostfs-sdk-go/user"
	"github.com/TrueCloudLab/frostfs-sdk-go/version"
)

// Object represents in-memory structure of the FrostFS object.
// Type is compatible with FrostFS API V2 protocol.
//
// Instance can be created depending on scenario:
//   - InitCreation (an object to be placed in container);
//   - New (blank instance, usually needed for decoding);
//   - NewFromV2 (when working under FrostFS API V2 protocol).
type Object object.Object

// RequiredFields contains the minimum set of object data that must be set
// by the FrostFS user at the stage of creation.
type RequiredFields struct {
	// Identifier of the FrostFS container associated with the object.
	Container cid.ID

	// Object owner's user ID in the FrostFS system.
	Owner user.ID
}

// InitCreation initializes the object instance with minimum set of required fields.
// Object is expected (but not required) to be blank. Object must not be nil.
func InitCreation(dst *Object, rf RequiredFields) {
	dst.SetContainerID(rf.Container)
	dst.SetOwnerID(&rf.Owner)
}

// NewFromV2 wraps v2 Object message to Object.
func NewFromV2(oV2 *object.Object) *Object {
	return (*Object)(oV2)
}

// New creates and initializes blank Object.
//
// Works similar as NewFromV2(new(Object)).
func New() *Object {
	return NewFromV2(new(object.Object))
}

// ToV2 converts Object to v2 Object message.
func (o *Object) ToV2() *object.Object {
	return (*object.Object)(o)
}

// MarshalHeaderJSON marshals object's header
// into JSON format.
func (o *Object) MarshalHeaderJSON() ([]byte, error) {
	return (*object.Object)(o).GetHeader().MarshalJSON()
}

func (o *Object) setHeaderField(setter func(*object.Header)) {
	obj := (*object.Object)(o)
	h := obj.GetHeader()

	if h == nil {
		h = new(object.Header)
		obj.SetHeader(h)
	}

	setter(h)
}

func (o *Object) setSplitFields(setter func(*object.SplitHeader)) {
	o.setHeaderField(func(h *object.Header) {
		split := h.GetSplit()
		if split == nil {
			split = new(object.SplitHeader)
			h.SetSplit(split)
		}

		setter(split)
	})
}

// ID returns object identifier.
func (o *Object) ID() (v oid.ID, isSet bool) {
	v2 := (*object.Object)(o)
	if id := v2.GetObjectID(); id != nil {
		_ = v.ReadFromV2(*v2.GetObjectID())
		isSet = true
	}

	return
}

// SetID sets object identifier.
func (o *Object) SetID(v oid.ID) {
	var v2 refs.ObjectID
	v.WriteToV2(&v2)

	(*object.Object)(o).
		SetObjectID(&v2)
}

// Signature returns signature of the object identifier.
func (o *Object) Signature() *frostfscrypto.Signature {
	sigv2 := (*object.Object)(o).GetSignature()
	if sigv2 == nil {
		return nil
	}

	var sig frostfscrypto.Signature
	sig.ReadFromV2(*sigv2) // FIXME(@cthulhu-rider): #226 handle error

	return &sig
}

// SetSignature sets signature of the object identifier.
func (o *Object) SetSignature(v *frostfscrypto.Signature) {
	var sigv2 *refs.Signature

	if v != nil {
		sigv2 = new(refs.Signature)

		v.WriteToV2(sigv2)
	}

	(*object.Object)(o).SetSignature(sigv2)
}

// Payload returns payload bytes.
func (o *Object) Payload() []byte {
	return (*object.Object)(o).GetPayload()
}

// SetPayload sets payload bytes.
func (o *Object) SetPayload(v []byte) {
	(*object.Object)(o).SetPayload(v)
}

// Version returns version of the object.
func (o *Object) Version() *version.Version {
	var ver version.Version
	if verV2 := (*object.Object)(o).GetHeader().GetVersion(); verV2 != nil {
		ver.ReadFromV2(*verV2) // FIXME(@cthulhu-rider): #226 handle error
	}
	return &ver
}

// SetVersion sets version of the object.
func (o *Object) SetVersion(v *version.Version) {
	var verV2 refs.Version
	v.WriteToV2(&verV2)

	o.setHeaderField(func(h *object.Header) {
		h.SetVersion(&verV2)
	})
}

// PayloadSize returns payload length of the object.
func (o *Object) PayloadSize() uint64 {
	return (*object.Object)(o).
		GetHeader().
		GetPayloadLength()
}

// SetPayloadSize sets payload length of the object.
func (o *Object) SetPayloadSize(v uint64) {
	o.setHeaderField(func(h *object.Header) {
		h.SetPayloadLength(v)
	})
}

// ContainerID returns identifier of the related container.
func (o *Object) ContainerID() (v cid.ID, isSet bool) {
	v2 := (*object.Object)(o)

	cidV2 := v2.GetHeader().GetContainerID()
	if cidV2 != nil {
		_ = v.ReadFromV2(*cidV2)
		isSet = true
	}

	return
}

// SetContainerID sets identifier of the related container.
func (o *Object) SetContainerID(v cid.ID) {
	var cidV2 refs.ContainerID
	v.WriteToV2(&cidV2)

	o.setHeaderField(func(h *object.Header) {
		h.SetContainerID(&cidV2)
	})
}

// OwnerID returns identifier of the object owner.
func (o *Object) OwnerID() *user.ID {
	var id user.ID

	m := (*object.Object)(o).GetHeader().GetOwnerID()
	if m != nil {
		_ = id.ReadFromV2(*m)
	}

	return &id
}

// SetOwnerID sets identifier of the object owner.
func (o *Object) SetOwnerID(v *user.ID) {
	o.setHeaderField(func(h *object.Header) {
		var m refs.OwnerID
		v.WriteToV2(&m)

		h.SetOwnerID(&m)
	})
}

// CreationEpoch returns epoch number in which object was created.
func (o *Object) CreationEpoch() uint64 {
	return (*object.Object)(o).
		GetHeader().
		GetCreationEpoch()
}

// SetCreationEpoch sets epoch number in which object was created.
func (o *Object) SetCreationEpoch(v uint64) {
	o.setHeaderField(func(h *object.Header) {
		h.SetCreationEpoch(v)
	})
}

// PayloadChecksum returns checksum of the object payload and
// bool that indicates checksum presence in the object.
//
// Zero Object does not have payload checksum.
//
// See also SetPayloadChecksum.
func (o *Object) PayloadChecksum() (checksum.Checksum, bool) {
	var v checksum.Checksum
	v2 := (*object.Object)(o)

	if hash := v2.GetHeader().GetPayloadHash(); hash != nil {
		v.ReadFromV2(*hash) // FIXME(@cthulhu-rider): #226 handle error
		return v, true
	}

	return v, false
}

// SetPayloadChecksum sets checksum of the object payload.
//
// See also PayloadChecksum.
func (o *Object) SetPayloadChecksum(v checksum.Checksum) {
	var v2 refs.Checksum
	v.WriteToV2(&v2)

	o.setHeaderField(func(h *object.Header) {
		h.SetPayloadHash(&v2)
	})
}

// PayloadHomomorphicHash returns homomorphic hash of the object
// payload and bool that indicates checksum presence in the object.
//
// Zero Object does not have payload homomorphic checksum.
//
// See also SetPayloadHomomorphicHash.
func (o *Object) PayloadHomomorphicHash() (checksum.Checksum, bool) {
	var v checksum.Checksum
	v2 := (*object.Object)(o)

	if hash := v2.GetHeader().GetHomomorphicHash(); hash != nil {
		v.ReadFromV2(*hash) // FIXME(@cthulhu-rider): #226 handle error
		return v, true
	}

	return v, false
}

// SetPayloadHomomorphicHash sets homomorphic hash of the object payload.
//
// See also PayloadHomomorphicHash.
func (o *Object) SetPayloadHomomorphicHash(v checksum.Checksum) {
	var v2 refs.Checksum
	v.WriteToV2(&v2)

	o.setHeaderField(func(h *object.Header) {
		h.SetHomomorphicHash(&v2)
	})
}

// Attributes returns object attributes.
func (o *Object) Attributes() []Attribute {
	attrs := (*object.Object)(o).
		GetHeader().
		GetAttributes()

	res := make([]Attribute, len(attrs))

	for i := range attrs {
		res[i] = *NewAttributeFromV2(&attrs[i])
	}

	return res
}

// SetAttributes sets object attributes.
func (o *Object) SetAttributes(v ...Attribute) {
	attrs := make([]object.Attribute, len(v))

	for i := range v {
		attrs[i] = *v[i].ToV2()
	}

	o.setHeaderField(func(h *object.Header) {
		h.SetAttributes(attrs)
	})
}

// PreviousID returns identifier of the previous sibling object.
func (o *Object) PreviousID() (v oid.ID, isSet bool) {
	v2 := (*object.Object)(o)

	v2Prev := v2.GetHeader().GetSplit().GetPrevious()
	if v2Prev != nil {
		_ = v.ReadFromV2(*v2Prev)
		isSet = true
	}

	return
}

// SetPreviousID sets identifier of the previous sibling object.
func (o *Object) SetPreviousID(v oid.ID) {
	var v2 refs.ObjectID
	v.WriteToV2(&v2)

	o.setSplitFields(func(split *object.SplitHeader) {
		split.SetPrevious(&v2)
	})
}

// Children return list of the identifiers of the child objects.
func (o *Object) Children() []oid.ID {
	v2 := (*object.Object)(o)
	ids := v2.GetHeader().GetSplit().GetChildren()

	var (
		id  oid.ID
		res = make([]oid.ID, len(ids))
	)

	for i := range ids {
		_ = id.ReadFromV2(ids[i])
		res[i] = id
	}

	return res
}

// SetChildren sets list of the identifiers of the child objects.
func (o *Object) SetChildren(v ...oid.ID) {
	var (
		v2  refs.ObjectID
		ids = make([]refs.ObjectID, len(v))
	)

	for i := range v {
		v[i].WriteToV2(&v2)
		ids[i] = v2
	}

	o.setSplitFields(func(split *object.SplitHeader) {
		split.SetChildren(ids)
	})
}

// NotificationInfo groups information about object notification
// that can be written to object.
//
// Topic is an optional field.
type NotificationInfo struct {
	ni object.NotificationInfo
}

// Epoch returns object notification tick
// epoch.
func (n NotificationInfo) Epoch() uint64 {
	return n.ni.Epoch()
}

// SetEpoch sets object notification tick
// epoch.
func (n *NotificationInfo) SetEpoch(epoch uint64) {
	n.ni.SetEpoch(epoch)
}

// Topic return optional object notification
// topic.
func (n NotificationInfo) Topic() string {
	return n.ni.Topic()
}

// SetTopic sets optional object notification
// topic.
func (n *NotificationInfo) SetTopic(topic string) {
	n.ni.SetTopic(topic)
}

// NotificationInfo returns notification info
// read from the object structure.
// Returns any error that appeared during notification
// information parsing.
func (o *Object) NotificationInfo() (*NotificationInfo, error) {
	ni, err := object.GetNotificationInfo((*object.Object)(o))
	if err != nil {
		return nil, err
	}

	return &NotificationInfo{
		ni: *ni,
	}, nil
}

// SetNotification writes NotificationInfo to the object structure.
func (o *Object) SetNotification(ni NotificationInfo) {
	object.WriteNotificationInfo((*object.Object)(o), ni.ni)
}

// SplitID return split identity of split object. If object is not split
// returns nil.
func (o *Object) SplitID() *SplitID {
	return NewSplitIDFromV2(
		(*object.Object)(o).
			GetHeader().
			GetSplit().
			GetSplitID(),
	)
}

// SetSplitID sets split identifier for the split object.
func (o *Object) SetSplitID(id *SplitID) {
	o.setSplitFields(func(split *object.SplitHeader) {
		split.SetSplitID(id.ToV2())
	})
}

// ParentID returns identifier of the parent object.
func (o *Object) ParentID() (v oid.ID, isSet bool) {
	v2 := (*object.Object)(o)

	v2Par := v2.GetHeader().GetSplit().GetParent()
	if v2Par != nil {
		_ = v.ReadFromV2(*v2Par)
		isSet = true
	}

	return
}

// SetParentID sets identifier of the parent object.
func (o *Object) SetParentID(v oid.ID) {
	var v2 refs.ObjectID
	v.WriteToV2(&v2)

	o.setSplitFields(func(split *object.SplitHeader) {
		split.SetParent(&v2)
	})
}

// Parent returns parent object w/o payload.
func (o *Object) Parent() *Object {
	h := (*object.Object)(o).
		GetHeader().
		GetSplit()

	parSig := h.GetParentSignature()
	parHdr := h.GetParentHeader()

	if parSig == nil && parHdr == nil {
		return nil
	}

	oV2 := new(object.Object)
	oV2.SetObjectID(h.GetParent())
	oV2.SetSignature(parSig)
	oV2.SetHeader(parHdr)

	return NewFromV2(oV2)
}

// SetParent sets parent object w/o payload.
func (o *Object) SetParent(v *Object) {
	o.setSplitFields(func(split *object.SplitHeader) {
		split.SetParent((*object.Object)(v).GetObjectID())
		split.SetParentSignature((*object.Object)(v).GetSignature())
		split.SetParentHeader((*object.Object)(v).GetHeader())
	})
}

func (o *Object) initRelations() {
	o.setHeaderField(func(h *object.Header) {
		h.SetSplit(new(object.SplitHeader))
	})
}

func (o *Object) resetRelations() {
	o.setHeaderField(func(h *object.Header) {
		h.SetSplit(nil)
	})
}

// SessionToken returns token of the session
// within which object was created.
func (o *Object) SessionToken() *session.Object {
	tokv2 := (*object.Object)(o).GetHeader().GetSessionToken()
	if tokv2 == nil {
		return nil
	}

	var res session.Object

	_ = res.ReadFromV2(*tokv2)

	return &res
}

// SetSessionToken sets token of the session
// within which object was created.
func (o *Object) SetSessionToken(v *session.Object) {
	o.setHeaderField(func(h *object.Header) {
		var tokv2 *v2session.Token

		if v != nil {
			tokv2 = new(v2session.Token)
			v.WriteToV2(tokv2)
		}

		h.SetSessionToken(tokv2)
	})
}

// Type returns type of the object.
func (o *Object) Type() Type {
	return TypeFromV2(
		(*object.Object)(o).
			GetHeader().
			GetObjectType(),
	)
}

// SetType sets type of the object.
func (o *Object) SetType(v Type) {
	o.setHeaderField(func(h *object.Header) {
		h.SetObjectType(v.ToV2())
	})
}

// CutPayload returns Object w/ empty payload.
//
// Changes of non-payload fields affect source object.
func (o *Object) CutPayload() *Object {
	ov2 := new(object.Object)
	*ov2 = *(*object.Object)(o)
	ov2.SetPayload(nil)

	return (*Object)(ov2)
}

func (o *Object) HasParent() bool {
	return (*object.Object)(o).
		GetHeader().
		GetSplit() != nil
}

// ResetRelations removes all fields of links with other objects.
func (o *Object) ResetRelations() {
	o.resetRelations()
}

// InitRelations initializes relation field.
func (o *Object) InitRelations() {
	o.initRelations()
}

// Marshal marshals object into a protobuf binary form.
func (o *Object) Marshal() ([]byte, error) {
	return (*object.Object)(o).StableMarshal(nil), nil
}

// Unmarshal unmarshals protobuf binary representation of object.
func (o *Object) Unmarshal(data []byte) error {
	err := (*object.Object)(o).Unmarshal(data)
	if err != nil {
		return err
	}

	return formatCheck((*object.Object)(o))
}

// MarshalJSON encodes object to protobuf JSON format.
func (o *Object) MarshalJSON() ([]byte, error) {
	return (*object.Object)(o).MarshalJSON()
}

// UnmarshalJSON decodes object from protobuf JSON format.
func (o *Object) UnmarshalJSON(data []byte) error {
	err := (*object.Object)(o).UnmarshalJSON(data)
	if err != nil {
		return err
	}

	return formatCheck((*object.Object)(o))
}

var errOIDNotSet = errors.New("object ID is not set")
var errCIDNotSet = errors.New("container ID is not set")

func formatCheck(v2 *object.Object) error {
	var (
		oID oid.ID
		cID cid.ID
	)

	oidV2 := v2.GetObjectID()
	if oidV2 == nil {
		return errOIDNotSet
	}

	err := oID.ReadFromV2(*oidV2)
	if err != nil {
		return fmt.Errorf("could not convert V2 object ID: %w", err)
	}

	cidV2 := v2.GetHeader().GetContainerID()
	if cidV2 == nil {
		return errCIDNotSet
	}

	err = cID.ReadFromV2(*cidV2)
	if err != nil {
		return fmt.Errorf("could not convert V2 container ID: %w", err)
	}

	if prev := v2.GetHeader().GetSplit().GetPrevious(); prev != nil {
		err = oID.ReadFromV2(*prev)
		if err != nil {
			return fmt.Errorf("could not convert previous object ID: %w", err)
		}
	}

	if parent := v2.GetHeader().GetSplit().GetParent(); parent != nil {
		err = oID.ReadFromV2(*parent)
		if err != nil {
			return fmt.Errorf("could not convert parent object ID: %w", err)
		}
	}

	return nil
}
