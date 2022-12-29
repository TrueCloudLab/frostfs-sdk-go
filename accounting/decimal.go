package accounting

import "github.com/TrueCloudLab/frostfs-api-go/v2/accounting"

// Decimal represents decimal number for accounting operations.
//
// Decimal is mutually compatible with github.com/TrueCloudLab/frostfs-api-go/v2/accounting.Decimal
// message. See ReadFromV2 / WriteToV2 methods.
//
// Instances can be created using built-in var declaration.
//
// Note that direct typecast is not safe and may result in loss of compatibility:
//
//	_ = Decimal(accounting.Decimal{}) // not recommended
type Decimal accounting.Decimal

// ReadFromV2 reads Decimal from the accounting.Decimal message. Checks if the
// message conforms to FrostFS API V2 protocol.
//
// See also WriteToV2.
func (d *Decimal) ReadFromV2(m accounting.Decimal) error {
	*d = Decimal(m)
	return nil
}

// WriteToV2 writes Decimal to the accounting.Decimal message.
// The message must not be nil.
//
// See also ReadFromV2.
func (d Decimal) WriteToV2(m *accounting.Decimal) {
	*m = (accounting.Decimal)(d)
}

// Value returns value of the decimal number.
//
// Zero Decimal has zero value.
//
// See also SetValue.
func (d Decimal) Value() int64 {
	return (*accounting.Decimal)(&d).GetValue()
}

// SetValue sets value of the decimal number.
//
// See also Value.
func (d *Decimal) SetValue(v int64) {
	(*accounting.Decimal)(d).SetValue(v)
}

// Precision returns precision of the decimal number.
//
// Zero Decimal has zero precision.
//
// See also SetPrecision.
func (d Decimal) Precision() uint32 {
	return (*accounting.Decimal)(&d).GetPrecision()
}

// SetPrecision sets precision of the decimal number.
//
// See also Precision.
func (d *Decimal) SetPrecision(p uint32) {
	(*accounting.Decimal)(d).SetPrecision(p)
}
