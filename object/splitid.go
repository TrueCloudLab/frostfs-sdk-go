package object

import (
	"github.com/google/uuid"
)

// SplitID is a UUIDv4 used as attribute in split objects.
type SplitID struct {
	uuid uuid.UUID
}

// NewSplitID returns UUID representation of splitID attribute.
//
// Defaults:
//   - id: random UUID.
func NewSplitID() *SplitID {
	return &SplitID{
		uuid: uuid.New(),
	}
}

// NewSplitIDFromV2 returns parsed UUID from bytes.
// If v is invalid UUIDv4 byte sequence, then function returns nil.
//
// Nil converts to nil.
func NewSplitIDFromV2(v []byte) *SplitID {
	if v == nil {
		return nil
	}

	id := uuid.New()

	err := id.UnmarshalBinary(v)
	if err != nil {
		return nil
	}

	return &SplitID{
		uuid: id,
	}
}

// Parse converts UUIDv4 string representation into SplitID.
func (id *SplitID) Parse(s string) (err error) {
	id.uuid, err = uuid.Parse(s)
	if err != nil {
		return err
	}

	return nil
}

// String returns UUIDv4 string representation of SplitID.
func (id *SplitID) String() string {
	if id == nil {
		return ""
	}

	return id.uuid.String()
}

// SetUUID sets pre created UUID structure as SplitID.
func (id *SplitID) SetUUID(v uuid.UUID) {
	if id != nil {
		id.uuid = v
	}
}

// ToV2 converts SplitID to a representation of SplitID in frostfs-api v2.
//
// Nil SplitID converts to nil.
func (id *SplitID) ToV2() []byte {
	if id == nil {
		return nil
	}

	data, _ := id.uuid.MarshalBinary() // err is always nil

	return data
}
