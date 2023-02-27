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

// shared parameters of GET/HEAD/RANGE.
type prmObjectRead struct {
	meta v2session.RequestMetaHeader

	raw bool

	addr v2refs.Address
}

// WithXHeaders specifies list of extended headers (string key-value pairs)
// to be attached to the request. Must have an even length.
//
// Slice must not be mutated until the operation completes.
func (x *prmObjectRead) WithXHeaders(hs ...string) {
	writeXHeadersToMeta(hs, &x.meta)
}

// MarkRaw marks an intent to read physically stored object.
func (x *prmObjectRead) MarkRaw() {
	x.raw = true
}

// MarkLocal tells the server to execute the operation locally.
func (x *prmObjectRead) MarkLocal() {
	x.meta.SetTTL(1)
}

// WithinSession specifies session within which object should be read.
//
// Creator of the session acquires the authorship of the request.
// This may affect the execution of an operation (e.g. access control).
//
// Must be signed.
func (x *prmObjectRead) WithinSession(t session.Object) {
	var tokv2 v2session.Token
	t.WriteToV2(&tokv2)
	x.meta.SetSessionToken(&tokv2)
}

// WithBearerToken attaches bearer token to be used for the operation.
//
// If set, underlying eACL rules will be used in access control.
//
// Must be signed.
func (x *prmObjectRead) WithBearerToken(t bearer.Token) {
	var v2token acl.BearerToken
	t.WriteToV2(&v2token)
	x.meta.SetBearerToken(&v2token)
}

// FromContainer specifies FrostFS container of the object.
// Required parameter.
func (x *prmObjectRead) FromContainer(id cid.ID) {
	var cnrV2 v2refs.ContainerID
	id.WriteToV2(&cnrV2)
	x.addr.SetContainerID(&cnrV2)
}

// ByID specifies identifier of the requested object.
// Required parameter.
func (x *prmObjectRead) ByID(id oid.ID) {
	var objV2 v2refs.ObjectID
	id.WriteToV2(&objV2)
	x.addr.SetObjectID(&objV2)
}

// PrmObjectGet groups parameters of ObjectGetInit operation.
type PrmObjectGet struct {
	prmObjectRead

	key *ecdsa.PrivateKey
}

// ResObjectGet groups the final result values of ObjectGetInit operation.
type ResObjectGet struct {
	statusRes
}

// ObjectReader is designed to read one object from FrostFS system.
//
// Must be initialized using Client.ObjectGetInit, any other
// usage is unsafe.
type ObjectReader struct {
	cancelCtxStream context.CancelFunc

	client *Client
	stream interface {
		Read(resp *v2object.GetResponse) error
	}

	res ResObjectGet
	err error

	tailPayload []byte

	remainingPayloadLen int
}

// UseKey specifies private key to sign the requests.
// If key is not provided, then Client default key is used.
func (x *PrmObjectGet) UseKey(key ecdsa.PrivateKey) {
	x.key = &key
}

// ReadHeader reads header of the object. Result means success.
// Failure reason can be received via Close.
func (x *ObjectReader) ReadHeader(dst *object.Object) bool {
	var resp v2object.GetResponse
	x.err = x.stream.Read(&resp)
	if x.err != nil {
		return false
	}

	x.res.st, x.err = x.client.processResponse(&resp)
	if x.err != nil || !apistatus.IsSuccessful(x.res.st) {
		return false
	}

	var partInit *v2object.GetObjectPartInit

	switch v := resp.GetBody().GetObjectPart().(type) {
	default:
		x.err = fmt.Errorf("unexpected message instead of heading part: %T", v)
		return false
	case *v2object.SplitInfo:
		x.err = object.NewSplitInfoError(object.NewSplitInfoFromV2(v))
		return false
	case *v2object.GetObjectPartInit:
		partInit = v
	}

	var objv2 v2object.Object

	objv2.SetObjectID(partInit.GetObjectID())
	objv2.SetHeader(partInit.GetHeader())
	objv2.SetSignature(partInit.GetSignature())

	x.remainingPayloadLen = int(objv2.GetHeader().GetPayloadLength())

	*dst = *object.NewFromV2(&objv2) // need smth better

	return true
}

