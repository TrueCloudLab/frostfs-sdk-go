/*
Package reputation collects functionality related to the FrostFS reputation system.

The functionality is based on the system described in the FrostFS specification.

Trust type represents simple instances of trust values. PeerToPeerTrust extends
Trust to support the direction of trust, i.e. from whom to whom. GlobalTrust
is designed as a global measure of trust in a network member. See the docs
for each type for details.

Instances can be also used to process FrostFS API V2 protocol messages
(see neo.fs.v2.reputation package in https://github.com/TrueCloudLab/frostfs-api).

On client side:

	import "github.com/TrueCloudLab/frostfs-api-go/v2/reputation"

	var msg reputation.GlobalTrust
	trust.WriteToV2(&msg)

	// send trust

On server side:

	// recv msg

	var trust reputation.GlobalTrust
	trust.ReadFromV2(msg)

	// process trust

Using package types in an application is recommended to potentially work with
different protocol versions with which these types are compatible.
*/
package reputation
