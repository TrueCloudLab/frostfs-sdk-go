package object

import (
	"encoding/json"
	"strconv"

	v2object "github.com/TrueCloudLab/frostfs-api-go/v2/object"
	cid "github.com/TrueCloudLab/frostfs-sdk-go/container/id"
	oid "github.com/TrueCloudLab/frostfs-sdk-go/object/id"
	"github.com/TrueCloudLab/frostfs-sdk-go/user"
	"github.com/TrueCloudLab/frostfs-sdk-go/version"
)

// SearchMatchType indicates match operation on specified header.
type SearchMatchType uint32

const (
	MatchUnknown SearchMatchType = iota
	MatchStringEqual
	MatchStringNotEqual
	MatchNotPresent
	MatchCommonPrefix
)

func (m SearchMatchType) ToV2() v2object.MatchType {
	switch m {
	case MatchStringEqual:
		return v2object.MatchStringEqual
	case MatchStringNotEqual:
		return v2object.MatchStringNotEqual
	case MatchNotPresent:
		return v2object.MatchNotPresent
	case MatchCommonPrefix:
		return v2object.MatchCommonPrefix
	default:
		return v2object.MatchUnknown
	}
}

func SearchMatchFromV2(t v2object.MatchType) (m SearchMatchType) {
	switch t {
	case v2object.MatchStringEqual:
		m = MatchStringEqual
	case v2object.MatchStringNotEqual:
		m = MatchStringNotEqual
	case v2object.MatchNotPresent:
		m = MatchNotPresent
	case v2object.MatchCommonPrefix:
		m = MatchCommonPrefix
	default:
		m = MatchUnknown
	}

	return m
}

// String returns string representation of SearchMatchType.
//
// String mapping:
//   - MatchStringEqual: STRING_EQUAL;
//   - MatchStringNotEqual: STRING_NOT_EQUAL;
//   - MatchNotPresent: NOT_PRESENT;
//   - MatchCommonPrefix: COMMON_PREFIX;
//   - MatchUnknown, default: MATCH_TYPE_UNSPECIFIED.
func (m SearchMatchType) String() string {
	return m.ToV2().String()
}

// FromString parses SearchMatchType from a string representation.
// It is a reverse action to String().
//
// Returns true if s was parsed successfully.
func (m *SearchMatchType) FromString(s string) bool {
	var g v2object.MatchType

	ok := g.FromString(s)

	if ok {
		*m = SearchMatchFromV2(g)
	}

	return ok
}

type stringEncoder interface {
	EncodeToString() string
}

type SearchFilter struct {
	header filterKey
	value  stringEncoder
	op     SearchMatchType
}

type staticStringer string

type filterKey struct {
	typ filterKeyType

	str string
}

// enumeration of reserved filter keys.
type filterKeyType int

type SearchFilters []SearchFilter

const (
	_ filterKeyType = iota
	fKeyVersion
	fKeyObjectID
	fKeyContainerID
	fKeyOwnerID
	fKeyCreationEpoch
	fKeyPayloadLength
	fKeyPayloadHash
	fKeyType
	fKeyHomomorphicHash
	fKeyParent
	fKeySplitID
	fKeyPropRoot
	fKeyPropPhy
)

func (k filterKey) String() string {
	switch k.typ {
	default:
		return k.str
	case fKeyVersion:
		return v2object.FilterHeaderVersion
	case fKeyObjectID:
		return v2object.FilterHeaderObjectID
	case fKeyContainerID:
		return v2object.FilterHeaderContainerID
	case fKeyOwnerID:
		return v2object.FilterHeaderOwnerID
	case fKeyCreationEpoch:
		return v2object.FilterHeaderCreationEpoch
	case fKeyPayloadLength:
		return v2object.FilterHeaderPayloadLength
	case fKeyPayloadHash:
		return v2object.FilterHeaderPayloadHash
	case fKeyType:
		return v2object.FilterHeaderObjectType
	case fKeyHomomorphicHash:
		return v2object.FilterHeaderHomomorphicHash
	case fKeyParent:
		return v2object.FilterHeaderParent
	case fKeySplitID:
		return v2object.FilterHeaderSplitID
	case fKeyPropRoot:
		return v2object.FilterPropertyRoot
	case fKeyPropPhy:
		return v2object.FilterPropertyPhy
	}
}

func (s staticStringer) EncodeToString() string {
	return string(s)
}

func (f *SearchFilter) Header() string {
	return f.header.String()
}

func (f *SearchFilter) Value() string {
	return f.value.EncodeToString()
}