func (x *ObjectReader) readChunk(buf []byte) (int, bool) {
	var read int

	// read remaining tail
	read = copy(buf, x.tailPayload)

	x.tailPayload = x.tailPayload[read:]

	if len(buf) == read {
		return read, true
	}

	var chunk []byte
	var lastRead int

	for {
		var resp v2object.GetResponse
		x.err = x.stream.Read(&resp)
		if x.err != nil {
			return read, false
		}

		x.res.st, x.err = x.client.processResponse(&resp)
		if x.err != nil || !apistatus.IsSuccessful(x.res.st) {
			return read, false
		}

		part := resp.GetBody().GetObjectPart()
		partChunk, ok := part.(*v2object.GetObjectPartChunk)
		if !ok {
			x.err = fmt.Errorf("unexpected message instead of chunk part: %T", part)
			return read, false
		}

		// read new chunk
		chunk = partChunk.GetChunk()
		if len(chunk) == 0 {
			// just skip empty chunks since they are not prohibited by protocol
			continue
		}

		lastRead = copy(buf[read:], chunk)

		read += lastRead

		if read == len(buf) {
			// save the tail
			x.tailPayload = append(x.tailPayload, chunk[lastRead:]...)

			return read, true
		}
	}
}

// ReadChunk reads another chunk of the object payload. Works similar to
// io.Reader.Read but returns success flag instead of error.
//
// Failure reason can be received via Close.
func (x *ObjectReader) ReadChunk(buf []byte) (int, bool) {
	return x.readChunk(buf)
}

func (x *ObjectReader) close(ignoreEOF bool) (*ResObjectGet, error) {
	defer x.cancelCtxStream()

	if x.err != nil {
		if !errors.Is(x.err, io.EOF) {
			return nil, x.err
		} else if !ignoreEOF {
			if x.remainingPayloadLen > 0 {
				return nil, io.ErrUnexpectedEOF
			}

			return nil, io.EOF
		}
	}

	return &x.res, nil
}

// Close ends reading the object and returns the result of the operation
// along with the final results. Must be called after using the ObjectReader.
//
// Exactly one return value is non-nil. By default, server status is returned in res structure.
// Any client's internal or transport errors are returned as Go built-in error.
// If Client is tuned to resolve FrostFS API statuses, then FrostFS failures
// codes are returned as error.
//
// Return errors:
//
//	*object.SplitInfoError (returned on virtual objects with PrmObjectGet.MakeRaw).
//
// Return statuses:
//   - global (see Client docs);
//   - *apistatus.ContainerNotFound;
//   - *apistatus.ObjectNotFound;
//   - *apistatus.ObjectAccessDenied;
//   - *apistatus.ObjectAlreadyRemoved;
//   - *apistatus.SessionTokenExpired.
func (x *ObjectReader) Close() (*ResObjectGet, error) {
	return x.close(true)
}

// Read implements io.Reader of the object payload.
func (x *ObjectReader) Read(p []byte) (int, error) {
	n, ok := x.readChunk(p)

	x.remainingPayloadLen -= n

	if !ok {
		res, err := x.close(false)
		if err != nil {
			return n, err
		}

		return n, apistatus.ErrFromStatus(res.Status())
	}

	if x.remainingPayloadLen < 0 {
		return n, errors.New("payload size overflow")
	}

	return n, nil
}

// ObjectGetInit initiates reading an object through a remote server using FrostFS API protocol.
//
// The call only opens the transmission channel, explicit fetching is done using the ObjectReader.
// Exactly one return value is non-nil. Resulting reader must be finally closed.
//
// Returns an error if parameters are set incorrectly (see PrmObjectGet docs).
// Context is required and must not be nil. It is used for network communication.
func (c *Client) ObjectGetInit(ctx context.Context, prm PrmObjectGet) (*ObjectReader, error) {
	// check parameters
	switch {
	case ctx == nil:
		return nil, errorMissingContext
	case prm.addr.GetContainerID() == nil:
		return nil, errorMissingContainer
	case prm.addr.GetObjectID() == nil:
		return nil, errorMissingObject
	}

	// form request body
	var body v2object.GetRequestBody

	body.SetRaw(prm.raw)
	body.SetAddress(&prm.addr)

	// form request
	var req v2object.GetRequest

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

	ctx, cancel := context.WithCancel(ctx)

	stream, err := rpcapi.GetObject(&c.c, &req, client.WithContext(ctx))
	if err != nil {
		cancel()
		return nil, fmt.Errorf("open stream: %w", err)
	}

	var r ObjectReader
	r.cancelCtxStream = cancel
	r.stream = stream
	r.client = c

	return &r, nil
}

// PrmObjectHead groups parameters of ObjectHead operation.
type PrmObjectHead struct {
	prmObjectRead

	keySet bool
	key    ecdsa.PrivateKey
}

// UseKey specifies private key to sign the requests.
// If key is not provided, then Client default key is used.
func (x *PrmObjectHead) UseKey(key ecdsa.PrivateKey) {
	x.keySet = true
	x.key = key
}

// ResObjectHead groups resulting values of ObjectHead operation.
type ResObjectHead struct {
	statusRes

	// requested object (response doesn't carry the ID)
	idObj oid.ID

	hdr *v2object.HeaderWithSignature
}

