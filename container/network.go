package container

import (
	"github.com/TrueCloudLab/frostfs-sdk-go/netmap"
)

// ApplyNetworkConfig applies network configuration to the
// container. Changes the container if it does not satisfy
// network configuration.
func ApplyNetworkConfig(cnr *Container, cfg netmap.NetworkInfo) {
	if cfg.HomomorphicHashingDisabled() {
		DisableHomomorphicHashing(cnr)
	}
}

// AssertNetworkConfig checks if a container matches passed
// network configuration.
func AssertNetworkConfig(cnr Container, cfg netmap.NetworkInfo) bool {
	return IsHomomorphicHashingDisabled(cnr) == cfg.HomomorphicHashingDisabled()
}
