/*
Package accountingtest provides functions for convenient testing of accounting package API.

Note that importing the package into source files is highly discouraged.

Random instance generation functions can be useful when testing expects any value, e.g.:

	import accountingtest "github.com/TrueCloudLab/frostfs-sdk-go/accounting/test"

	dec := accountingtest.Decimal()
	// test the value
*/
package accountingtest
