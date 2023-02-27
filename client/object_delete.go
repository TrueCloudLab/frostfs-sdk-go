package client

import (
	"context"
	"crypto/ecdsa"
	"fmt"

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
	oid "github.com/TrueCloudLab/frostfs-sdk-go/object/id"
	"github.com/TrueCloudLab/frostfs-sdk-go/session"
)

// PrmObjectDelete groups parameters of ObjectDelete operation.
type PrmObjectDelete struct {
	meta v2session.RequestMetaHeader

	body v2object.DeleteRequestBody

	addr v2refs.Address

	keySet bool
	key    ecdsa.PrivateKey
}

// WithinSession specifies session within which object should be read.
//
// Creator of the session acquires the authorship of the request.
// This may affect the execution of an operation (e.g. access control).
//
// Must be signed.
func (x *PrmObjectDelete) WithinSession(t session.Object) {
	var tv2 v2session.Token
	t.WriteToV2(&tv2)

	x.meta.SetSessionToken(&tv2)
}

// WithBearerToken attaches bearer token to be used for the operation.
//
// If set, underlying eACL rules will be used in access control.
//
// Must be signed.
func (x *PrmObjectDelete) WithBearerToken(t bearer.Token) {
	var v2token acl.BearerToken
	t.WriteToV2(&v2token)
	x.meta.SetBearerToken(&v2token)
}

// FromContainer specifies FrostFS container of the object.
// Required parameter.
func (x *PrmObjectDelete) FromContainer(id cid.ID) {
	var cidV2 v2refs.ContainerID
	id.WriteToV2(&cidV2)

	x.addr.SetContainerID(&cidV2)
}

// ByID specifies identifier of the requested object.
// Required parameter.
func (x *PrmObjectDelete) ByID(id oid.ID) {
	var idV2 v2refs.ObjectID
	id.WriteToV2(&idV2)

	x.addr.SetObjectID(&idV2)
}

// UseKey specifies private key to sign the requests.
// If key is not provided, then Client default key is used.
func (x *PrmObjectDelete) UseKey(key ecdsa.PrivateKey) {
	x.keySet = true
	x.key = key
}

// WithXHeaders specifies list of extended headers (string key-value pairs)
// to be attached to the request. Must have an even length.
//
// Slice must not be mutated until the operation completes.
func (x *PrmObjectDelete) WithXHeaders(hs ...string) {
	writeXHeadersToMeta(hs, &x.meta)
}

// ResObjectDelete groups resulting values of ObjectDelete operation.
type ResObjectDelete struct {
	statusRes

	tomb oid.ID
}

// Tombstone returns identifier of the created tombstone object.
func (x ResObjectDelete) Tombstone() oid.ID {
	return x.tomb
}

// ObjectDelete marks an object for deletion from the container using FrostFS API protocol.
// As a marker, a special unit called a tombstone is placed in the container.
// It confirms the user's intent to delete the object, and is itself a container object.
// Explicit deletion is done asynchronously, and is generally not guaranteed.
//
// Returns a list of checksums in raw form: the format of hashes and their number
// is left for the caller to check. Client preserves the order of the server's response.
//
// Exactly one return value is non-nil. By default, server status is returned in res structure.
// Any client's internal or transport errors are returned as `error`,
// If PrmInit.ResolveFrostFSFailures has been called, unsuccessful
// FrostFS status codes are returned as `error`, otherwise, are included
// in the returned result structure.
//
// Returns an error if parameters are set incorrectly (see PrmObjectDelete docs).
// Context is required and must not be nil. It is used for network communication.
//
// Return statuses:
//   - global (see Client docs)
//   - *apistatus.ContainerNotFound;
//   - *apistatus.ObjectAccessDenied;
//   - *apistatus.ObjectLocked;
//   - *apistatus.SessionTokenExpired.
func (c *Client) ObjectDelete(ctx context.Context, prm PrmObjectDelete) (*ResObjectDelete, error) {
	switch {
	case ctx == nil:
		return nil, errorMissingContext
	case prm.addr.GetContainerID() == nil:
		return nil, errorMissingContainer
	case prm.addr.GetObjectID() == nil:
		return nil, errorMissingObject
	}

	// form request body
	prm.body.SetAddress(&prm.addr)

	// form request
	var req v2object.DeleteRequest
	req.SetBody(&prm.body)
	c.prepareRequest(&req, &prm.meta)

	key := c.prm.key
	if prm.keySet {
		key = prm.key
	}

	err := signature.SignServiceMessage(&key, &req)
	if err != nil {
		return nil, fmt.Errorf("sign request: %w", err)
	}

	resp, err := rpcapi.DeleteObject(&c.c, &req, client.WithContext(ctx))
	if err != nil {
		return nil, err
	}

	var res ResObjectDelete
	res.st, err = c.processResponse(resp)
	if err != nil {
		return nil, err
	}

	if !apistatus.IsSuccessful(res.st) {
		return &res, nil
	}

	const fieldTombstone = "tombstone"

	idTombV2 := resp.GetBody().GetTombstone().GetObjectID()
	if idTombV2 == nil {
		return nil, newErrMissingResponseField(fieldTombstone)
	}

	err = res.tomb.ReadFromV2(*idTombV2)
	if err != nil {
		return nil, newErrInvalidResponseField(fieldTombstone, err)
	}

	return &res, nil
}
