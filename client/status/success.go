package apistatus

import (
	"github.com/TrueCloudLab/frostfs-api-go/v2/status"
)

// SuccessDefaultV2 represents Status instance of default success. Implements StatusV2.
type SuccessDefaultV2 struct {
	isNil bool

	v2 *status.Status
}

// implements local interface defined in FromStatusV2 func.
func (x *SuccessDefaultV2) fromStatusV2(st *status.Status) {
	x.isNil = st == nil
	x.v2 = st
}

// ToStatusV2 implements StatusV2 interface method.
// If the value was returned by FromStatusV2, returns the source message.
// Otherwise, returns message with
//   - code: OK;
//   - string message: empty;
//   - details: empty.
func (x SuccessDefaultV2) ToStatusV2() *status.Status {
	if x.isNil || x.v2 != nil {
		return x.v2
	}

	return newStatusV2WithLocalCode(status.OK, status.GlobalizeSuccess)
}
