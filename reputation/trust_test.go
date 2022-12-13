package reputation_test

import (
	"testing"

	"github.com/TrueCloudLab/frostfs-api-go/v2/refs"
	v2reputation "github.com/TrueCloudLab/frostfs-api-go/v2/reputation"
	frostfsecdsa "github.com/TrueCloudLab/frostfs-sdk-go/crypto/ecdsa"
	"github.com/TrueCloudLab/frostfs-sdk-go/reputation"
	reputationtest "github.com/TrueCloudLab/frostfs-sdk-go/reputation/test"
	"github.com/TrueCloudLab/frostfs-sdk-go/version"
	"github.com/nspcc-dev/neo-go/pkg/crypto/keys"
	"github.com/stretchr/testify/require"
)

func TestTrust_Peer(t *testing.T) {
	var trust reputation.Trust

	require.Zero(t, trust.Peer())

	trust = reputationtest.Trust()

	peer := reputationtest.PeerID()

	trust.SetPeer(peer)

	var peerV2 v2reputation.PeerID
	peer.WriteToV2(&peerV2)

	var trustV2 v2reputation.Trust
	trust.WriteToV2(&trustV2)

	require.Equal(t, &peerV2, trustV2.GetPeer())

	var val2 reputation.Trust
	require.NoError(t, val2.ReadFromV2(trustV2))

	require.Equal(t, peer, val2.Peer())
}

func TestTrust_Value(t *testing.T) {
	var val reputation.Trust

	require.Zero(t, val.Value())

	val = reputationtest.Trust()

	const value = 0.75

	val.SetValue(value)

	var trustV2 v2reputation.Trust
	val.WriteToV2(&trustV2)

	require.EqualValues(t, value, trustV2.GetValue())

	var val2 reputation.Trust
	require.NoError(t, val2.ReadFromV2(trustV2))

	require.EqualValues(t, value, val2.Value())
}

func TestPeerToPeerTrust_TrustingPeer(t *testing.T) {
	var val reputation.PeerToPeerTrust

	require.Zero(t, val.TrustingPeer())

	val = reputationtest.PeerToPeerTrust()

	peer := reputationtest.PeerID()

	val.SetTrustingPeer(peer)

	var peerV2 v2reputation.PeerID
	peer.WriteToV2(&peerV2)

	var trustV2 v2reputation.PeerToPeerTrust
	val.WriteToV2(&trustV2)

	require.Equal(t, &peerV2, trustV2.GetTrustingPeer())

	var val2 reputation.PeerToPeerTrust
	require.NoError(t, val2.ReadFromV2(trustV2))

	require.Equal(t, peer, val2.TrustingPeer())
}

func TestPeerToPeerTrust_Trust(t *testing.T) {
	var val reputation.PeerToPeerTrust

	require.Zero(t, val.Trust())

	val = reputationtest.PeerToPeerTrust()

	trust := reputationtest.Trust()

	val.SetTrust(trust)

	var trustV2 v2reputation.Trust
	trust.WriteToV2(&trustV2)

	var valV2 v2reputation.PeerToPeerTrust
	val.WriteToV2(&valV2)

	require.Equal(t, &trustV2, valV2.GetTrust())

	var val2 reputation.PeerToPeerTrust
	require.NoError(t, val2.ReadFromV2(valV2))

	require.Equal(t, trust, val2.Trust())
}

func TestGlobalTrust_Init(t *testing.T) {
	var val reputation.GlobalTrust
	val.Init()

	var valV2 v2reputation.GlobalTrust
	val.WriteToV2(&valV2)

	var verV2 refs.Version
	version.Current().WriteToV2(&verV2)

	require.Equal(t, &verV2, valV2.GetVersion())
}

func TestGlobalTrust_Manager(t *testing.T) {
	var val reputation.GlobalTrust

	require.Zero(t, val.Manager())

	val = reputationtest.SignedGlobalTrust()

	peer := reputationtest.PeerID()

	val.SetManager(peer)

	var peerV2 v2reputation.PeerID
	peer.WriteToV2(&peerV2)

	var trustV2 v2reputation.GlobalTrust
	val.WriteToV2(&trustV2)

	require.Equal(t, &peerV2, trustV2.GetBody().GetManager())

	var val2 reputation.GlobalTrust
	require.NoError(t, val2.ReadFromV2(trustV2))

	require.Equal(t, peer, val2.Manager())
}

func TestGlobalTrust_Trust(t *testing.T) {
	var val reputation.GlobalTrust

	require.Zero(t, val.Trust())

	val = reputationtest.SignedGlobalTrust()

	trust := reputationtest.Trust()

	val.SetTrust(trust)

	var trustV2 v2reputation.Trust
	trust.WriteToV2(&trustV2)

	var valV2 v2reputation.GlobalTrust
	val.WriteToV2(&valV2)

	require.Equal(t, &trustV2, valV2.GetBody().GetTrust())

	var val2 reputation.GlobalTrust
	require.NoError(t, val2.ReadFromV2(valV2))

	require.Equal(t, trust, val2.Trust())
}

func TestGlobalTrust_Sign(t *testing.T) {
	k, err := keys.NewPrivateKey()
	require.NoError(t, err)

	val := reputationtest.GlobalTrust()

	require.False(t, val.VerifySignature())

	require.NoError(t, val.Sign(frostfsecdsa.Signer(k.PrivateKey)))

	var valV2 v2reputation.GlobalTrust
	val.WriteToV2(&valV2)

	require.NotZero(t, valV2.GetSignature())

	var val2 reputation.GlobalTrust
	require.NoError(t, val2.ReadFromV2(valV2))

	require.True(t, val2.VerifySignature())
}

func TestGlobalTrustEncoding(t *testing.T) {
	val := reputationtest.SignedGlobalTrust()

	t.Run("binary", func(t *testing.T) {
		data := val.Marshal()

		var val2 reputation.GlobalTrust
		require.NoError(t, val2.Unmarshal(data))

		require.Equal(t, val, val2)
	})
}
