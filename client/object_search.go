package client

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"io"

	"github.com/TrueCloudLab/frostfs-api-go/v2/acl"
	v2object "github.com/TrueCloudLab/frostfs-api-go/v2/object"
	v2refs "github.com/TrueCloudLab/frostfs-api-go/v2/refs"
	rpcapi "github.com/TrueCloudLab/frostfs-api-go/v2/rpc"
	"github.com/TrueCloudLab/frostfs-api-go/v2/rpc/client"
	v2session "github.com/TrueCloudLab/frostfs-api-go/v2/session"
	"github.com/TrueCloudLab/frostfs-api-go/v2/signature"
	"github.com/TrueCloudLab/frostfs-sdk-go/bearer"
	apistatus "github.com/TrueCloudLab/frostfs-sdk-go/client/status"
	cid "github.com/TrueCloudLab/frostfs-sdk-go/container/id"
	"github.com/TrueCloudLab/frostfs-sdk-go/object"
	oid "github.com/TrueCloudLab/frostfs-sdk-go/object/id"
	"github.com/TrueCloudLab/frostfs-sdk-go/session"
)

// PrmObjectSearch groups parameters of ObjectSearch operation.
type PrmObjectSearch struct {
	meta v2session.RequestMetaHeader

	key *ecdsa.PrivateKey

	cnrSet bool
	cnrID  cid.ID

	filters object.SearchFilters
}

// MarkLocal tells the server to execute the operation locally.
func (x *PrmObjectSearch) MarkLocal() {
	x.meta.SetTTL(1)
}

// WithinSession specifies session within which the search query must be executed.
//
// Creator of the session acquires the authorship of the request.
// This may affect the execution of an operation (e.g. access control).
//
// Must be signed.
func (x *PrmObjectSearch) WithinSession(t session.Object) {
	var tokv2 v2session.Token
	t.WriteToV2(&tokv2)
	x.meta.SetSessionToken(&tokv2)
}

// WithBearerToken attaches bearer token to be used for the operation.
//
// If set, underlying eACL rules will be used in access control.
//
// Must be signed.
func (x *PrmObjectSearch) WithBearerToken(t bearer.Token) {
	var v2token acl.BearerToken
	t.WriteToV2(&v2token)
	x.meta.SetBearerToken(&v2token)
}

// WithXHeaders specifies list of extended headers (string key-value pairs)
// to be attached to the request. Must have an even length.
//
// Slice must not be mutated until the operation completes.
func (x *PrmObjectSearch) WithXHeaders(hs ...string) {
	writeXHeadersToMeta(hs, &x.meta)
}

// UseKey specifies private key to sign the requests.
// If key is not provided, then Client default key is used.
func (x *PrmObjectSearch) UseKey(key ecdsa.PrivateKey) {
	x.key = &key
}

// InContainer specifies the container in which to look for objects.
// Required parameter.
func (x *PrmObjectSearch) InContainer(id cid.ID) {
	x.cnrID = id
	x.cnrSet = true
}

// SetFilters sets filters by which to select objects. All container objects
// match unset/empty filters.
func (x *PrmObjectSearch) SetFilters(filters object.SearchFilters) {
	x.filters = filters
}

// ResObjectSearch groups the final result values of ObjectSearch operation.
type ResObjectSearch struct {
	statusRes
}

// ObjectListReader is designed to read list of object identifiers from FrostFS system.
//
// Must be initialized using Client.ObjectSearch, any other usage is unsafe.
type ObjectListReader struct {
	client          *Client
	cancelCtxStream context.CancelFunc
	err             error
	res             ResObjectSearch
	stream          interface {
		Read(resp *v2object.SearchResponse) error
	}
	tail []v2refs.ObjectID
}

