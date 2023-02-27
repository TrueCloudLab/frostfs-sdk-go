package client

import (
	"context"
	"errors"
	"fmt"
	"testing"

	v2netmap "github.com/TrueCloudLab/frostfs-api-go/v2/netmap"
	"github.com/TrueCloudLab/frostfs-api-go/v2/session"
	"github.com/TrueCloudLab/frostfs-api-go/v2/signature"
	apistatus "github.com/TrueCloudLab/frostfs-sdk-go/client/status"
	"github.com/TrueCloudLab/frostfs-sdk-go/netmap"
	"github.com/stretchr/testify/require"
)

type serverNetMap struct {
	errTransport error

	signResponse bool

	statusOK bool

	setNetMap bool
	netMap    v2netmap.NetMap
}

func (x *serverNetMap) netMapSnapshot(ctx context.Context, req v2netmap.SnapshotRequest) (*v2netmap.SnapshotResponse, error) {
	err := signature.VerifyServiceMessage(&req)
	if err != nil {
		return nil, err
	}

	if x.errTransport != nil {
		return nil, x.errTransport
	}

	var body v2netmap.SnapshotResponseBody

	if x.setNetMap {
		body.SetNetMap(&x.netMap)
	}

	var meta session.ResponseMetaHeader

	if !x.statusOK {
		meta.SetStatus(statusErr.ToStatusV2())
	}

	var resp v2netmap.SnapshotResponse
	resp.SetBody(&body)
	resp.SetMetaHeader(&meta)

	if x.signResponse {
		err = signature.SignServiceMessage(key, &resp)
		if err != nil {
			panic(fmt.Sprintf("sign response: %v", err))
		}
	}

	return &resp, nil
}

func TestClient_NetMapSnapshot(t *testing.T) {
	var err error
	var prm PrmNetMapSnapshot
	var res *ResNetMapSnapshot
	var srv serverNetMap
	c := newClient(&srv)
	ctx := context.Background()

	// missing context
	//nolint:staticcheck
	_, err = c.NetMapSnapshot(nil, prm)
	require.ErrorIs(t, err, errorMissingContext, "")

	// request signature
	srv.errTransport = errors.New("any error")

	_, err = c.NetMapSnapshot(ctx, prm)
	require.ErrorIs(t, err, srv.errTransport)

	srv.errTransport = nil

	// unsigned response
	_, err = c.NetMapSnapshot(ctx, prm)
	require.Error(t, err)

	srv.signResponse = true

	// status failure
	res, err = c.NetMapSnapshot(ctx, prm)
	require.NoError(t, err)
	assertStatusErr(t, res)

	srv.statusOK = true

	// missing netmap field
	_, err = c.NetMapSnapshot(ctx, prm)
	require.Error(t, err)

	srv.setNetMap = true

	// invalid network map
	var netMap netmap.NetMap

	var node netmap.NodeInfo
	// TODO: #260 use instance corrupter

	var nodeV2 v2netmap.NodeInfo

	node.WriteToV2(&nodeV2)
	require.Error(t, new(netmap.NodeInfo).ReadFromV2(nodeV2))

	netMap.SetNodes([]netmap.NodeInfo{node})
	netMap.WriteToV2(&srv.netMap)

	_, err = c.NetMapSnapshot(ctx, prm)
	require.Error(t, err)

	// correct network map
	// TODO: #260 use instance normalizer
	node.SetPublicKey([]byte{1, 2, 3})
	node.SetNetworkEndpoints("1", "2", "3")

	node.WriteToV2(&nodeV2)
	require.NoError(t, new(netmap.NodeInfo).ReadFromV2(nodeV2))

	netMap.SetNodes([]netmap.NodeInfo{node})
	netMap.WriteToV2(&srv.netMap)

	res, err = c.NetMapSnapshot(ctx, prm)
	require.NoError(t, err)
	require.True(t, apistatus.IsSuccessful(res.Status()))
	require.Equal(t, netMap, res.NetMap())
}
