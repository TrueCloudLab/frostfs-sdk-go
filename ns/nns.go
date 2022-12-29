package ns

import (
	"context"
	"errors"
	"fmt"
	"net/url"

	"github.com/TrueCloudLab/frostfs-contract/nns"
	"github.com/TrueCloudLab/frostfs-sdk-go/container"
	cid "github.com/TrueCloudLab/frostfs-sdk-go/container/id"
	"github.com/nspcc-dev/neo-go/pkg/core/state"
	"github.com/nspcc-dev/neo-go/pkg/neorpc/result"
	"github.com/nspcc-dev/neo-go/pkg/rpcclient"
	"github.com/nspcc-dev/neo-go/pkg/rpcclient/invoker"
	"github.com/nspcc-dev/neo-go/pkg/rpcclient/unwrap"
	"github.com/nspcc-dev/neo-go/pkg/util"
	"github.com/nspcc-dev/neo-go/pkg/vm/stackitem"
)

// NNS looks up FrostFS names using Neo Name Service.
//
// Instances are created with a variable declaration. Before work, the connection
// to the NNS server MUST be established using Dial method.
type NNS struct {
	nnsContract util.Uint160

	invoker interface {
		Call(contract util.Uint160, operation string, params ...interface{}) (*result.Invoke, error)
	}
}

// Dial connects to the address of the NNS server. If fails, the instance
// MUST NOT be used.
//
// If URL address scheme is 'ws' or 'wss', then WebSocket protocol is used,
// otherwise HTTP.
func (n *NNS) Dial(address string) error {
	// multiSchemeClient unites invoker.RPCInvoke and common interface of
	// rpcclient.Client and rpcclient.WSClient. Interface is anonymous
	// according to assumption that common interface of these client types
	// is not required by design and may diverge with changes.
	var multiSchemeClient interface {
		invoker.RPCInvoke
		// Init turns client to "ready-to-work" state.
		Init() error
		// GetContractStateByID returns state of the NNS contract on 1 input.
		GetContractStateByID(int32) (*state.Contract, error)
	}
	var err error

	uri, err := url.Parse(address)
	if err == nil && (uri.Scheme == "ws" || uri.Scheme == "wss") {
		multiSchemeClient, err = rpcclient.NewWS(context.Background(), address, rpcclient.Options{})
		if err != nil {
			return fmt.Errorf("create Neo WebSocket client: %w", err)
		}
	} else {
		multiSchemeClient, err = rpcclient.New(context.Background(), address, rpcclient.Options{})
		if err != nil {
			return fmt.Errorf("create Neo HTTP client: %w", err)
		}
	}

	if err = multiSchemeClient.Init(); err != nil {
		return fmt.Errorf("initialize Neo client: %w", err)
	}

	nnsContract, err := multiSchemeClient.GetContractStateByID(1)
	if err != nil {
		return fmt.Errorf("get NNS contract state: %w", err)
	}

	n.invoker = invoker.New(multiSchemeClient, nil)
	n.nnsContract = nnsContract.Hash

	return nil
}

// ResolveContainerDomain looks up for NNS TXT records for the given container domain
// by calling `resolve` method of NNS contract. Returns the first record which represents
// valid container ID in a string format. Otherwise, returns an error.
//
// ResolveContainerDomain MUST NOT be called before successful Dial.
//
// See also https://docs.neo.org/docs/en-us/reference/nns.html.
func (n *NNS) ResolveContainerDomain(domain container.Domain) (cid.ID, error) {
	item, err := unwrap.Item(n.invoker.Call(n.nnsContract, "resolve",
		domain.Name()+"."+domain.Zone(), int64(nns.TXT),
	))
	if err != nil {
		return cid.ID{}, fmt.Errorf("contract invocation: %w", err)
	}

	if _, ok := item.(stackitem.Null); !ok {
		arr, ok := item.Value().([]stackitem.Item)
		if !ok {
			// unexpected for types from stackitem package
			return cid.ID{}, errors.New("invalid cast to stack item slice")
		}

		var id cid.ID

		for i := range arr {
			bs, err := arr[i].TryBytes()
			if err != nil {
				return cid.ID{}, fmt.Errorf("convert array item to byte slice: %w", err)
			}

			err = id.DecodeString(string(bs))
			if err == nil {
				return id, nil
			}
		}
	}

	return cid.ID{}, errNotFound
}
