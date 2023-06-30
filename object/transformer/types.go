package transformer

import (
	"io"

	"github.com/TrueCloudLab/frostfs-sdk-go/object"
	oid "github.com/TrueCloudLab/frostfs-sdk-go/object/id"
)

// AccessIdentifiers represents group of the object identifiers
// that are returned after writing the object.
// Consists of the ID of the stored object and the ID of the parent object.
type AccessIdentifiers struct {
	ParentID     *oid.ID
	SelfID       oid.ID
	ParentHeader *object.Object
}

// EpochSource is a source for the current epoch.
type EpochSource interface {
	CurrentEpoch() uint64
}

// ObjectTarget is an interface of the object writer.
type ObjectTarget interface {
	// WriteHeader writes object header w/ payload part.
	// The payload of the object may be incomplete.
	//
	// Must be called exactly once. Control remains with the caller.
	// Missing a call or re-calling can lead to undefined behavior
	// that depends on the implementation.
	//
	// Must not be called after Close call.
	WriteHeader(*object.Object) error

	// Write writes object payload chunk.
	//
	// Can be called multiple times.
	//
	// Must not be called after Close call.
	io.Writer

	// Close is used to finish object writing.
	//
	// Close must return access identifiers of the object
	// that has been written.
	//
	// Must be called no more than once. Control remains with the caller.
	// Re-calling can lead to undefined behavior
	// that depends on the implementation.
	Close() (*AccessIdentifiers, error)
}

// TargetInitializer represents ObjectTarget constructor.
type TargetInitializer func() ObjectTarget
