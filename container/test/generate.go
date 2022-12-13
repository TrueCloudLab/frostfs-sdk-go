package containertest

import (
	"math/rand"

	"github.com/TrueCloudLab/frostfs-sdk-go/container"
	"github.com/TrueCloudLab/frostfs-sdk-go/container/acl"
	cidtest "github.com/TrueCloudLab/frostfs-sdk-go/container/id/test"
	netmaptest "github.com/TrueCloudLab/frostfs-sdk-go/netmap/test"
	usertest "github.com/TrueCloudLab/frostfs-sdk-go/user/test"
)

// Container returns random container.Container.
func Container() (x container.Container) {
	owner := usertest.ID()

	x.Init()
	x.SetAttribute("some attribute", "value")
	x.SetOwner(*owner)
	x.SetBasicACL(BasicACL())
	x.SetPlacementPolicy(netmaptest.PlacementPolicy())

	return x
}

// SizeEstimation returns random container.SizeEstimation.
func SizeEstimation() (x container.SizeEstimation) {
	x.SetContainer(cidtest.ID())
	x.SetEpoch(rand.Uint64())
	x.SetValue(rand.Uint64())

	return x
}

// BasicACL returns random acl.Basic.
func BasicACL() (x acl.Basic) {
	x.FromBits(rand.Uint32())
	return
}
