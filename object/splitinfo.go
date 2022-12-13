package object

import (
	"errors"
	"fmt"

	"github.com/TrueCloudLab/frostfs-api-go/v2/object"
	"github.com/TrueCloudLab/frostfs-api-go/v2/refs"
	oid "github.com/TrueCloudLab/frostfs-sdk-go/object/id"
)

type SplitInfo object.SplitInfo

// NewSplitInfoFromV2 wraps v2 SplitInfo message to SplitInfo.
//
// Nil object.SplitInfo converts to nil.
func NewSplitInfoFromV2(v2 *object.SplitInfo) *SplitInfo {
	return (*SplitInfo)(v2)
}

// NewSplitInfo creates and initializes blank SplitInfo.
//
// Defaults:
//   - splitID: nil;
//   - lastPart nil;
//   - link: nil.
func NewSplitInfo() *SplitInfo {
	return NewSplitInfoFromV2(new(object.SplitInfo))
}

// ToV2 converts SplitInfo to v2 SplitInfo message.
//
// Nil SplitInfo converts to nil.
func (s *SplitInfo) ToV2() *object.SplitInfo {
	return (*object.SplitInfo)(s)
}

func (s *SplitInfo) SplitID() *SplitID {
	return NewSplitIDFromV2(
		(*object.SplitInfo)(s).GetSplitID())
}

func (s *SplitInfo) SetSplitID(v *SplitID) {
	(*object.SplitInfo)(s).SetSplitID(v.ToV2())
}

func (s SplitInfo) LastPart() (v oid.ID, isSet bool) {
	v2 := (object.SplitInfo)(s)

	lpV2 := v2.GetLastPart()
	if lpV2 != nil {
		_ = v.ReadFromV2(*lpV2)
		isSet = true
	}

	return
}

func (s *SplitInfo) SetLastPart(v oid.ID) {
	var idV2 refs.ObjectID
	v.WriteToV2(&idV2)

	(*object.SplitInfo)(s).SetLastPart(&idV2)
}

func (s SplitInfo) Link() (v oid.ID, isSet bool) {
	v2 := (object.SplitInfo)(s)

	linkV2 := v2.GetLink()
	if linkV2 != nil {
		_ = v.ReadFromV2(*linkV2)
		isSet = true
	}

	return
}

func (s *SplitInfo) SetLink(v oid.ID) {
	var idV2 refs.ObjectID
	v.WriteToV2(&idV2)

	(*object.SplitInfo)(s).SetLink(&idV2)
}

func (s *SplitInfo) Marshal() ([]byte, error) {
	return (*object.SplitInfo)(s).StableMarshal(nil), nil
}

func (s *SplitInfo) Unmarshal(data []byte) error {
	err := (*object.SplitInfo)(s).Unmarshal(data)
	if err != nil {
		return err
	}

	return formatCheckSI((*object.SplitInfo)(s))
}

// MarshalJSON implements json.Marshaler.
func (s *SplitInfo) MarshalJSON() ([]byte, error) {
	return (*object.SplitInfo)(s).MarshalJSON()
}

// UnmarshalJSON implements json.Unmarshaler.
func (s *SplitInfo) UnmarshalJSON(data []byte) error {
	err := (*object.SplitInfo)(s).UnmarshalJSON(data)
	if err != nil {
		return err
	}

	return formatCheckSI((*object.SplitInfo)(s))
}

var errSplitInfoMissingFields = errors.New("neither link object ID nor last part object ID is set")

func formatCheckSI(v2 *object.SplitInfo) error {
	link := v2.GetLink()
	lastPart := v2.GetLastPart()
	if link == nil && lastPart == nil {
		return errSplitInfoMissingFields
	}

	var oID oid.ID

	if link != nil {
		err := oID.ReadFromV2(*link)
		if err != nil {
			return fmt.Errorf("could not convert link object ID: %w", err)
		}
	}

	if lastPart != nil {
		err := oID.ReadFromV2(*lastPart)
		if err != nil {
			return fmt.Errorf("could not convert last part object ID: %w", err)
		}
	}

	return nil
}
