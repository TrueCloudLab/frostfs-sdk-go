package apistatus

import (
	"github.com/TrueCloudLab/frostfs-api-go/v2/container"
	"github.com/TrueCloudLab/frostfs-api-go/v2/status"
)

// ContainerNotFound describes status of the failure because of the missing container.
// Instances provide Status and StatusV2 interfaces.
type ContainerNotFound struct {
	v2 status.Status
}

const defaultContainerNotFoundMsg = "container not found"

func (x ContainerNotFound) Error() string {
	msg := x.v2.Message()
	if msg == "" {
		msg = defaultContainerNotFoundMsg
	}

	return errMessageStatusV2(
		globalizeCodeV2(container.StatusNotFound, container.GlobalizeFail),
		msg,
	)
}

// implements local interface defined in FromStatusV2 func.
func (x *ContainerNotFound) fromStatusV2(st *status.Status) {
	x.v2 = *st
}

// ToStatusV2 implements StatusV2 interface method.
// If the value was returned by FromStatusV2, returns the source message.
// Otherwise, returns message with
//   - code: CONTAINER_NOT_FOUND;
//   - string message: "container not found";
//   - details: empty.
func (x ContainerNotFound) ToStatusV2() *status.Status {
	x.v2.SetCode(globalizeCodeV2(container.StatusNotFound, container.GlobalizeFail))
	x.v2.SetMessage(defaultContainerNotFoundMsg)
	return &x.v2
}

// EACLNotFound describes status of the failure because of the missing eACL
// table.
// Instances provide Status and StatusV2 interfaces.
type EACLNotFound struct {
	v2 status.Status
}

const defaultEACLNotFoundMsg = "eACL not found"

func (x EACLNotFound) Error() string {
	msg := x.v2.Message()
	if msg == "" {
		msg = defaultEACLNotFoundMsg
	}

	return errMessageStatusV2(
		globalizeCodeV2(container.StatusEACLNotFound, container.GlobalizeFail),
		msg,
	)
}

// implements local interface defined in FromStatusV2 func.
func (x *EACLNotFound) fromStatusV2(st *status.Status) {
	x.v2 = *st
}

// ToStatusV2 implements StatusV2 interface method.
// If the value was returned by FromStatusV2, returns the source message.
// Otherwise, returns message with
//   - code: EACL_NOT_FOUND;
//   - string message: "eACL not found";
//   - details: empty.
func (x EACLNotFound) ToStatusV2() *status.Status {
	x.v2.SetCode(globalizeCodeV2(container.StatusEACLNotFound, container.GlobalizeFail))
	x.v2.SetMessage(defaultEACLNotFoundMsg)
	return &x.v2
}
