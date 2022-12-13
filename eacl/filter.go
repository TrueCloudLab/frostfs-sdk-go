package eacl

import (
	"strconv"

	v2acl "github.com/TrueCloudLab/frostfs-api-go/v2/acl"
)

// Filter defines check conditions if request header is matched or not. Matched
// header means that request should be processed according to ContainerEACL action.
//
// Filter is compatible with v2 acl.EACLRecord.Filter message.
type Filter struct {
	from    FilterHeaderType
	matcher Match
	key     filterKey
	value   stringEncoder
}

type staticStringer string

type u64Stringer uint64

type filterKey struct {
	typ filterKeyType

	str string
}

// enumeration of reserved filter keys.
type filterKeyType int

const (
	_ filterKeyType = iota
	fKeyObjVersion
	fKeyObjID
	fKeyObjContainerID
	fKeyObjOwnerID
	fKeyObjCreationEpoch
	fKeyObjPayloadLength
	fKeyObjPayloadHash
	fKeyObjType
	fKeyObjHomomorphicHash
	fKeyObjLast // helper, used in tests
)

func (s staticStringer) EncodeToString() string {
	return string(s)
}

func (u u64Stringer) EncodeToString() string {
	return strconv.FormatUint(uint64(u), 10)
}

// Value returns filtered string value.
func (f Filter) Value() string {
	return f.value.EncodeToString()
}

// Matcher returns filter Match type.
func (f Filter) Matcher() Match {
	return f.matcher
}

// Key returns key to the filtered header.
func (f Filter) Key() string {
	return f.key.String()
}

// From returns FilterHeaderType that defined which header will be filtered.
func (f Filter) From() FilterHeaderType {
	return f.from
}

// ToV2 converts Filter to v2 acl.EACLRecord.Filter message.
//
// Nil Filter converts to nil.
func (f *Filter) ToV2() *v2acl.HeaderFilter {
	if f == nil {
		return nil
	}

	filter := new(v2acl.HeaderFilter)
	filter.SetValue(f.value.EncodeToString())
	filter.SetKey(f.key.String())
	filter.SetMatchType(f.matcher.ToV2())
	filter.SetHeaderType(f.from.ToV2())

	return filter
}

func (k filterKey) String() string {
	switch k.typ {
	default:
		return k.str
	case fKeyObjVersion:
		return v2acl.FilterObjectVersion
	case fKeyObjID:
		return v2acl.FilterObjectID
	case fKeyObjContainerID:
		return v2acl.FilterObjectContainerID
	case fKeyObjOwnerID:
		return v2acl.FilterObjectOwnerID
	case fKeyObjCreationEpoch:
		return v2acl.FilterObjectCreationEpoch
	case fKeyObjPayloadLength:
		return v2acl.FilterObjectPayloadLength
	case fKeyObjPayloadHash:
		return v2acl.FilterObjectPayloadHash
	case fKeyObjType:
		return v2acl.FilterObjectType
	case fKeyObjHomomorphicHash:
		return v2acl.FilterObjectHomomorphicHash
	}
}

func (k *filterKey) fromString(s string) {
	switch s {
	default:
		k.typ, k.str = 0, s
	case v2acl.FilterObjectVersion:
		k.typ, k.str = fKeyObjVersion, ""
	case v2acl.FilterObjectID:
		k.typ, k.str = fKeyObjID, ""
	case v2acl.FilterObjectContainerID:
		k.typ, k.str = fKeyObjContainerID, ""
	case v2acl.FilterObjectOwnerID:
		k.typ, k.str = fKeyObjOwnerID, ""
	case v2acl.FilterObjectCreationEpoch:
		k.typ, k.str = fKeyObjCreationEpoch, ""
	case v2acl.FilterObjectPayloadLength:
		k.typ, k.str = fKeyObjPayloadLength, ""
	case v2acl.FilterObjectPayloadHash:
		k.typ, k.str = fKeyObjPayloadHash, ""
	case v2acl.FilterObjectType:
		k.typ, k.str = fKeyObjType, ""
	case v2acl.FilterObjectHomomorphicHash:
		k.typ, k.str = fKeyObjHomomorphicHash, ""
	}
}

// NewFilter creates, initializes and returns blank Filter instance.
//
// Defaults:
//   - header type: HeaderTypeUnknown;
//   - matcher: MatchUnknown;
//   - key: "";
//   - value: "".
func NewFilter() *Filter {
	return NewFilterFromV2(new(v2acl.HeaderFilter))
}

// NewFilterFromV2 converts v2 acl.EACLRecord.Filter message to Filter.
func NewFilterFromV2(filter *v2acl.HeaderFilter) *Filter {
	f := new(Filter)

	if filter == nil {
		return f
	}

	f.from = FilterHeaderTypeFromV2(filter.GetHeaderType())
	f.matcher = MatchFromV2(filter.GetMatchType())
	f.key.fromString(filter.GetKey())
	f.value = staticStringer(filter.GetValue())

	return f
}

// Marshal marshals Filter into a protobuf binary form.
func (f *Filter) Marshal() ([]byte, error) {
	return f.ToV2().StableMarshal(nil), nil
}

// Unmarshal unmarshals protobuf binary representation of Filter.
func (f *Filter) Unmarshal(data []byte) error {
	fV2 := new(v2acl.HeaderFilter)
	if err := fV2.Unmarshal(data); err != nil {
		return err
	}

	*f = *NewFilterFromV2(fV2)

	return nil
}

// MarshalJSON encodes Filter to protobuf JSON format.
func (f *Filter) MarshalJSON() ([]byte, error) {
	return f.ToV2().MarshalJSON()
}

// UnmarshalJSON decodes Filter from protobuf JSON format.
func (f *Filter) UnmarshalJSON(data []byte) error {
	fV2 := new(v2acl.HeaderFilter)
	if err := fV2.UnmarshalJSON(data); err != nil {
		return err
	}

	*f = *NewFilterFromV2(fV2)

	return nil
}

// equalFilters compares Filter with each other.
func equalFilters(f1, f2 Filter) bool {
	return f1.From() == f2.From() &&
		f1.Matcher() == f2.Matcher() &&
		f1.Key() == f2.Key() &&
		f1.Value() == f2.Value()
}