// Read reads another list of the object identifiers. Works similar to
// io.Reader.Read but copies oid.ID and returns success flag instead of error.
//
// Failure reason can be received via Close.
//
// Panics if buf has zero length.
func (x *ObjectListReader) Read(buf []oid.ID) (int, bool) {
	if len(buf) == 0 {
		panic("empty buffer in ObjectListReader.ReadList")
	}

	read := copyIDBuffers(buf, x.tail)
	x.tail = x.tail[read:]

	if len(buf) == read {
		return read, true
	}

	for {
		var resp v2object.SearchResponse
		x.err = x.stream.Read(&resp)
		if x.err != nil {
			return read, false
		}

		x.res.st, x.err = x.client.processResponse(&resp)
		if x.err != nil || !apistatus.IsSuccessful(x.res.st) {
			return read, false
		}

		// read new chunk of objects
		ids := resp.GetBody().GetIDList()
		if len(ids) == 0 {
			// just skip empty lists since they are not prohibited by protocol
			continue
		}

		ln := copyIDBuffers(buf[read:], ids)
		read += ln

		if read == len(buf) {
			// save the tail
			x.tail = append(x.tail, ids[ln:]...)

			return read, true
		}
	}
}

func copyIDBuffers(dst []oid.ID, src []v2refs.ObjectID) int {
	var i int
	for ; i < len(dst) && i < len(src); i++ {
		_ = dst[i].ReadFromV2(src[i])
	}
	return i
}

// Iterate iterates over the list of found object identifiers.
// f can return true to stop iteration earlier.
//
// Returns an error if object can't be read.
func (x *ObjectListReader) Iterate(f func(oid.ID) bool) error {
	buf := make([]oid.ID, 1)

	for {
		// Do not check first return value because `len(buf) == 1`,
		// so false means nothing was read.
		_, ok := x.Read(buf)
		if !ok {
			res, err := x.Close()
			if err != nil {
				return err
			}
			return apistatus.ErrFromStatus(res.Status())
		}
		if f(buf[0]) {
			return nil
		}
	}
}

// Close ends reading list of the matched objects and returns the result of the operation
// along with the final results. Must be called after using the ObjectListReader.
//
// Exactly one return value is non-nil. By default, server status is returned in res structure.
// Any client's internal or transport errors are returned as Go built-in error.
// If Client is tuned to resolve FrostFS API statuses, then FrostFS failures
// codes are returned as error.
//
// Return statuses:
//   - global (see Client docs);
//   - *apistatus.ContainerNotFound;
//   - *apistatus.ObjectAccessDenied;
//   - *apistatus.SessionTokenExpired.
func (x *ObjectListReader) Close() (*ResObjectSearch, error) {
	defer x.cancelCtxStream()

	if x.err != nil && !errors.Is(x.err, io.EOF) {
		return nil, x.err
	}

	return &x.res, nil
}

// ObjectSearchInit initiates object selection through a remote server using FrostFS API protocol.
//
// The call only opens the transmission channel, explicit fetching of matched objects
// is done using the ObjectListReader. Exactly one return value is non-nil.
// Resulting reader must be finally closed.
//
// Returns an error if parameters are set incorrectly (see PrmObjectSearch docs).
// Context is required and must not be nil. It is used for network communication.
func (c *Client) ObjectSearchInit(ctx context.Context, prm PrmObjectSearch) (*ObjectListReader, error) {
	// check parameters
	switch {
	case ctx == nil:
		return nil, errorMissingContext
	case !prm.cnrSet:
		return nil, errorMissingContainer
	}

	var cidV2 v2refs.ContainerID
	prm.cnrID.WriteToV2(&cidV2)

	var body v2object.SearchRequestBody
	body.SetVersion(1)
	body.SetContainerID(&cidV2)
	body.SetFilters(prm.filters.ToV2())

	// init reader
	var req v2object.SearchRequest
	req.SetBody(&body)
	c.prepareRequest(&req, &prm.meta)

	key := prm.key
	if key == nil {
		key = &c.prm.key
	}

	err := signature.SignServiceMessage(key, &req)
	if err != nil {
		return nil, fmt.Errorf("sign request: %w", err)
	}

	var r ObjectListReader
	ctx, r.cancelCtxStream = context.WithCancel(ctx)

	r.stream, err = rpcapi.SearchObjects(&c.c, &req, client.WithContext(ctx))
	if err != nil {
		return nil, fmt.Errorf("open stream: %w", err)
	}
	r.client = c

	return &r, nil
}
