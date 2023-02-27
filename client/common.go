package client

import (
	"crypto/ecdsa"
	"errors"
	"fmt"

	"github.com/TrueCloudLab/frostfs-api-go/v2/refs"
	"github.com/TrueCloudLab/frostfs-api-go/v2/rpc/client"
	v2session "github.com/TrueCloudLab/frostfs-api-go/v2/session"
	"github.com/TrueCloudLab/frostfs-api-go/v2/signature"
	apistatus "github.com/TrueCloudLab/frostfs-sdk-go/client/status"
	"github.com/TrueCloudLab/frostfs-sdk-go/version"
)

// common interface of resulting structures with API status.
type resCommon interface {
	setStatus(apistatus.Status)
}

// structure is embedded to all resulting types in order to inherit status-related methods.
type statusRes struct {
	st apistatus.Status
}

// setStatus implements resCommon interface method.
func (x *statusRes) setStatus(st apistatus.Status) {
	x.st = st
}

// Status returns server's status return.
//
// Use apistatus package functionality to handle the status.
func (x statusRes) Status() apistatus.Status {
	return x.st
}

// groups meta parameters shared between all Client operations.
type prmCommonMeta struct {
	// FrostFS request X-Headers
	xHeaders []string
}

// WithXHeaders specifies list of extended headers (string key-value pairs)
// to be attached to the request. Must have an even length.
//
// Slice must not be mutated until the operation completes.
func (x *prmCommonMeta) WithXHeaders(hs ...string) {
	if len(hs)%2 != 0 {
		panic("slice of X-Headers with odd length")
	}

	x.xHeaders = hs
}

func writeXHeadersToMeta(xHeaders []string, h *v2session.RequestMetaHeader) {
	if len(xHeaders) == 0 {
		return
	}

	if len(xHeaders)%2 != 0 {
		panic("slice of X-Headers with odd length")
	}

	hs := make([]v2session.XHeader, len(xHeaders)/2)
	for i := 0; i < len(xHeaders); i += 2 {
		hs[i].SetKey(xHeaders[i])
		hs[i].SetValue(xHeaders[i+1])
	}

	h.SetXHeaders(hs)
}

// error messages.
var (
	errorMissingContext       = errors.New("missing context")
	errorMissingContainer     = errors.New("missing container")
	errorMissingObject        = errors.New("missing object")
	errorAccountNotSet        = errors.New("account not set")
	errorServerAddrUnset      = errors.New("server address is unset or empty")
	errorNonPositiveTimeout   = errors.New("non-positive timeout")
	errorEACLTableNotSet      = errors.New("eACL table not set")
	errorMissingAnnouncements = errors.New("missing announcements")
	errorZeroRangeLength      = errors.New("zero range length")
	errorMissingRanges        = errors.New("missing ranges")
	errorZeroEpoch            = errors.New("zero epoch")
	errorMissingTrusts        = errors.New("missing trusts")
	errorTrustNotSet          = errors.New("current trust value not set")
)

// groups all the details required to send a single request and process a response to it.
type contextCall struct {
	// ==================================================
	// state vars that do not require explicit initialization

	// final error to be returned from client method
	err error

	// received response
	resp responseV2

	// ==================================================
	// shared parameters which are set uniformly on all calls

	// request signing key
	key ecdsa.PrivateKey

	// callback prior to processing the response by the client
	callbackResp func(ResponseMetaInfo) error

	// if set, protocol errors will be expanded into a final error
	resolveAPIFailures bool

	// FrostFS network magic
	netMagic uint64

	// Meta parameters
	meta prmCommonMeta

	// ==================================================
	// custom call parameters

	// structure of the call result
	statusRes resCommon

	// request to be signed with a key and sent
	req request

	// function to send a request (unary) and receive a response
	call func() (responseV2, error)

	// function to send the request (req field)
	wReq func() error

	// function to recv the response (resp field)
	rResp func() error

	// function to close the message stream
	closer func() error

	// function of writing response fields to the resulting structure (optional)
	result func(v2 responseV2)
}

type request interface {
	GetMetaHeader() *v2session.RequestMetaHeader
	SetMetaHeader(*v2session.RequestMetaHeader)
	SetVerificationHeader(*v2session.RequestVerificationHeader)
}

// sets needed fields of the request meta header.
func (x contextCall) prepareRequest() {
	meta := x.req.GetMetaHeader()
	if meta == nil {
		meta = new(v2session.RequestMetaHeader)
		x.req.SetMetaHeader(meta)
	}

	if meta.GetTTL() == 0 {
		meta.SetTTL(2)
	}

	if meta.GetVersion() == nil {
		var verV2 refs.Version
		version.Current().WriteToV2(&verV2)
		meta.SetVersion(&verV2)
	}

	meta.SetNetworkMagic(x.netMagic)

	writeXHeadersToMeta(x.meta.xHeaders, meta)
}

