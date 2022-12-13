/*
Package checksumtest provides functions for convenient testing of checksum package API.

Note that importing the package into source files is highly discouraged.

Random instance generation functions can be useful when testing expects any value, e.g.:

	import checksumtest "github.com/TrueCloudLab/frostfs-sdk-go/checksum/test"

	cs := checksumtest.Checksum()
	// test the value
*/
package checksumtest
