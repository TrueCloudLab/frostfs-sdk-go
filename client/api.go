package client

import (
	"context"
	"fmt"

	v2netmap "github.com/TrueCloudLab/frostfs-api-go/v2/netmap"
	rpcapi "github.com/TrueCloudLab/frostfs-api-go/v2/rpc"
	"github.com/TrueCloudLab/frostfs-api-go/v2/rpc/client"
)

// interface of FrostFS API server. Exists for test purposes only.
type frostFSAPIServer interface {
	netMapSnapshot(context.Context, v2netmap.SnapshotRequest) (*v2netmap.SnapshotResponse, error)
}

// wrapper over real client connection which communicates over FrostFS API protocol.
// Provides frostFSAPIServer for Client instances used in real applications.
type coreServer client.Client

// unifies errors of all RPC.
func rpcErr(e error) error {
	return fmt.Errorf("rpc failure: %w", e)
}

// executes NetmapService.NetmapSnapshot RPC declared in FrostFS API protocol
// using underlying client.Client.
func (x *coreServer) netMapSnapshot(ctx context.Context, req v2netmap.SnapshotRequest) (*v2netmap.SnapshotResponse, error) {
	resp, err := rpcapi.NetMapSnapshot((*client.Client)(x), &req, client.WithContext(ctx))
	if err != nil {
		return nil, rpcErr(err)
	}

	return resp, nil
}
