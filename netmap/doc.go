/*
Package netmap provides functionality for working with information about the
NeoFS network, primarily a layer of storage nodes.

The package concentrates all the characteristics of NeoFS networks.

NetMap represents NeoFS network map - one of the main technologies used to
store data in the system. It is composed of information about all storage nodes
(NodeInfo type) in a particular network. NetMap methods allow you to impose
container storage policies (PlacementPolicy type) on a fixed composition of
nodes for selecting nodes corresponding to the placement rules chosen by the
container creator.

NetworkInfo type is dedicated to descriptive characterization of network state
and settings.

Instances can be also used to process NeoFS API V2 protocol messages
(see neo.fs.v2.netmap package in https://github.com/TrueCloudLab/frostfs-api).

On client side:

	import "github.com/TrueCloudLab/frostfs-api-go/v2/netmap"

	var msg netmap.NodeInfo
	info.WriteToV2(&msg)

	// send msg

On server side:

	// recv msg

	var info NodeInfo

	err := info.ReadFromV2(msg)
	// ...

	// process info

Using package types in an application is recommended to potentially work with
different protocol versions with which these types are compatible.
*/
package netmap
