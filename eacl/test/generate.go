package eacltest

import (
	cidtest "github.com/TrueCloudLab/frostfs-sdk-go/container/id/test"
	"github.com/TrueCloudLab/frostfs-sdk-go/eacl"
	usertest "github.com/TrueCloudLab/frostfs-sdk-go/user/test"
	versiontest "github.com/TrueCloudLab/frostfs-sdk-go/version/test"
)

// Target returns random eacl.Target.
func Target() *eacl.Target {
	x := eacl.NewTarget()

	x.SetRole(eacl.RoleSystem)
	x.SetBinaryKeys([][]byte{
		{1, 2, 3},
		{4, 5, 6},
	})

	return x
}

// Record returns random eacl.Record.
func Record() *eacl.Record {
	x := eacl.NewRecord()

	x.SetAction(eacl.ActionAllow)
	x.SetOperation(eacl.OperationRangeHash)
	x.SetTargets(*Target(), *Target())
	x.AddObjectContainerIDFilter(eacl.MatchStringEqual, cidtest.ID())
	x.AddObjectOwnerIDFilter(eacl.MatchStringNotEqual, usertest.ID())

	return x
}

func Table() *eacl.Table {
	x := eacl.NewTable()

	x.SetCID(cidtest.ID())
	x.AddRecord(Record())
	x.AddRecord(Record())
	x.SetVersion(versiontest.Version())

	return x
}
