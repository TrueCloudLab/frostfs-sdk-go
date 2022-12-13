/*
Package subnetidtest provides functions for convenient testing of subnetid package API.

Note that importing the package into source files is highly discouraged.

Random instance generation functions can be useful when testing expects any value, e.g.:

	import subnetidtest "github.com/TrueCloudLab/frostfs-sdk-go/suibnet/id/test"

	value := subnetidtest.ID()
	// test the value
*/
package subnetidtest
