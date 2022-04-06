package pool

import (
	"testing"
	"time"

	"github.com/nspcc-dev/neo-go/pkg/crypto/keys"
	"github.com/nspcc-dev/neofs-sdk-go/session"
	sessiontest "github.com/nspcc-dev/neofs-sdk-go/session/test"
	"github.com/stretchr/testify/require"
)

func TestSessionCache_GetAccessTime(t *testing.T) {
	const key = "Foo"

	cache, err := NewCache()
	require.NoError(t, err)

	cache.Put(key, nil)

	t1, ok := cache.GetAccessTime(key)
	require.True(t, ok)

	time.Sleep(10 * time.Millisecond)

	t2, ok := cache.GetAccessTime(key)
	require.True(t, ok)
	require.Equal(t, t1, t2)

	_ = cache.Get(key)
	t3, ok := cache.GetAccessTime(key)
	require.True(t, ok)
	require.NotEqual(t, t1, t3)
}

func TestSessionCache_GetUnmodifiedToken(t *testing.T) {
	const key = "Foo"
	target := sessiontest.Token()

	pk, err := keys.NewPrivateKey()
	require.NoError(t, err)

	check := func(t *testing.T, tok *session.Token, extra string) {
		require.Empty(t, tok.Signature().Sign(), extra)
		require.Nil(t, tok.Context(), extra)
	}

	cache, err := NewCache()
	require.NoError(t, err)

	cache.Put(key, target)
	value := cache.Get(key)
	check(t, value, "before sign")

	err = value.Sign(&pk.PrivateKey)
	require.NoError(t, err)

	octx := sessiontest.ObjectContext()
	value.SetContext(octx)

	value = cache.Get(key)
	check(t, value, "after sign")
}
