/*
Package client provides FrostFS API client implementation.

The main component is Client type. It is a virtual connection to the network
and provides methods for executing operations on the server.

Create client instance:

	var c client.Client

Initialize client state:

	var prm client.PrmInit
	prm.SetDefaultPrivateKey(key)
	// ...

	c.Init(prm)

Connect to the FrostFS server:

	var prm client.PrmDial
	prm.SetServerURI("localhost:8080")
	prm.SetDefaultPrivateKey(key)
	// ...

	err := c.Dial(prm)
	// ...

Execute FrostFS operation on the server:

	var prm client.PrmContainerPut
	prm.SetContainer(cnr)
	// ...

	res, err := c.ContainerPut(context.Background(), prm)
	err := c.Dial(dialPrm)
	if err == nil {
		err = apistatus.ErrFromStatus(res.Status())
	}
	// ...

Consume custom service of the server:

	syntax = "proto3";

	service CustomService {
		rpc CustomRPC(CustomRPCRequest) returns (CustomRPCResponse);
	}

	import "github.com/TrueCloudLab/frostfs-api-go/v2/rpc/client"
	import "github.com/TrueCloudLab/frostfs-api-go/v2/rpc/common"

	req := new(CustomRPCRequest)
	// ...
	resp := new(CustomRPCResponse)

	err := c.ExecRaw(func(c *client.Client) error {
		return client.SendUnary(c, common.CallMethodInfo{
			Service: "CustomService",
			Name:    "CustomRPC",
		}, req, resp)
	})
	// ...

Close the connection:

	err := c.Close()
	// ...

Note that it's not allowed to override Client behaviour directly: the parameters
for the all operations are write-only and the results of the all operations are
read-only. To be able to override client behavior (e.g. for tests), abstract it
with an interface:

	import "github.com/TrueCloudLab/frostfs-sdk-go/client"

	type FrostFSClient interface {
		// Operations according to the application needs
		CreateContainer(context.Context, container.Container) error
		// ...
	}

	type client struct {
		c *client.Client
	}

	func (x *client) CreateContainer(context.Context, container.Container) error {
		// ...
	}
*/
package client
