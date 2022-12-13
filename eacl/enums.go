package eacl

import (
	v2acl "github.com/TrueCloudLab/frostfs-api-go/v2/acl"
)

// Action taken if ContainerEACL record matched request.
// Action is compatible with v2 acl.Action enum.
type Action uint32

const (
	// ActionUnknown is an Action value used to mark action as undefined.
	ActionUnknown Action = iota

	// ActionAllow is an Action value that allows access to the operation from context.
	ActionAllow

	// ActionDeny is an Action value that denies access to the operation from context.
	ActionDeny
)

// Operation is a object service method to match request.
// Operation is compatible with v2 acl.Operation enum.
type Operation uint32

const (
	// OperationUnknown is an Operation value used to mark operation as undefined.
	OperationUnknown Operation = iota

	// OperationGet is an object get Operation.
	OperationGet

	// OperationHead is an Operation of getting the object header.
	OperationHead

	// OperationPut is an object put Operation.
	OperationPut

	// OperationDelete is an object delete Operation.
	OperationDelete

	// OperationSearch is an object search Operation.
	OperationSearch

	// OperationRange is an object payload range retrieval Operation.
	OperationRange

	// OperationRangeHash is an object payload range hashing Operation.
	OperationRangeHash
)

// Role is a group of request senders to match request.
// Role is compatible with v2 acl.Role enum.
type Role uint32

const (
	// RoleUnknown is a Role value used to mark role as undefined.
	RoleUnknown Role = iota

	// RoleUser is a group of senders that contains only key of container owner.
	RoleUser

	// RoleSystem is a group of senders that contains keys of container nodes and
	// inner ring nodes.
	RoleSystem

	// RoleOthers is a group of senders that contains none of above keys.
	RoleOthers
)

// Match is binary operation on filer name and value to check if request is matched.
// Match is compatible with v2 acl.MatchType enum.
type Match uint32

const (
	// MatchUnknown is a Match value used to mark matcher as undefined.
	MatchUnknown Match = iota

	// MatchStringEqual is a Match of string equality.
	MatchStringEqual

	// MatchStringNotEqual is a Match of string inequality.
	MatchStringNotEqual
)

// FilterHeaderType indicates source of headers to make matches.
// FilterHeaderType is compatible with v2 acl.HeaderType enum.
type FilterHeaderType uint32

const (
	// HeaderTypeUnknown is a FilterHeaderType value used to mark header type as undefined.
	HeaderTypeUnknown FilterHeaderType = iota

	// HeaderFromRequest is a FilterHeaderType for request X-Header.
	HeaderFromRequest

	// HeaderFromObject is a FilterHeaderType for object header.
	HeaderFromObject

	// HeaderFromService is a FilterHeaderType for service header.
	HeaderFromService
)

// ToV2 converts Action to v2 Action enum value.
func (a Action) ToV2() v2acl.Action {
	switch a {
	case ActionAllow:
		return v2acl.ActionAllow
	case ActionDeny:
		return v2acl.ActionDeny
	default:
		return v2acl.ActionUnknown
	}
}

// ActionFromV2 converts v2 Action enum value to Action.
func ActionFromV2(action v2acl.Action) (a Action) {
	switch action {
	case v2acl.ActionAllow:
		a = ActionAllow
	case v2acl.ActionDeny:
		a = ActionDeny
	default:
		a = ActionUnknown
	}

	return a
}

// String returns string representation of Action.
//
// String mapping:
//   - ActionAllow: ALLOW;
//   - ActionDeny: DENY;
//   - ActionUnknown, default: ACTION_UNSPECIFIED.
func (a Action) String() string {
	return a.ToV2().String()
}

// FromString parses Action from a string representation.
// It is a reverse action to String().
//
// Returns true if s was parsed successfully.
func (a *Action) FromString(s string) bool {
	var g v2acl.Action

	ok := g.FromString(s)

	if ok {
		*a = ActionFromV2(g)
	}

	return ok
}

// ToV2 converts Operation to v2 Operation enum value.
func (o Operation) ToV2() v2acl.Operation {
	switch o {
	case OperationGet:
		return v2acl.OperationGet
	case OperationHead:
		return v2acl.OperationHead
	case OperationPut:
		return v2acl.OperationPut
	case OperationDelete:
		return v2acl.OperationDelete
	case OperationSearch:
		return v2acl.OperationSearch
	case OperationRange:
		return v2acl.OperationRange
	case OperationRangeHash:
		return v2acl.OperationRangeHash
	default:
		return v2acl.OperationUnknown
	}
}

// OperationFromV2 converts v2 Operation enum value to Operation.
func OperationFromV2(operation v2acl.Operation) (o Operation) {
	switch operation {
	case v2acl.OperationGet:
		o = OperationGet
	case v2acl.OperationHead:
		o = OperationHead
	case v2acl.OperationPut:
		o = OperationPut
	case v2acl.OperationDelete:
		o = OperationDelete
	case v2acl.OperationSearch:
		o = OperationSearch
	case v2acl.OperationRange:
		o = OperationRange
	case v2acl.OperationRangeHash:
		o = OperationRangeHash
	default:
		o = OperationUnknown
	}

	return o
}

