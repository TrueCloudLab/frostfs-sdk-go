package client

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"testing"

	apistatus "github.com/TrueCloudLab/frostfs-sdk-go/client/status"
	"github.com/stretchr/testify/require"
)

/*
File contains common functionality used for client package testing.
*/

var key, _ = ecdsa.GenerateKey(elliptic.P256(), rand.Reader)

var statusErr apistatus.ServerInternal

func init() {
	statusErr.SetMessage("test status error")
}

func assertStatusErr(tb testing.TB, res interface{ Status() apistatus.Status }) {
	require.IsType(tb, &statusErr, res.Status())
	require.Equal(tb, statusErr.Message(), res.Status().(*apistatus.ServerInternal).Message())
}

func newClient(server neoFSAPIServer) *Client {
	var prm PrmInit
	prm.SetDefaultPrivateKey(*key)

	var c Client
	c.Init(prm)
	c.setNeoFSAPIServer(server)

	return &c
}

func TestClient_DialContext(t *testing.T) {
	var c Client

	// try to connect to any host
	var prm PrmDial
	prm.SetServerURI("localhost:8080")

	assert := func(ctx context.Context, errExpected error) {
		// use the particular context
		prm.SetContext(ctx)

		// expect particular context error according to Dial docs
		require.ErrorIs(t, c.Dial(prm), errExpected)
	}

	// create pre-abandoned context
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	assert(ctx, context.Canceled)

	// create "pre-deadlined" context
	ctx, cancel = context.WithTimeout(context.Background(), 0)
	defer cancel()

	assert(ctx, context.DeadlineExceeded)
}
