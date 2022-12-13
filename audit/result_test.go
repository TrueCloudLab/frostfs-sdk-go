package audit_test

import (
	"bytes"
	"testing"

	"github.com/TrueCloudLab/frostfs-sdk-go/audit"
	audittest "github.com/TrueCloudLab/frostfs-sdk-go/audit/test"
	cidtest "github.com/TrueCloudLab/frostfs-sdk-go/container/id/test"
	oid "github.com/TrueCloudLab/frostfs-sdk-go/object/id"
	oidtest "github.com/TrueCloudLab/frostfs-sdk-go/object/id/test"
	"github.com/stretchr/testify/require"
)

func TestResultData(t *testing.T) {
	var r audit.Result

	countSG := func(passed bool, f func(oid.ID)) int {
		called := 0

		ff := func(arg oid.ID) bool {
			called++

			if f != nil {
				f(arg)
			}

			return true
		}

		if passed {
			r.IteratePassedStorageGroups(ff)
		} else {
			r.IterateFailedStorageGroups(ff)
		}

		return called
	}

	countPassSG := func(f func(oid.ID)) int { return countSG(true, f) }
	countFailSG := func(f func(oid.ID)) int { return countSG(false, f) }

	countNodes := func(passed bool, f func([]byte)) int {
		called := 0

		ff := func(arg []byte) bool {
			called++

			if f != nil {
				f(arg)
			}

			return true
		}

		if passed {
			r.IteratePassedStorageNodes(ff)
		} else {
			r.IterateFailedStorageNodes(ff)
		}

		return called
	}

	countPassNodes := func(f func([]byte)) int { return countNodes(true, f) }
	countFailNodes := func(f func([]byte)) int { return countNodes(false, f) }

	require.Zero(t, r.Epoch())
	_, set := r.Container()
	require.False(t, set)
	require.Nil(t, r.AuditorKey())
	require.False(t, r.Completed())
	require.Zero(t, r.RequestsPoR())
	require.Zero(t, r.RetriesPoR())
	require.Zero(t, countPassSG(nil))
	require.Zero(t, countFailSG(nil))
	require.Zero(t, countPassNodes(nil))
	require.Zero(t, countFailNodes(nil))

	epoch := uint64(13)
	r.ForEpoch(epoch)
	require.Equal(t, epoch, r.Epoch())

	cnr := cidtest.ID()
	r.ForContainer(cnr)
	cID, set := r.Container()
	require.True(t, set)
	require.Equal(t, cnr, cID)

	key := []byte{1, 2, 3}
	r.SetAuditorKey(key)
	require.Equal(t, key, r.AuditorKey())

	r.Complete()
	require.True(t, r.Completed())

	requests := uint32(2)
	r.SetRequestsPoR(requests)
	require.Equal(t, requests, r.RequestsPoR())

	retries := uint32(1)
	r.SetRetriesPoR(retries)
	require.Equal(t, retries, r.RetriesPoR())

	passSG1, passSG2 := oidtest.ID(), oidtest.ID()
	r.SubmitPassedStorageGroup(passSG1)
	r.SubmitPassedStorageGroup(passSG2)

	called1, called2 := false, false

	require.EqualValues(t, 2, countPassSG(func(id oid.ID) {
		if id.Equals(passSG1) {
			called1 = true
		} else if id.Equals(passSG2) {
			called2 = true
		}
	}))
	require.True(t, called1)
	require.True(t, called2)

	failSG1, failSG2 := oidtest.ID(), oidtest.ID()
	r.SubmitFailedStorageGroup(failSG1)
	r.SubmitFailedStorageGroup(failSG2)

	called1, called2 = false, false

	require.EqualValues(t, 2, countFailSG(func(id oid.ID) {
		if id.Equals(failSG1) {
			called1 = true
		} else if id.Equals(failSG2) {
			called2 = true
		}
	}))
	require.True(t, called1)
	require.True(t, called2)

	hit := uint32(1)
	r.SetHits(hit)
	require.Equal(t, hit, r.Hits())

	miss := uint32(2)
	r.SetMisses(miss)
	require.Equal(t, miss, r.Misses())

	fail := uint32(3)
	r.SetFailures(fail)
	require.Equal(t, fail, r.Failures())

	passNodes := [][]byte{{1}, {2}}
	r.SubmitPassedStorageNodes(passNodes)

	called1, called2 = false, false

	require.EqualValues(t, 2, countPassNodes(func(arg []byte) {
		if bytes.Equal(arg, passNodes[0]) {
			called1 = true
		} else if bytes.Equal(arg, passNodes[1]) {
			called2 = true
		}
	}))
	require.True(t, called1)
	require.True(t, called2)

	failNodes := [][]byte{{3}, {4}}
	r.SubmitFailedStorageNodes(failNodes)

	called1, called2 = false, false

	require.EqualValues(t, 2, countFailNodes(func(arg []byte) {
		if bytes.Equal(arg, failNodes[0]) {
			called1 = true
		} else if bytes.Equal(arg, failNodes[1]) {
			called2 = true
		}
	}))
	require.True(t, called1)
	require.True(t, called2)
}

func TestResultEncoding(t *testing.T) {
	r := *audittest.Result()

	t.Run("binary", func(t *testing.T) {
		data := r.Marshal()

		var r2 audit.Result
		require.NoError(t, r2.Unmarshal(data))

		require.Equal(t, r, r2)
	})
}
