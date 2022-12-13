/*
Package usertest provides functions for convenient testing of user package API.

Note that importing the package into source files is highly discouraged.

Random instance generation functions can be useful when testing expects any value, e.g.:

	import usertest "github.com/TrueCloudLab/frostfs-sdk-go/user/test"

	id := usertest.ID()
	// test the value
*/
package usertest
