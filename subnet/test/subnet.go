package subnettest

import (
	"github.com/TrueCloudLab/frostfs-sdk-go/subnet"
	subnetidtest "github.com/TrueCloudLab/frostfs-sdk-go/subnet/id/test"
	usertest "github.com/TrueCloudLab/frostfs-sdk-go/user/test"
)

// Info generates and returns random subnet.Info.
func Info() (x subnet.Info) {
	x.SetID(subnetidtest.ID())
	x.SetOwner(*usertest.ID())
	return
}