// ReadHeader reads header of the requested object.
// Returns false if header is missing in the response (not read).
func (x *ResObjectHead) ReadHeader(dst *object.Object) bool {
	if x.hdr == nil {
		return false
	}

	var objv2 v2object.Object

	objv2.SetHeader(x.hdr.GetHeader())
	objv2.SetSignature(x.hdr.GetSignature())

	obj := object.NewFromV2(&objv2)
	obj.SetID(x.idObj)

	*dst = *obj

	return true
}

// ObjectHead reads object header through a remote server using FrostFS API protocol.
//
// Exactly one return value is non-nil. By default, server status is returned in res structure.
// Any client's internal or transport errors are returned as `error`,
// If PrmInit.ResolveFrostFSFailures has been called, unsuccessful
// FrostFS status codes are returned as `error`, otherwise, are included
// in the returned result structure.
//
// Returns an error if parameters are set incorrectly (see PrmObjectHead docs).
// Context is required and must not be nil. It is used for network communication.
//
// Return errors:
//
//	*object.SplitInfoError (returned on virtual objects with PrmObjectHead.MakeRaw).
//
// Return statuses:
//   - global (see Client docs);
//   - *apistatus.ContainerNotFound;
//   - *apistatus.ObjectNotFound;
//   - *apistatus.ObjectAccessDenied;
//   - *apistatus.ObjectAlreadyRemoved;
//   - *apistatus.SessionTokenExpired.
func (c *Client) ObjectHead(ctx context.Context, prm PrmObjectHead) (*ResObjectHead, error) {
	switch {
	case ctx == nil:
		return nil, errorMissingContext
	case prm.addr.GetContainerID() == nil:
		return nil, errorMissingContainer
	case prm.addr.GetObjectID() == nil:
		return nil, errorMissingObject
	}

	var body v2object.HeadRequestBody
	body.SetRaw(prm.raw)
	body.SetAddress(&prm.addr)

	var req v2object.HeadRequest
	req.SetBody(&body)
	c.prepareRequest(&req, &prm.meta)

	key := c.prm.key
	if prm.keySet {
		key = prm.key
	}

	// sign the request
	err := signature.SignServiceMessage(&key, &req)
	if err != nil {
		return nil, fmt.Errorf("sign request: %w", err)
	}

	resp, err := rpcapi.HeadObject(&c.c, &req, client.WithContext(ctx))
	if err != nil {
		return nil, fmt.Errorf("write request: %w", err)
	}

	var res ResObjectHead
	res.st, err = c.processResponse(resp)
	if err != nil {
		return nil, err
	}

	if !apistatus.IsSuccessful(res.st) {
		return &res, nil
	}

	_ = res.idObj.ReadFromV2(*prm.addr.GetObjectID())

	switch v := resp.GetBody().GetHeaderPart().(type) {
	default:
		return nil, fmt.Errorf("unexpected header type %T", v)
	case *v2object.SplitInfo:
		return nil, object.NewSplitInfoError(object.NewSplitInfoFromV2(v))
	case *v2object.HeaderWithSignature:
		res.hdr = v
	}

	return &res, nil
}

// PrmObjectRange groups parameters of ObjectRange operation.
type PrmObjectRange struct {
	prmObjectRead

	key *ecdsa.PrivateKey

	rng v2object.Range
}

// SetOffset sets offset of the payload range to be read.
// Zero by default.
func (x *PrmObjectRange) SetOffset(off uint64) {
	x.rng.SetOffset(off)
}

// SetLength sets length of the payload range to be read.
// Must be positive.
func (x *PrmObjectRange) SetLength(ln uint64) {
	x.rng.SetLength(ln)
}

// UseKey specifies private key to sign the requests.
// If key is not provided, then Client default key is used.
func (x *PrmObjectRange) UseKey(key ecdsa.PrivateKey) {
	x.key = &key
}

// ResObjectRange groups the final result values of ObjectRange operation.
type ResObjectRange struct {
	statusRes
}

// ObjectRangeReader is designed to read payload range of one object
// from FrostFS system.
//
// Must be initialized using Client.ObjectRangeInit, any other
// usage is unsafe.
type ObjectRangeReader struct {
	cancelCtxStream context.CancelFunc

	client *Client

	res ResObjectRange
	err error

	stream interface {
		Read(resp *v2object.GetRangeResponse) error
	}

	tailPayload []byte

	remainingPayloadLen int
}

