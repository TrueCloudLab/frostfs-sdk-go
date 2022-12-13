package frostfsecdsa

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"encoding/base64"
	"fmt"

	"github.com/TrueCloudLab/frostfs-api-go/v2/util/signature/walletconnect"
	frostfscrypto "github.com/TrueCloudLab/frostfs-sdk-go/crypto"
	"github.com/nspcc-dev/neo-go/pkg/crypto/keys"
)

// SignerWalletConnect is similar to SignerRFC6979 with 2 changes:
// 1. The data is base64 encoded before signing/verifying.
// 2. The signature is a concatenation of the signature itself and 16-byte salt.
//
// Instances MUST be initialized from ecdsa.PrivateKey using type conversion.
type SignerWalletConnect ecdsa.PrivateKey

// Scheme returns frostfscrypto.ECDSA_WALLETCONNECT.
// Implements frostfscrypto.Signer.
func (x SignerWalletConnect) Scheme() frostfscrypto.Scheme {
	return frostfscrypto.ECDSA_WALLETCONNECT
}

// Sign signs data using ECDSA algorithm with SHA-512 hashing.
// Implements frostfscrypto.Signer.
func (x SignerWalletConnect) Sign(data []byte) ([]byte, error) {
	b64 := make([]byte, base64.StdEncoding.EncodedLen(len(data)))
	base64.StdEncoding.Encode(b64, data)
	return walletconnect.Sign((*ecdsa.PrivateKey)(&x), b64)
}

// Public initializes PublicKey and returns it as frostfscrypto.PublicKey.
// Implements frostfscrypto.Signer.
func (x SignerWalletConnect) Public() frostfscrypto.PublicKey {
	return (*PublicKeyWalletConnect)(&x.PublicKey)
}

// PublicKeyWalletConnect is a wrapper over ecdsa.PublicKey used for NeoFS needs.
// Provides frostfscrypto.PublicKey interface.
//
// Instances MUST be initialized from ecdsa.PublicKey using type conversion.
type PublicKeyWalletConnect ecdsa.PublicKey

// MaxEncodedSize returns size of the compressed ECDSA public key.
func (x PublicKeyWalletConnect) MaxEncodedSize() int {
	return 33
}

// Encode encodes ECDSA public key in compressed form into buf.
// Uses exactly MaxEncodedSize bytes of the buf.
//
// Encode panics if buf length is less than MaxEncodedSize.
//
// See also Decode.
func (x PublicKeyWalletConnect) Encode(buf []byte) int {
	if len(buf) < 33 {
		panic(fmt.Sprintf("too short buffer %d", len(buf)))
	}

	return copy(buf, (*keys.PublicKey)(&x).Bytes())
}

// Decode decodes compressed binary representation of the PublicKeyWalletConnect.
//
// See also Encode.
func (x *PublicKeyWalletConnect) Decode(data []byte) error {
	pub, err := keys.NewPublicKeyFromBytes(data, elliptic.P256())
	if err != nil {
		return err
	}

	*x = (PublicKeyWalletConnect)(*pub)

	return nil
}

// Verify verifies data signature calculated by ECDSA algorithm with SHA-512 hashing.
func (x PublicKeyWalletConnect) Verify(data, signature []byte) bool {
	b64 := make([]byte, base64.StdEncoding.EncodedLen(len(data)))
	base64.StdEncoding.Encode(b64, data)
	return walletconnect.Verify((*ecdsa.PublicKey)(&x), b64, signature)
}
