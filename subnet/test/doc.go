/*
Package subnettest provides functions for convenient testing of subnet package API.

Note that importing the package into source files is highly discouraged.

Random instance generation functions can be useful when testing expects any value, e.g.:

	import subnettest "github.com/TrueCloudLab/frostfs-sdk-go/suibnet/test"

	value := subnettest.Info()
	// test the value
*/
package subnettest
