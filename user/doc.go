/*
Package user provides functionality related to NeoFS users.

User identity is reflected in ID type. Each user has its own unique identifier
within the same network.

NeoFS user identification is compatible with Neo accounts:

	import "github.com/nspcc-dev/neo-go/pkg/crypto/keys"
	import "github.com/nspcc-dev/neo-go/pkg/crypto/hash"

	var id user.ID

	var scriptHash util.Uint160 // user account in NeoFS
	id.SetScriptHash(scriptHash)

	var key keys.PublicKey // user's public key
	user.IDFromKey(&id, k.PrivateKey.PublicKey)

ID is compatible with the NeoFS Smart Contract API:

	var id user.ID
	// ...

	wallet := id.WalletBytes()

	// use wallet in call

Encoding/decoding mechanisms are used to transfer identifiers:

	var id user.ID
	// ...

	s := id.EncodeToString() // on transmitter
	err = id.DecodeString(s) // on receiver

Instances can be also used to process NeoFS API protocol messages
(see neo.fs.v2.refs package in https://github.com/TrueCloudLab/frostfs-api).

On client side:

	import "github.com/TrueCloudLab/frostfs-api-go/v2/refs"

	var msg refs.OwnerID
	id.WriteToV2(&msg)

	// send msg

On server side:

	// recv msg

	var id user.ID

	err := id.ReadFromV2(msg)
	// ...

	// process id
*/
package user
