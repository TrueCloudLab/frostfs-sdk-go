package frostfscrypto_test

import (
	"math/rand"
	"testing"

	"github.com/TrueCloudLab/frostfs-api-go/v2/refs"
	frostfscrypto "github.com/TrueCloudLab/frostfs-sdk-go/crypto"
	frostfsecdsa "github.com/TrueCloudLab/frostfs-sdk-go/crypto/ecdsa"
	"github.com/nspcc-dev/neo-go/pkg/crypto/keys"
	"github.com/stretchr/testify/require"
)

func TestSignature(t *testing.T) {
	data := make([]byte, 512)
	rand.Read(data)

	k, err := keys.NewPrivateKey()
	require.NoError(t, err)

	var s frostfscrypto.Signature
	var m refs.Signature

	for _, f := range []func() frostfscrypto.Signer{
		func() frostfscrypto.Signer {
			return frostfsecdsa.Signer(k.PrivateKey)
		},
		func() frostfscrypto.Signer {
			return frostfsecdsa.SignerRFC6979(k.PrivateKey)
		},
		func() frostfscrypto.Signer {
			return frostfsecdsa.SignerWalletConnect(k.PrivateKey)
		},
	} {
		signer := f()

		err := s.Calculate(signer, data)
		require.NoError(t, err)

		s.WriteToV2(&m)

		require.NoError(t, s.ReadFromV2(m))

		valid := s.Verify(data)
		require.True(t, valid, "type %T", signer)
	}
}
