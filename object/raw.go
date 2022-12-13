package object

import (
	"github.com/TrueCloudLab/frostfs-api-go/v2/object"
)

// RawObject represents v2-compatible NeoFS object that provides
// a convenient interface to fill in the fields of
// an object in isolation from its internal structure.
//
// Deprecated: use Object type instead.
type RawObject = Object

// NewRawFromV2 wraps v2 Object message to Object.
//
// Deprecated: (v1.0.0) use NewFromV2 function instead.
func NewRawFromV2(oV2 *object.Object) *Object {
	return NewFromV2(oV2)
}

// NewRawFrom wraps Object instance to Object.
//
// Deprecated: (v1.0.0) function is no-op.
func NewRawFrom(obj *Object) *Object {
	return obj
}

// NewRaw creates and initializes blank Object.
//
// Deprecated: (v1.0.0) use New instead.
func NewRaw() *Object {
	return New()
}

// Object returns object instance.
//
// Deprecated: (v1.0.0) method is no-op, use arg directly.
func (o *Object) Object() *Object {
	return o
}
