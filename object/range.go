package object

import (
	"github.com/TrueCloudLab/frostfs-api-go/v2/object"
)

// Range represents v2-compatible object payload range.
type Range object.Range

// NewRangeFromV2 wraps v2 Range message to Range.
//
// Nil object.Range converts to nil.
func NewRangeFromV2(rV2 *object.Range) *Range {
	return (*Range)(rV2)
}

// NewRange creates and initializes blank Range.
//
// Defaults:
//   - offset: 0;
//   - length: 0.
func NewRange() *Range {
	return NewRangeFromV2(new(object.Range))
}

// ToV2 converts Range to v2 Range message.
//
// Nil Range converts to nil.
func (r *Range) ToV2() *object.Range {
	return (*object.Range)(r)
}

// GetLength returns payload range size.
func (r *Range) GetLength() uint64 {
	return (*object.Range)(r).GetLength()
}

// SetLength sets payload range size.
func (r *Range) SetLength(v uint64) {
	(*object.Range)(r).SetLength(v)
}

// GetOffset sets payload range offset from start.
func (r *Range) GetOffset() uint64 {
	return (*object.Range)(r).GetOffset()
}

// SetOffset gets payload range offset from start.
func (r *Range) SetOffset(v uint64) {
	(*object.Range)(r).SetOffset(v)
}