func (c *Client) prepareRequest(req request, meta *v2session.RequestMetaHeader) {
	ttl := meta.GetTTL()
	if ttl == 0 {
		ttl = 2
	}

	verV2 := meta.GetVersion()
	if verV2 == nil {
		verV2 = new(refs.Version)
		version.Current().WriteToV2(verV2)
	}

	meta.SetTTL(ttl)
	meta.SetVersion(verV2)
	meta.SetNetworkMagic(c.prm.netMagic)

	req.SetMetaHeader(meta)
}

// prepares, signs and writes the request. Result means success.
// If failed, contextCall.err contains the reason.
func (x *contextCall) writeRequest() bool {
	x.prepareRequest()

	x.req.SetVerificationHeader(nil)

	// sign the request
	x.err = signature.SignServiceMessage(&x.key, x.req)
	if x.err != nil {
		x.err = fmt.Errorf("sign request: %w", x.err)
		return false
	}

	x.err = x.wReq()
	if x.err != nil {
		x.err = fmt.Errorf("write request: %w", x.err)
		return false
	}

	return true
}

// performs common actions of response processing and writes any problem as a result status or client error
// (in both cases returns false).
//
// Actions:
//   - verify signature (internal);
//   - call response callback (internal);
//   - unwrap status error (optional).
func (x *contextCall) processResponse() bool {
	// call response callback if set
	if x.callbackResp != nil {
		x.err = x.callbackResp(ResponseMetaInfo{
			key:   x.resp.GetVerificationHeader().GetBodySignature().GetKey(),
			epoch: x.resp.GetMetaHeader().GetEpoch(),
		})
		if x.err != nil {
			x.err = fmt.Errorf("response callback error: %w", x.err)
			return false
		}
	}

	// note that we call response callback before signature check since it is expected more lightweight
	// while verification needs marshaling

	// verify response signature
	x.err = signature.VerifyServiceMessage(x.resp)
	if x.err != nil {
		x.err = fmt.Errorf("invalid response signature: %w", x.err)
		return false
	}

	// get result status
	st := apistatus.FromStatusV2(x.resp.GetMetaHeader().GetStatus())

	// unwrap unsuccessful status and return it
	// as error if client has been configured so
	successfulStatus := apistatus.IsSuccessful(st)

	if x.resolveAPIFailures {
		x.err = apistatus.ErrFromStatus(st)
	} else {
		x.statusRes.setStatus(st)
	}

	return successfulStatus
}

// processResponse verifies response signature and converts status to an error if needed.
func (c *Client) processResponse(resp responseV2) (apistatus.Status, error) {
	err := signature.VerifyServiceMessage(resp)
	if err != nil {
		return nil, fmt.Errorf("invalid response signature: %w", err)
	}

	st := apistatus.FromStatusV2(resp.GetMetaHeader().GetStatus())
	if c.prm.resolveFrostFSErrors {
		return st, apistatus.ErrFromStatus(st)
	}
	return st, nil
}

// reads response (if rResp is set) and processes it. Result means success.
// If failed, contextCall.err (or statusRes if resolveAPIFailures is set) contains the reason.
func (x *contextCall) readResponse() bool {
	if x.rResp != nil {
		x.err = x.rResp()
		if x.err != nil {
			x.err = fmt.Errorf("read response: %w", x.err)
			return false
		}
	}

	return x.processResponse()
}

// closes the message stream (if closer is set) and writes the results (if result is set).
// Return means success. If failed, contextCall.err contains the reason.
func (x *contextCall) close() bool {
	if x.closer != nil {
		x.err = x.closer()
		if x.err != nil {
			x.err = fmt.Errorf("close RPC: %w", x.err)
			return false
		}
	}

	// write response to resulting structure
	if x.result != nil {
		x.result(x.resp)
	}

	return x.err == nil
}

// goes through all stages of sending a request and processing a response. Returns true if successful.
// If failed, contextCall.err contains the reason.
func (x *contextCall) processCall() bool {
	// set request writer
	x.wReq = func() error {
		var err error
		x.resp, err = x.call()
		return err
	}

	// write request
	ok := x.writeRequest()
	if !ok {
		return false
	}

	// read response
	ok = x.readResponse()
	if !ok {
		return x.err == nil
	}

	// close and write response to resulting structure
	ok = x.close()
	if !ok {
		return false
	}

	return x.err == nil
}

// initializes static cross-call parameters inherited from client.
func (c *Client) initCallContext(ctx *contextCall) {
	ctx.key = c.prm.key
	ctx.resolveAPIFailures = c.prm.resolveFrostFSErrors
	ctx.callbackResp = c.prm.cbRespInfo
	ctx.netMagic = c.prm.netMagic
}

// ExecRaw executes f with underlying github.com/TrueCloudLab/frostfs-api-go/v2/rpc/client.Client
// instance. Communicate over the Protocol Buffers protocol in a more flexible way:
// most often used to transmit data over a fixed version of the FrostFS protocol, as well
// as to support custom services.
//
// The f must not manipulate the client connection passed into it.
//
// Like all other operations, must be called after connecting to the server and
// before closing the connection.
//
// See also Dial and Close.
// See also github.com/TrueCloudLab/frostfs-api-go/v2/rpc/client package docs.
func (c *Client) ExecRaw(f func(client *client.Client) error) error {
	return f(&c.c)
}
