/*
Package frostfscrypto collects NeoFS cryptographic primitives.

Signer type unifies entities for signing NeoFS data.

	// instantiate Signer
	// select data to be signed

	var sig Signature

	err := sig.Calculate(signer, data)
	// ...

	// attach signature to the request

SDK natively supports several signature schemes that are implemented
in nested packages.

PublicKey allows to verify signatures.

	// get signature to be verified
	// compose signed data

	isValid := sig.Verify(data)
	// ...

Signature can be also used to process NeoFS API V2 protocol messages
(see neo.fs.v2.refs package in https://github.com/TrueCloudLab/frostfs-api).

On client side:

	import "github.com/TrueCloudLab/frostfs-api-go/v2/refs"

	var msg refs.Signature
	sig.WriteToV2(&msg)

	// send msg

On server side:

	// recv msg

	var sig frostfscrypto.Signature
	sig.ReadFromV2(msg)

	// process sig

Using package types in an application is recommended to potentially work with
different protocol versions with which these types are compatible.
*/
package frostfscrypto
