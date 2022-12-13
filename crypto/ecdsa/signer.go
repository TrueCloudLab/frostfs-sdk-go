package frostfsecdsa

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha512"

	frostfscrypto "github.com/TrueCloudLab/frostfs-sdk-go/crypto"
	"github.com/nspcc-dev/neo-go/pkg/crypto/keys"
)

// Signer wraps ecdsa.PrivateKey and represents signer based on ECDSA with
// SHA-512 hashing. Provides frostfscrypto.Signer interface.
//
// Instances MUST be initialized from ecdsa.PrivateKey using type conversion.
type Signer ecdsa.PrivateKey

// Scheme returns frostfscrypto.ECDSA_SHA512.
// Implements frostfscrypto.Signer.
func (x Signer) Scheme() frostfscrypto.Scheme {
	return frostfscrypto.ECDSA_SHA512
}

// Sign signs data using ECDSA algorithm with SHA-512 hashing.
// Implements frostfscrypto.Signer.
func (x Signer) Sign(data []byte) ([]byte, error) {
	h := sha512.Sum512(data)
	r, s, err := ecdsa.Sign(rand.Reader, (*ecdsa.PrivateKey)(&x), h[:])
	if err != nil {
		return nil, err
	}

	params := elliptic.P256().Params()
	curveOrderByteSize := params.P.BitLen() / 8

	buf := make([]byte, 1+curveOrderByteSize*2)
	buf[0] = 4

	_ = r.FillBytes(buf[1 : 1+curveOrderByteSize])
	_ = s.FillBytes(buf[1+curveOrderByteSize:])

	return buf, nil
}

// Public initializes PublicKey and returns it as frostfscrypto.PublicKey.
// Implements frostfscrypto.Signer.
func (x Signer) Public() frostfscrypto.PublicKey {
	return (*PublicKey)(&x.PublicKey)
}

// SignerRFC6979 wraps ecdsa.PrivateKey and represents signer based on deterministic
// ECDSA with SHA-256 hashing (RFC 6979). Provides frostfscrypto.Signer interface.
//
// Instances SHOULD be initialized from ecdsa.PrivateKey using type conversion.
type SignerRFC6979 ecdsa.PrivateKey

// Scheme returns frostfscrypto.ECDSA_DETERMINISTIC_SHA256.
// Implements frostfscrypto.Signer.
func (x SignerRFC6979) Scheme() frostfscrypto.Scheme {
	return frostfscrypto.ECDSA_DETERMINISTIC_SHA256
}

// Sign signs data using deterministic ECDSA algorithm with SHA-256 hashing.
// Implements frostfscrypto.Signer.
//
// See also RFC 6979.
func (x SignerRFC6979) Sign(data []byte) ([]byte, error) {
	p := keys.PrivateKey{PrivateKey: (ecdsa.PrivateKey)(x)}
	return p.Sign(data), nil
}

// Public initializes PublicKeyRFC6979 and returns it as frostfscrypto.PublicKey.
// Implements frostfscrypto.Signer.
func (x SignerRFC6979) Public() frostfscrypto.PublicKey {
	return (*PublicKeyRFC6979)(&x.PublicKey)
}
