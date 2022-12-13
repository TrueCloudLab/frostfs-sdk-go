package audittest

import (
	"github.com/TrueCloudLab/frostfs-sdk-go/audit"
	cidtest "github.com/TrueCloudLab/frostfs-sdk-go/container/id/test"
	oidtest "github.com/TrueCloudLab/frostfs-sdk-go/object/id/test"
)

// Result returns random audit.Result.
func Result() *audit.Result {
	var x audit.Result

	x.ForContainer(cidtest.ID())
	x.SetAuditorKey([]byte("key"))
	x.Complete()
	x.ForEpoch(44)
	x.SetHits(55)
	x.SetMisses(66)
	x.SetFailures(77)
	x.SetRequestsPoR(88)
	x.SetRequestsPoR(99)
	x.SubmitFailedStorageNodes([][]byte{
		[]byte("node1"),
		[]byte("node2"),
	})
	x.SubmitPassedStorageNodes([][]byte{
		[]byte("node3"),
		[]byte("node4"),
	})
	x.SubmitPassedStorageGroup(oidtest.ID())
	x.SubmitPassedStorageGroup(oidtest.ID())
	x.SubmitFailedStorageGroup(oidtest.ID())
	x.SubmitFailedStorageGroup(oidtest.ID())

	return &x
}
