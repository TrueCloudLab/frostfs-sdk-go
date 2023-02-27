package client

import (
	"context"
	"crypto/ecdsa"

	"github.com/TrueCloudLab/frostfs-api-go/v2/refs"
	rpcapi "github.com/TrueCloudLab/frostfs-api-go/v2/rpc"
	"github.com/TrueCloudLab/frostfs-api-go/v2/rpc/client"
	v2session "github.com/TrueCloudLab/frostfs-api-go/v2/session"
	"github.com/TrueCloudLab/frostfs-sdk-go/user"
)

// PrmSessionCreate groups parameters of SessionCreate operation.
type PrmSessionCreate struct {
	prmCommonMeta

	exp uint64

	keySet bool
	key    ecdsa.PrivateKey
}

// SetExp sets number of the last NepFS epoch in the lifetime of the session after which it will be expired.
func (x *PrmSessionCreate) SetExp(exp uint64) {
	x.exp = exp
}

// UseKey specifies private key to sign the requests and compute token owner.
// If key is not provided, then Client default key is used.
func (x *PrmSessionCreate) UseKey(key ecdsa.PrivateKey) {
	x.keySet = true
	x.key = key
}

// ResSessionCreate groups resulting values of SessionCreate operation.
type ResSessionCreate struct {
	statusRes

	id []byte

	sessionKey []byte
}

func (x *ResSessionCreate) setID(id []byte) {
	x.id = id
}

// ID returns identifier of the opened session in a binary FrostFS API protocol format.
//
// Client doesn't retain value so modification is safe.
func (x ResSessionCreate) ID() []byte {
	return x.id
}

func (x *ResSessionCreate) setSessionKey(key []byte) {
	x.sessionKey = key
}

// PublicKey returns public key of the opened session in a binary FrostFS API protocol format.
func (x ResSessionCreate) PublicKey() []byte {
	return x.sessionKey
}

// SessionCreate opens a session with the node server on the remote endpoint.
// The session lifetime coincides with the server lifetime. Results can be written
// to session token which can be later attached to the requests.
//
// Exactly one return value is non-nil. By default, server status is returned in res structure.
// Any client's internal or transport errors are returned as `error`.
// If PrmInit.ResolveFrostFSFailures has been called, unsuccessful
// FrostFS status codes are returned as `error`, otherwise, are included
// in the returned result structure.
//
// Returns an error if parameters are set incorrectly (see PrmSessionCreate docs).
// Context is required and must not be nil. It is used for network communication.
//
// Return statuses:
//   - global (see Client docs).
func (c *Client) SessionCreate(ctx context.Context, prm PrmSessionCreate) (*ResSessionCreate, error) {
	// check context
	if ctx == nil {
		return nil, errorMissingContext
	}

	ownerKey := c.prm.key.PublicKey
	if prm.keySet {
		ownerKey = prm.key.PublicKey
	}
	var ownerID user.ID
	user.IDFromKey(&ownerID, ownerKey)

	var ownerIDV2 refs.OwnerID
	ownerID.WriteToV2(&ownerIDV2)

	// form request body
	reqBody := new(v2session.CreateRequestBody)
	reqBody.SetOwnerID(&ownerIDV2)
	reqBody.SetExpiration(prm.exp)

	// for request
	var req v2session.CreateRequest

	req.SetBody(reqBody)

	// init call context

	var (
		cc  contextCall
		res ResSessionCreate
	)

	c.initCallContext(&cc)
	if prm.keySet {
		cc.key = prm.key
	}

	cc.meta = prm.prmCommonMeta
	cc.req = &req
	cc.statusRes = &res
	cc.call = func() (responseV2, error) {
		return rpcapi.CreateSession(&c.c, &req, client.WithContext(ctx))
	}
	cc.result = func(r responseV2) {
		resp := r.(*v2session.CreateResponse)

		body := resp.GetBody()

		res.setID(body.GetID())
		res.setSessionKey(body.GetSessionKey())
	}

	// process call
	if !cc.processCall() {
		return nil, cc.err
	}

	return &res, nil
}
