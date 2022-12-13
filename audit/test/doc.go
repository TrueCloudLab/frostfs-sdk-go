/*
Package audittest provides functions for convenient testing of audit package API.

Note that importing the package into source files is highly discouraged.

Random instance generation functions can be useful when testing expects any value, e.g.:

	import audittest "github.com/TrueCloudLab/frostfs-sdk-go/audit/test"

	dec := audittest.Result()
	// test the value
*/
package audittest
