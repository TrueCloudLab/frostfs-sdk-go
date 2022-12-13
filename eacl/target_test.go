package eacl

import (
	"crypto/ecdsa"
	"testing"

	"github.com/TrueCloudLab/frostfs-api-go/v2/acl"
	v2acl "github.com/TrueCloudLab/frostfs-api-go/v2/acl"
	"github.com/nspcc-dev/neo-go/pkg/crypto/keys"
	"github.com/stretchr/testify/require"
)

func TestTarget(t *testing.T) {
	pubs := []*ecdsa.PublicKey{
		randomPublicKey(t),
		randomPublicKey(t),
	}

	target := NewTarget()
	target.SetRole(RoleSystem)
	SetTargetECDSAKeys(target, pubs...)

	v2 := target.ToV2()
	require.NotNil(t, v2)
	require.Equal(t, v2acl.RoleSystem, v2.GetRole())
	require.Len(t, v2.GetKeys(), len(pubs))
	for i, key := range v2.GetKeys() {
		require.Equal(t, key, (*keys.PublicKey)(pubs[i]).Bytes())
	}

	newTarget := NewTargetFromV2(v2)
	require.Equal(t, target, newTarget)

	t.Run("from nil v2 target", func(t *testing.T) {
		require.Equal(t, new(Target), NewTargetFromV2(nil))
	})
}

func TestTargetEncoding(t *testing.T) {
	tar := NewTarget()
	tar.SetRole(RoleSystem)
	SetTargetECDSAKeys(tar, randomPublicKey(t))

	t.Run("binary", func(t *testing.T) {
		data, err := tar.Marshal()
		require.NoError(t, err)

		tar2 := NewTarget()
		require.NoError(t, tar2.Unmarshal(data))

		require.Equal(t, tar, tar2)
	})

	t.Run("json", func(t *testing.T) {
		data, err := tar.MarshalJSON()
		require.NoError(t, err)

		tar2 := NewTarget()
		require.NoError(t, tar2.UnmarshalJSON(data))

		require.Equal(t, tar, tar2)
	})
}

func TestTarget_ToV2(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		var x *Target

		require.Nil(t, x.ToV2())
	})

	t.Run("default values", func(t *testing.T) {
		target := NewTarget()

		// check initial values
		require.Equal(t, RoleUnknown, target.Role())
		require.Nil(t, target.BinaryKeys())

		// convert to v2 message
		targetV2 := target.ToV2()

		require.Equal(t, acl.RoleUnknown, targetV2.GetRole())
		require.Nil(t, targetV2.GetKeys())
	})
}