// String returns string representation of Operation.
//
// String mapping:
//   - OperationGet: GET;
//   - OperationHead: HEAD;
//   - OperationPut: PUT;
//   - OperationDelete: DELETE;
//   - OperationSearch: SEARCH;
//   - OperationRange: GETRANGE;
//   - OperationRangeHash: GETRANGEHASH;
//   - OperationUnknown, default: OPERATION_UNSPECIFIED.
func (o Operation) String() string {
	return o.ToV2().String()
}

// FromString parses Operation from a string representation.
// It is a reverse action to String().
//
// Returns true if s was parsed successfully.
func (o *Operation) FromString(s string) bool {
	var g v2acl.Operation

	ok := g.FromString(s)

	if ok {
		*o = OperationFromV2(g)
	}

	return ok
}

// ToV2 converts Role to v2 Role enum value.
func (r Role) ToV2() v2acl.Role {
	switch r {
	case RoleUser:
		return v2acl.RoleUser
	case RoleSystem:
		return v2acl.RoleSystem
	case RoleOthers:
		return v2acl.RoleOthers
	default:
		return v2acl.RoleUnknown
	}
}

// RoleFromV2 converts v2 Role enum value to Role.
func RoleFromV2(role v2acl.Role) (r Role) {
	switch role {
	case v2acl.RoleUser:
		r = RoleUser
	case v2acl.RoleSystem:
		r = RoleSystem
	case v2acl.RoleOthers:
		r = RoleOthers
	default:
		r = RoleUnknown
	}

	return r
}

// String returns string representation of Role.
//
// String mapping:
//   - RoleUser: USER;
//   - RoleSystem: SYSTEM;
//   - RoleOthers: OTHERS;
//   - RoleUnknown, default: ROLE_UNKNOWN.
func (r Role) String() string {
	return r.ToV2().String()
}

// FromString parses Role from a string representation.
// It is a reverse action to String().
//
// Returns true if s was parsed successfully.
func (r *Role) FromString(s string) bool {
	var g v2acl.Role

	ok := g.FromString(s)

	if ok {
		*r = RoleFromV2(g)
	}

	return ok
}

// ToV2 converts Match to v2 MatchType enum value.
func (m Match) ToV2() v2acl.MatchType {
	switch m {
	case MatchStringEqual:
		return v2acl.MatchTypeStringEqual
	case MatchStringNotEqual:
		return v2acl.MatchTypeStringNotEqual
	default:
		return v2acl.MatchTypeUnknown
	}
}

// MatchFromV2 converts v2 MatchType enum value to Match.
func MatchFromV2(match v2acl.MatchType) (m Match) {
	switch match {
	case v2acl.MatchTypeStringEqual:
		m = MatchStringEqual
	case v2acl.MatchTypeStringNotEqual:
		m = MatchStringNotEqual
	default:
		m = MatchUnknown
	}

	return m
}

// String returns string representation of Match.
//
// String mapping:
//   - MatchStringEqual: STRING_EQUAL;
//   - MatchStringNotEqual: STRING_NOT_EQUAL;
//   - MatchUnknown, default: MATCH_TYPE_UNSPECIFIED.
func (m Match) String() string {
	return m.ToV2().String()
}

// FromString parses Match from a string representation.
// It is a reverse action to String().
//
// Returns true if s was parsed successfully.
func (m *Match) FromString(s string) bool {
	var g v2acl.MatchType

	ok := g.FromString(s)

	if ok {
		*m = MatchFromV2(g)
	}

	return ok
}

// ToV2 converts FilterHeaderType to v2 HeaderType enum value.
func (h FilterHeaderType) ToV2() v2acl.HeaderType {
	switch h {
	case HeaderFromRequest:
		return v2acl.HeaderTypeRequest
	case HeaderFromObject:
		return v2acl.HeaderTypeObject
	case HeaderFromService:
		return v2acl.HeaderTypeService
	default:
		return v2acl.HeaderTypeUnknown
	}
}

// FilterHeaderTypeFromV2 converts v2 HeaderType enum value to FilterHeaderType.
func FilterHeaderTypeFromV2(header v2acl.HeaderType) (h FilterHeaderType) {
	switch header {
	case v2acl.HeaderTypeRequest:
		h = HeaderFromRequest
	case v2acl.HeaderTypeObject:
		h = HeaderFromObject
	case v2acl.HeaderTypeService:
		h = HeaderFromService
	default:
		h = HeaderTypeUnknown
	}

	return h
}

// String returns string representation of FilterHeaderType.
//
// String mapping:
//   - HeaderFromRequest: REQUEST;
//   - HeaderFromObject: OBJECT;
//   - HeaderTypeUnknown, default: HEADER_UNSPECIFIED.
func (h FilterHeaderType) String() string {
	return h.ToV2().String()
}

// FromString parses FilterHeaderType from a string representation.
// It is a reverse action to String().
//
// Returns true if s was parsed successfully.
func (h *FilterHeaderType) FromString(s string) bool {
	var g v2acl.HeaderType

	ok := g.FromString(s)

	if ok {
		*h = FilterHeaderTypeFromV2(g)
	}

	return ok
}
