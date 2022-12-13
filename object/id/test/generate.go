package oidtest

import (
	"crypto/sha256"
	"math/rand"

	cidtest "github.com/TrueCloudLab/frostfs-sdk-go/container/id/test"
	oid "github.com/TrueCloudLab/frostfs-sdk-go/object/id"
)

// ID returns random oid.ID.
func ID() oid.ID {
	checksum := [sha256.Size]byte{}

	rand.Read(checksum[:])

	return idWithChecksum(checksum)
}

// idWithChecksum returns oid.ID initialized
// with specified checksum.
func idWithChecksum(cs [sha256.Size]byte) oid.ID {
	var id oid.ID
	id.SetSHA256(cs)

	return id
}

// Address returns random oid.Address.
func Address() oid.Address {
	var x oid.Address

	x.SetContainer(cidtest.ID())
	x.SetObject(ID())

	return x
}
