package frostfsecdsa

import frostfscrypto "github.com/TrueCloudLab/frostfs-sdk-go/crypto"

func init() {
	frostfscrypto.RegisterScheme(frostfscrypto.ECDSA_SHA512, func() frostfscrypto.PublicKey {
		return new(PublicKey)
	})

	frostfscrypto.RegisterScheme(frostfscrypto.ECDSA_DETERMINISTIC_SHA256, func() frostfscrypto.PublicKey {
		return new(PublicKeyRFC6979)
	})

	frostfscrypto.RegisterScheme(frostfscrypto.ECDSA_WALLETCONNECT, func() frostfscrypto.PublicKey {
		return new(PublicKeyWalletConnect)
	})
}
