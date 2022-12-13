package eacl

import (
	cid "github.com/TrueCloudLab/frostfs-sdk-go/container/id"
)

// Header is an interface of string key-value header.
type Header interface {
	Key() string
	Value() string
}

// TypedHeaderSource is the interface that wraps
// method for selecting typed headers by type.
type TypedHeaderSource interface {
	// HeadersOfType returns the list of key-value headers
	// of particular type.
	//
	// It returns any problem encountered through the boolean
	// false value.
	HeadersOfType(FilterHeaderType) ([]Header, bool)
}

// ValidationUnit represents unit of check for Validator.
type ValidationUnit struct {
	cid *cid.ID

	role Role

	op Operation

	hdrSrc TypedHeaderSource

	key []byte

	table *Table
}

// WithContainerID configures ValidationUnit to use v as request's container ID.
func (u *ValidationUnit) WithContainerID(v *cid.ID) *ValidationUnit {
	if u != nil {
		u.cid = v
	}

	return u
}

// WithRole configures ValidationUnit to use v as request's role.
func (u *ValidationUnit) WithRole(v Role) *ValidationUnit {
	if u != nil {
		u.role = v
	}

	return u
}

// WithOperation configures ValidationUnit to use v as request's operation.
func (u *ValidationUnit) WithOperation(v Operation) *ValidationUnit {
	if u != nil {
		u.op = v
	}

	return u
}

// WithHeaderSource configures ValidationUnit to use v as a source of headers.
func (u *ValidationUnit) WithHeaderSource(v TypedHeaderSource) *ValidationUnit {
	if u != nil {
		u.hdrSrc = v
	}

	return u
}

// WithSenderKey configures ValidationUnit to use as sender's public key.
func (u *ValidationUnit) WithSenderKey(v []byte) *ValidationUnit {
	if u != nil {
		u.key = v
	}

	return u
}

// WithBearerToken configures ValidationUnit to use v as request's bearer token.
func (u *ValidationUnit) WithEACLTable(table *Table) *ValidationUnit {
	if u != nil {
		u.table = table
	}

	return u
}
