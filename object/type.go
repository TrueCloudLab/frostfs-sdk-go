package object

import (
	"github.com/TrueCloudLab/frostfs-api-go/v2/object"
)

type Type object.Type

const (
	TypeRegular Type = iota
	TypeTombstone
	TypeStorageGroup
	TypeLock
)

func (t Type) ToV2() object.Type {
	return object.Type(t)
}

func TypeFromV2(t object.Type) Type {
	return Type(t)
}

// String returns string representation of Type.
//
// String mapping:
//   - TypeTombstone: TOMBSTONE;
//   - TypeStorageGroup: STORAGE_GROUP;
//   - TypeLock: LOCK;
//   - TypeRegular, default: REGULAR.
func (t Type) String() string {
	return t.ToV2().String()
}

// FromString parses Type from a string representation.
// It is a reverse action to String().
//
// Returns true if s was parsed successfully.
func (t *Type) FromString(s string) bool {
	var g object.Type

	ok := g.FromString(s)

	if ok {
		*t = TypeFromV2(g)
	}

	return ok
}
