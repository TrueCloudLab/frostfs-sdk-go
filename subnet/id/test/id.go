package subnetidtest

import (
	"math/rand"

	subnetid "github.com/TrueCloudLab/frostfs-sdk-go/subnet/id"
)

// ID generates and returns random subnetid.ID.
func ID() (x subnetid.ID) {
	x.SetNumeric(rand.Uint32())
	return
}
