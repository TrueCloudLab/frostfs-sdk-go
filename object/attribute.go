package object

import (
	"github.com/TrueCloudLab/frostfs-api-go/v2/object"
)

// Attribute represents v2-compatible object attribute.
type Attribute object.Attribute

// NewAttributeFromV2 wraps v2 Attribute message to Attribute.
//
// Nil object.Attribute converts to nil.
func NewAttributeFromV2(aV2 *object.Attribute) *Attribute {
	return (*Attribute)(aV2)
}

// NewAttribute creates and initializes blank Attribute.
//
// Works similar as NewAttributeFromV2(new(Attribute)).
//
// Defaults:
//   - key: "";
//   - value: "".
func NewAttribute() *Attribute {
	return NewAttributeFromV2(new(object.Attribute))
}

// Key returns key to the object attribute.
func (a *Attribute) Key() string {
	return (*object.Attribute)(a).GetKey()
}

// SetKey sets key to the object attribute.
func (a *Attribute) SetKey(v string) {
	(*object.Attribute)(a).SetKey(v)
}

// Value return value of the object attribute.
func (a *Attribute) Value() string {
	return (*object.Attribute)(a).GetValue()
}

// SetValue sets value of the object attribute.
func (a *Attribute) SetValue(v string) {
	(*object.Attribute)(a).SetValue(v)
}

// ToV2 converts Attribute to v2 Attribute message.
//
// Nil Attribute converts to nil.
func (a *Attribute) ToV2() *object.Attribute {
	return (*object.Attribute)(a)
}

// Marshal marshals Attribute into a protobuf binary form.
func (a *Attribute) Marshal() ([]byte, error) {
	return (*object.Attribute)(a).StableMarshal(nil), nil
}

// Unmarshal unmarshals protobuf binary representation of Attribute.
func (a *Attribute) Unmarshal(data []byte) error {
	return (*object.Attribute)(a).Unmarshal(data)
}

// MarshalJSON encodes Attribute to protobuf JSON format.
func (a *Attribute) MarshalJSON() ([]byte, error) {
	return (*object.Attribute)(a).MarshalJSON()
}

// UnmarshalJSON decodes Attribute from protobuf JSON format.
func (a *Attribute) UnmarshalJSON(data []byte) error {
	return (*object.Attribute)(a).UnmarshalJSON(data)
}
