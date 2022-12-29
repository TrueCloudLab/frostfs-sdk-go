/*
Package pool provides a wrapper for several FrostFS API clients.

The main component is Pool type. It is a virtual connection to the network
and provides methods for executing operations on the server. It also supports
a weighted random selection of the underlying client to make requests.

Create pool instance with 3 nodes connection.
This InitParameters will make pool use 192.168.130.71 node while it is healthy. Otherwise, it will make the pool use
192.168.130.72 for 90% of requests and 192.168.130.73 for remaining 10%.
:

	var prm pool.InitParameters
	prm.SetKey(key)
	prm.AddNode(NewNodeParam(1, "192.168.130.71", 1))
	prm.AddNode(NewNodeParam(2, "192.168.130.72", 9))
	prm.AddNode(NewNodeParam(2, "192.168.130.73", 1))
	// ...

	p, err := pool.NewPool(prm)
	// ...

Connect to the FrostFS server:

	err := p.Dial(ctx)
	// ...

Execute FrostFS operation on the server:

	var prm pool.PrmContainerPut
	prm.SetContainer(cnr)
	// ...

	res, err := p.PutContainer(context.Background(), prm)
	// ...

Execute FrostFS operation on the server and check error:

	var prm pool.PrmObjectHead
	prm.SetAddress(addr)
	// ...

	res, err := p.HeadObject(context.Background(), prm)
	if client.IsErrObjectNotFound(err) {
		// ...
	}
	// ...

Close the connection:

	p.Close()
*/
package pool
