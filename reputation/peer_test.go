package reputation_test

import (
	"testing"

	v2reputation "github.com/TrueCloudLab/frostfs-api-go/v2/reputation"
	"github.com/TrueCloudLab/frostfs-sdk-go/reputation"
	reputationtest "github.com/TrueCloudLab/frostfs-sdk-go/reputation/test"
	"github.com/stretchr/testify/require"
)

func TestPeerID_PublicKey(t *testing.T) {
	var val reputation.PeerID

	require.Zero(t, val.PublicKey())

	key := []byte{3, 2, 1}

	val.SetPublicKey(key)

	var m v2reputation.PeerID
	val.WriteToV2(&m)

	require.Equal(t, key, m.GetPublicKey())

	var val2 reputation.PeerID
	require.NoError(t, val2.ReadFromV2(m))

	require.Equal(t, key, val.PublicKey())

	require.True(t, reputation.ComparePeerKey(val, key))
}

func TestPeerID_EncodeToString(t *testing.T) {
	val := reputationtest.PeerID()
	var val2 reputation.PeerID

	require.NoError(t, val2.DecodeString(val.EncodeToString()))
	require.Equal(t, val, val2)
}
