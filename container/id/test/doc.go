/*
Package cidtest provides functions for convenient testing of cid package API.

Note that importing the package into source files is highly discouraged.

Random instance generation functions can be useful when testing expects any value, e.g.:

	import cidtest "github.com/TrueCloudLab/frostfs-sdk-go/container/id/test"

	cid := cidtest.ID()
	// test the value
*/
package cidtest
