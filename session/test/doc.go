/*
Package sessiontest provides functions for convenient testing of session package API.

Note that importing the package into source files is highly discouraged.

Random instance generation functions can be useful when testing expects any value, e.g.:

	import sessiontest "github.com/TrueCloudLab/frostfs-sdk-go/session/test"

	val := sessiontest.Container()
	// test the value
*/
package sessiontest
