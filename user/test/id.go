package usertest

import (
	"github.com/TrueCloudLab/frostfs-sdk-go/user"
	"github.com/nspcc-dev/neo-go/pkg/crypto/keys"
)

// ID returns random user.ID.
func ID() *user.ID {
	key, err := keys.NewPrivateKey()
	if err != nil {
		panic(err)
	}

	var x user.ID
	user.IDFromKey(&x, key.PrivateKey.PublicKey)

	return &x
}