func (x *ObjectRangeReader) readChunk(buf []byte) (int, bool) {
	var read int

	// read remaining tail
	read = copy(buf, x.tailPayload)

	x.tailPayload = x.tailPayload[read:]

	if len(buf) == read {
		return read, true
	}

	var partChunk *v2object.GetRangePartChunk
	var chunk []byte
	var lastRead int

	for {
		var resp v2object.GetRangeResponse
		x.err = x.stream.Read(&resp)
		if x.err != nil {
			return read, false
		}

		x.res.st, x.err = x.client.processResponse(&resp)
		if x.err != nil || !apistatus.IsSuccessful(x.res.st) {
			return read, false
		}

		// get chunk message
		switch v := resp.GetBody().GetRangePart().(type) {
		default:
			x.err = fmt.Errorf("unexpected message received: %T", v)
			return read, false
		case *v2object.SplitInfo:
			x.err = object.NewSplitInfoError(object.NewSplitInfoFromV2(v))
			return read, false
		case *v2object.GetRangePartChunk:
			partChunk = v
		}

		chunk = partChunk.GetChunk()
		if len(chunk) == 0 {
			// just skip empty chunks since they are not prohibited by protocol
			continue
		}

		lastRead = copy(buf[read:], chunk)

		read += lastRead

		if read == len(buf) {
			// save the tail
			x.tailPayload = append(x.tailPayload, chunk[lastRead:]...)

			return read, true
		}
	}
}

// ReadChunk reads another chunk of the object payload range.
// Works similar to io.Reader.Read but returns success flag instead of error.
//
// Failure reason can be received via Close.
func (x *ObjectRangeReader) ReadChunk(buf []byte) (int, bool) {
	return x.readChunk(buf)
}

func (x *ObjectRangeReader) close(ignoreEOF bool) (*ResObjectRange, error) {
	defer x.cancelCtxStream()

	if x.err != nil {
		if !errors.Is(x.err, io.EOF) {
			return nil, x.err
		} else if !ignoreEOF {
			if x.remainingPayloadLen > 0 {
				return nil, io.ErrUnexpectedEOF
			}

			return nil, io.EOF
		}
	}

	return &x.res, nil
}

// Close ends reading the payload range and returns the result of the operation
// along with the final results. Must be called after using the ObjectRangeReader.
//
// Exactly one return value is non-nil. By default, server status is returned in res structure.
// Any client's internal or transport errors are returned as Go built-in error.
// If Client is tuned to resolve FrostFS API statuses, then FrostFS failures
// codes are returned as error.
//
// Return errors:
//
//	*object.SplitInfoError (returned on virtual objects with PrmObjectRange.MakeRaw).
//
// Return statuses:
//   - global (see Client docs);
//   - *apistatus.ContainerNotFound;
//   - *apistatus.ObjectNotFound;
//   - *apistatus.ObjectAccessDenied;
//   - *apistatus.ObjectAlreadyRemoved;
//   - *apistatus.ObjectOutOfRange;
//   - *apistatus.SessionTokenExpired.
func (x *ObjectRangeReader) Close() (*ResObjectRange, error) {
	return x.close(true)
}

// Read implements io.Reader of the object payload.
func (x *ObjectRangeReader) Read(p []byte) (int, error) {
	n, ok := x.readChunk(p)

	x.remainingPayloadLen -= n

	if !ok {
		res, err := x.close(false)
		if err != nil {
			return n, err
		}

		return n, apistatus.ErrFromStatus(res.Status())
	}

	if x.remainingPayloadLen < 0 {
		return n, errors.New("payload range size overflow")
	}

	return n, nil
}

// ObjectRangeInit initiates reading an object's payload range through a remote
// server using FrostFS API protocol.
//
// The call only opens the transmission channel, explicit fetching is done using the ObjectRangeReader.
// Exactly one return value is non-nil. Resulting reader must be finally closed.
//
// Returns an error if parameters are set incorrectly (see PrmObjectRange docs).
// Context is required and must not be nil. It is used for network communication.
func (c *Client) ObjectRangeInit(ctx context.Context, prm PrmObjectRange) (*ObjectRangeReader, error) {
	// check parameters
	switch {
	case ctx == nil:
		return nil, errorMissingContext
	case prm.addr.GetContainerID() == nil:
		return nil, errorMissingContainer
	case prm.addr.GetObjectID() == nil:
		return nil, errorMissingObject
	case prm.rng.GetLength() == 0:
		return nil, errorZeroRangeLength
	}

	// form request body
	var body v2object.GetRangeRequestBody

	body.SetRaw(prm.raw)
	body.SetAddress(&prm.addr)
	body.SetRange(&prm.rng)

	// form request
	var req v2object.GetRangeRequest

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

	ctx, cancel := context.WithCancel(ctx)

	stream, err := rpcapi.GetObjectRange(&c.c, &req, client.WithContext(ctx))
	if err != nil {
		cancel()
		return nil, fmt.Errorf("open stream: %w", err)
	}

	var r ObjectRangeReader
	r.remainingPayloadLen = int(prm.rng.GetLength())
	r.cancelCtxStream = cancel
	r.stream = stream
	r.client = c

	return &r, nil
}
