package client

import (
	"crypto/ecdsa"
	"errors"
	"fmt"
	"io"
	"testing"

	v2object "github.com/TrueCloudLab/frostfs-api-go/v2/object"
	"github.com/TrueCloudLab/frostfs-api-go/v2/refs"
	signatureV2 "github.com/TrueCloudLab/frostfs-api-go/v2/signature"
	oid "github.com/TrueCloudLab/frostfs-sdk-go/object/id"
	oidtest "github.com/TrueCloudLab/frostfs-sdk-go/object/id/test"
	"github.com/nspcc-dev/neo-go/pkg/crypto/keys"
	"github.com/stretchr/testify/require"
)

func TestObjectSearch(t *testing.T) {
	ids := make([]oid.ID, 20)
	for i := range ids {
		ids[i] = oidtest.ID()
	}

	p, resp := testListReaderResponse(t)

	buf := make([]oid.ID, 2)
	checkRead := func(t *testing.T, expected []oid.ID) {
		n, ok := resp.Read(buf)
		require.True(t, ok == (len(expected) == len(buf)), "expected no error")
		require.Equal(t, len(expected), n, "expected %d items to be read", len(expected))
		require.Equal(t, expected, buf[:len(expected)])
	}

	// nil panic
	require.Panics(t, func() { resp.Read(nil) })

	// both ID fetched
	resp.stream = newSearchStream(p, nil, ids[:3])
	checkRead(t, ids[:2])

	// one ID cached, second fetched
	resp.stream = newSearchStream(p, nil, ids[3:6])
	checkRead(t, ids[2:4])

	// both ID cached
	resp.stream = nil // shouldn't be called, panic if so
	checkRead(t, ids[4:6])

	// both ID fetched in 2 requests, with empty one in the middle
	resp.stream = newSearchStream(p, nil, ids[6:7], nil, ids[7:8])
	checkRead(t, ids[6:8])

	// read from tail multiple times
	resp.stream = newSearchStream(p, nil, ids[8:11])
	buf = buf[:1]
	checkRead(t, ids[8:9])
	checkRead(t, ids[9:10])
	checkRead(t, ids[10:11])

	// handle EOF
	buf = buf[:2]
	resp.stream = newSearchStream(p, io.EOF, ids[11:12])
	checkRead(t, ids[11:12])
}

func TestObjectIterate(t *testing.T) {
	ids := make([]oid.ID, 3)
	for i := range ids {
		ids[i] = oidtest.ID()
	}

	t.Run("iterate all sequence", func(t *testing.T) {
		p, resp := testListReaderResponse(t)

		resp.stream = newSearchStream(p, io.EOF, ids[0:2], nil, ids[2:3])

		var actual []oid.ID
		require.NoError(t, resp.Iterate(func(id oid.ID) bool {
			actual = append(actual, id)
			return false
		}))
		require.Equal(t, ids[:3], actual)
	})
	t.Run("stop by return value", func(t *testing.T) {
		p, resp := testListReaderResponse(t)

		var actual []oid.ID
		resp.stream = &singleStreamResponder{key: p, idList: [][]oid.ID{ids}}
		require.NoError(t, resp.Iterate(func(id oid.ID) bool {
			actual = append(actual, id)
			return len(actual) == 2
		}))
		require.Equal(t, ids[:2], actual)
	})
	t.Run("stop after error", func(t *testing.T) {
		p, resp := testListReaderResponse(t)
		expectedErr := errors.New("test error")

		resp.stream = newSearchStream(p, expectedErr, ids[:2])

		var actual []oid.ID
		err := resp.Iterate(func(id oid.ID) bool {
			actual = append(actual, id)
			return false
		})
		require.True(t, errors.Is(err, expectedErr), "got: %v", err)
		require.Equal(t, ids[:2], actual)
	})
}

func testListReaderResponse(t *testing.T) (*ecdsa.PrivateKey, *ObjectListReader) {
	p, err := keys.NewPrivateKey()
	require.NoError(t, err)

	return &p.PrivateKey, &ObjectListReader{
		cancelCtxStream: func() {},
		client:          &Client{},
		tail:            nil,
	}
}

func newSearchStream(key *ecdsa.PrivateKey, endError error, idList ...[]oid.ID) *singleStreamResponder {
	return &singleStreamResponder{
		key:      key,
		endError: endError,
		idList:   idList,
	}
}

type singleStreamResponder struct {
	key      *ecdsa.PrivateKey
	n        int
	endError error
	idList   [][]oid.ID
}

func (s *singleStreamResponder) Read(resp *v2object.SearchResponse) error {
	if s.n >= len(s.idList) {
		if s.endError != nil {
			return s.endError
		}
		panic("unexpected call to `Read`")
	}

	var body v2object.SearchResponseBody

	if s.idList[s.n] != nil {
		ids := make([]refs.ObjectID, len(s.idList[s.n]))
		for i := range s.idList[s.n] {
			s.idList[s.n][i].WriteToV2(&ids[i])
		}
		body.SetIDList(ids)
	}
	resp.SetBody(&body)

	err := signatureV2.SignServiceMessage(s.key, resp)
	if err != nil {
		panic(fmt.Errorf("error: %w", err))
	}

	s.n++
	return nil
}