func (f *SearchFilter) Operation() SearchMatchType {
	return f.op
}

func NewSearchFilters() SearchFilters {
	return SearchFilters{}
}

func NewSearchFiltersFromV2(v2 []v2object.SearchFilter) SearchFilters {
	filters := make(SearchFilters, 0, len(v2))

	for i := range v2 {
		filters.AddFilter(
			v2[i].GetKey(),
			v2[i].GetValue(),
			SearchMatchFromV2(v2[i].GetMatchType()),
		)
	}

	return filters
}

func (f *SearchFilters) addFilter(op SearchMatchType, keyTyp filterKeyType, key string, val stringEncoder) {
	if *f == nil {
		*f = make(SearchFilters, 0, 1)
	}

	*f = append(*f, SearchFilter{
		header: filterKey{
			typ: keyTyp,
			str: key,
		},
		value: val,
		op:    op,
	})
}

func (f *SearchFilters) AddFilter(header, value string, op SearchMatchType) {
	f.addFilter(op, 0, header, staticStringer(value))
}

func (f *SearchFilters) addReservedFilter(op SearchMatchType, keyTyp filterKeyType, val stringEncoder) {
	f.addFilter(op, keyTyp, "", val)
}

// addFlagFilters adds filters that works like flags: they don't need to have
// specific match type or value. They processed by FrostFS nodes by the fact
// of presence in search query. E.g.: PHY, ROOT.
func (f *SearchFilters) addFlagFilter(keyTyp filterKeyType) {
	f.addFilter(MatchUnknown, keyTyp, "", staticStringer(""))
}

func (f *SearchFilters) AddObjectVersionFilter(op SearchMatchType, v version.Version) {
	f.addReservedFilter(op, fKeyVersion, staticStringer(version.EncodeToString(v)))
}

func (f *SearchFilters) AddObjectContainerIDFilter(m SearchMatchType, id cid.ID) {
	f.addReservedFilter(m, fKeyContainerID, id)
}

func (f *SearchFilters) AddObjectOwnerIDFilter(m SearchMatchType, id user.ID) {
	f.addReservedFilter(m, fKeyOwnerID, id)
}

func (f *SearchFilters) AddNotificationEpochFilter(epoch uint64) {
	f.addFilter(MatchStringEqual, 0, v2object.SysAttributeTickEpoch, staticStringer(strconv.FormatUint(epoch, 10)))
}

func (f SearchFilters) ToV2() []v2object.SearchFilter {
	result := make([]v2object.SearchFilter, len(f))

	for i := range f {
		result[i].SetKey(f[i].header.String())
		result[i].SetValue(f[i].value.EncodeToString())
		result[i].SetMatchType(f[i].op.ToV2())
	}

	return result
}

func (f *SearchFilters) addRootFilter() {
	f.addFlagFilter(fKeyPropRoot)
}

func (f *SearchFilters) AddRootFilter() {
	f.addRootFilter()
}

func (f *SearchFilters) addPhyFilter() {
	f.addFlagFilter(fKeyPropPhy)
}

func (f *SearchFilters) AddPhyFilter() {
	f.addPhyFilter()
}

// AddParentIDFilter adds filter by parent identifier.
func (f *SearchFilters) AddParentIDFilter(m SearchMatchType, id oid.ID) {
	f.addReservedFilter(m, fKeyParent, id)
}

// AddObjectIDFilter adds filter by object identifier.
func (f *SearchFilters) AddObjectIDFilter(m SearchMatchType, id oid.ID) {
	f.addReservedFilter(m, fKeyObjectID, id)
}

func (f *SearchFilters) AddSplitIDFilter(m SearchMatchType, id *SplitID) {
	f.addReservedFilter(m, fKeySplitID, staticStringer(id.String()))
}

// AddTypeFilter adds filter by object type.
func (f *SearchFilters) AddTypeFilter(m SearchMatchType, typ Type) {
	f.addReservedFilter(m, fKeyType, staticStringer(typ.String()))
}

// MarshalJSON encodes SearchFilters to protobuf JSON format.
func (f *SearchFilters) MarshalJSON() ([]byte, error) {
	return json.Marshal(f.ToV2())
}

// UnmarshalJSON decodes SearchFilters from protobuf JSON format.
func (f *SearchFilters) UnmarshalJSON(data []byte) error {
	var fsV2 []v2object.SearchFilter

	if err := json.Unmarshal(data, &fsV2); err != nil {
		return err
	}

	*f = NewSearchFiltersFromV2(fsV2)

	return nil
}
