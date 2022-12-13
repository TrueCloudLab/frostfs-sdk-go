/*
Package oidtest provides functions for convenient testing of oid package API.

Note that importing the package into source files is highly discouraged.

Random instance generation functions can be useful when testing expects any value, e.g.:

	import oidtest "github.com/TrueCloudLab/frostfs-sdk-go/object/id/test"

	value := oidtest.ID()
	// test the value
*/
package oidtest
