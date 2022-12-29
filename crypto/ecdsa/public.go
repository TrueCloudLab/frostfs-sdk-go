package frostfsecdsa

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/sha256"
	"crypto/sha512"
	"fmt"
	"math/big"

	"github.com/nspcc-dev/neo-go/pkg/crypto/keys"
)

// PublicKey is a wrapper over ecdsa.PublicKey used for FrostFS needs.
// Provides frostfscrypto.PublicKey interface.
//
// Instances MUST be initialized from ecdsa.PublicKey using type conversion.
type PublicKey ecdsa.PublicKey

// MaxEncodedSize returns size of the compressed ECDSA public key.
func (x PublicKey) MaxEncodedSize() int {
	return 33
}

// Encode encodes ECDSA public key in compressed form into buf.
// Uses exactly MaxEncodedSize bytes of the buf.
//
// Encode panics if buf length is less than MaxEncodedSize.
//
// See also Decode.
func (x PublicKey) Encode(buf []byte) int {
	if len(buf) < 33 {
		panic(fmt.Sprintf("too short buffer %d", len(buf)))
	}

	return copy(buf, (*keys.PublicKey)(&x).Bytes())
}

// Decode decodes compressed binary representation of the PublicKey.
//
// See also Encode.
func (x *PublicKey) Decode(data []byte) error {
	pub, err := keys.NewPublicKeyFromBytes(data, elliptic.P256())
	if err != nil {
		return err
	}

	*x = (PublicKey)(*pub)

	return nil
}

// similar to elliptic.Unmarshal but without IsOnCurve check.
func unmarshalXY(data []byte) (x *big.Int, y *big.Int) {
	if len(data) != 65 {
		return
	} else if data[0] != 4 { // uncompressed form
		return
	}

	p := elliptic.P256().Params().P
	x = new(big.Int).SetBytes(data[1:33])
	y = new(big.Int).SetBytes(data[33:])

	if x.Cmp(p) >= 0 || y.Cmp(p) >= 0 {
		x, y = nil, nil
	}

	return
}

// Verify verifies data signature calculated by ECDSA algorithm with SHA-512 hashing.
func (x PublicKey) Verify(data, signature []byte) bool {
	h := sha512.Sum512(data)
	r, s := unmarshalXY(signature)

	return r != nil && s != nil && ecdsa.Verify((*ecdsa.PublicKey)(&x), h[:], r, s)
}

// PublicKeyRFC6979 is a wrapper over ecdsa.PublicKey used for FrostFS needs.
// Provides frostfscrypto.PublicKey interface.
//
// Instances MUST be initialized from ecdsa.PublicKey using type conversion.
type PublicKeyRFC6979 ecdsa.PublicKey

// MaxEncodedSize returns size of the compressed ECDSA public key.
func (x PublicKeyRFC6979) MaxEncodedSize() int {
	return 33
}

// Encode encodes ECDSA public key in compressed form into buf.
// Uses exactly MaxEncodedSize bytes of the buf.
//
// Encode panics if buf length is less than MaxEncodedSize.
//
// See also Decode.
func (x PublicKeyRFC6979) Encode(buf []byte) int {
	if len(buf) < 33 {
		panic(fmt.Sprintf("too short buffer %d", len(buf)))
	}

	return copy(buf, (*keys.PublicKey)(&x).Bytes())
}

// Decode decodes binary representation of the ECDSA public key.
//
// See also Encode.
func (x *PublicKeyRFC6979) Decode(data []byte) error {
	pub, err := keys.NewPublicKeyFromBytes(data, elliptic.P256())
	if err != nil {
		return err
	}

	*x = (PublicKeyRFC6979)(*pub)

	return nil
}

// Verify verifies data signature calculated by deterministic ECDSA algorithm
// with SHA-256 hashing.
//
// See also RFC 6979.
func (x PublicKeyRFC6979) Verify(data, signature []byte) bool {
	h := sha256.Sum256(data)
	return (*keys.PublicKey)(&x).Verify(signature, h[:])
}
