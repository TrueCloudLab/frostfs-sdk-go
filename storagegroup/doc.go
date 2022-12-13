/*
Package storagegroup provides features to work with information that is
used for proof of storage in NeoFS system.

StorageGroup type groups verification values for Data Audit sessions:

	// receive sg info

	sg.ExpirationEpoch() // expiration of the storage group
	sg.Members() // objects in the group
	sg.ValidationDataHash() // hash for objects validation
	sg.ValidationDataSize() // total objects' payload size

Instances can be also used to process NeoFS API V2 protocol messages
(see neo.fs.v2.storagegroup package in https://github.com/TrueCloudLab/frostfs-api).

On client side:

	import "github.com/TrueCloudLab/frostfs-api-go/v2/storagegroup"

	var msg storagegroup.StorageGroup
	sg.WriteToV2(&msg)

	// send msg

On server side:

	// recv msg

	var sg StorageGroupDecimal
	sg.ReadFromV2(msg)

	// process sg

Using package types in an application is recommended to potentially work with
different protocol versions with which these types are compatible.
*/
package storagegroup
