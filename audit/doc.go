/*
Package audit provides features to process data audit in FrostFS system.

Result type groups values which can be gathered during data audit process:

	var res audit.Result
	res.ForEpoch(32)
	res.ForContainer(cnr)
	// ...
	res.Complete()

Result instances can be stored in a binary format. On reporter side:

	data := res.Marshal()
	// send data

On receiver side:

	var res audit.Result
	err := res.Unmarshal(data)
	// ...

Using package types in an application is recommended to potentially work with
different protocol versions with which these types are compatible.
*/
package audit
