package apistatus

import (
	"github.com/TrueCloudLab/frostfs-api-go/v2/session"
	"github.com/TrueCloudLab/frostfs-api-go/v2/status"
)

// SessionTokenNotFound describes status of the failure because of the missing session token.
// Instances provide Status and StatusV2 interfaces.
type SessionTokenNotFound struct {
	v2 status.Status
}

const defaultSessionTokenNotFoundMsg = "session token not found"

func (x SessionTokenNotFound) Error() string {
	msg := x.v2.Message()
	if msg == "" {
		msg = defaultSessionTokenNotFoundMsg
	}

	return errMessageStatusV2(
		globalizeCodeV2(session.StatusTokenNotFound, session.GlobalizeFail),
		msg,
	)
}

// implements local interface defined in FromStatusV2 func.
func (x *SessionTokenNotFound) fromStatusV2(st *status.Status) {
	x.v2 = *st
}

// ToStatusV2 implements StatusV2 interface method.
// If the value was returned by FromStatusV2, returns the source message.
// Otherwise, returns message with
//   - code: TOKEN_NOT_FOUND;
//   - string message: "session token not found";
//   - details: empty.
func (x SessionTokenNotFound) ToStatusV2() *status.Status {
	x.v2.SetCode(globalizeCodeV2(session.StatusTokenNotFound, session.GlobalizeFail))
	x.v2.SetMessage(defaultSessionTokenNotFoundMsg)
	return &x.v2
}

// SessionTokenExpired describes status of the failure because of the expired session token.
// Instances provide Status and StatusV2 interfaces.
type SessionTokenExpired struct {
	v2 status.Status
}

const defaultSessionTokenExpiredMsg = "expired session token"

func (x SessionTokenExpired) Error() string {
	msg := x.v2.Message()
	if msg == "" {
		msg = defaultSessionTokenExpiredMsg
	}

	return errMessageStatusV2(
		globalizeCodeV2(session.StatusTokenExpired, session.GlobalizeFail),
		msg,
	)
}

// implements local interface defined in FromStatusV2 func.
func (x *SessionTokenExpired) fromStatusV2(st *status.Status) {
	x.v2 = *st
}

// ToStatusV2 implements StatusV2 interface method.
// If the value was returned by FromStatusV2, returns the source message.
// Otherwise, returns message with
//   - code: TOKEN_EXPIRED;
//   - string message: "expired session token";
//   - details: empty.
func (x SessionTokenExpired) ToStatusV2() *status.Status {
	x.v2.SetCode(globalizeCodeV2(session.StatusTokenExpired, session.GlobalizeFail))
	x.v2.SetMessage(defaultSessionTokenExpiredMsg)
	return &x.v2
}
