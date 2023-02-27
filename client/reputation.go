package client

import (
	"context"

	v2reputation "github.com/TrueCloudLab/frostfs-api-go/v2/reputation"
	rpcapi "github.com/TrueCloudLab/frostfs-api-go/v2/rpc"
	"github.com/TrueCloudLab/frostfs-api-go/v2/rpc/client"
	"github.com/TrueCloudLab/frostfs-sdk-go/reputation"
)

// PrmAnnounceLocalTrust groups parameters of AnnounceLocalTrust operation.
type PrmAnnounceLocalTrust struct {
	prmCommonMeta

	epoch uint64

	trusts []reputation.Trust
}

// SetEpoch sets number of FrostFS epoch in which the trust was assessed.
// Required parameter, must not be zero.
func (x *PrmAnnounceLocalTrust) SetEpoch(epoch uint64) {
	x.epoch = epoch
}

// SetValues sets values describing trust of the client to the FrostFS network participants.
// Required parameter. Must not be empty.
//
// Must not be mutated before the end of the operation.
func (x *PrmAnnounceLocalTrust) SetValues(trusts []reputation.Trust) {
	x.trusts = trusts
}

// ResAnnounceLocalTrust groups results of AnnounceLocalTrust operation.
type ResAnnounceLocalTrust struct {
	statusRes
}

// AnnounceLocalTrust sends client's trust values to the FrostFS network participants.
//
// Exactly one return value is non-nil. By default, server status is returned in res structure.
// Any client's internal or transport errors are returned as `error`.
// If PrmInit.ResolveFrostFSFailures has been called, unsuccessful
// FrostFS status codes are returned as `error`, otherwise, are included
// in the returned result structure.
//
// Returns an error if parameters are set incorrectly (see PrmAnnounceLocalTrust docs).
// Context is required and must not be nil. It is used for network communication.
//
// Return statuses:
//   - global (see Client docs).
func (c *Client) AnnounceLocalTrust(ctx context.Context, prm PrmAnnounceLocalTrust) (*ResAnnounceLocalTrust, error) {
	// check parameters
	switch {
	case ctx == nil:
		return nil, errorMissingContext
	case prm.epoch == 0:
		return nil, errorZeroEpoch
	case len(prm.trusts) == 0:
		return nil, errorMissingTrusts
	}

	// form request body
	reqBody := new(v2reputation.AnnounceLocalTrustRequestBody)
	reqBody.SetEpoch(prm.epoch)

	trusts := make([]v2reputation.Trust, len(prm.trusts))

	for i := range prm.trusts {
		prm.trusts[i].WriteToV2(&trusts[i])
	}

	reqBody.SetTrusts(trusts)

	// form request
	var req v2reputation.AnnounceLocalTrustRequest

	req.SetBody(reqBody)

	// init call context

	var (
		cc  contextCall
		res ResAnnounceLocalTrust
	)

	c.initCallContext(&cc)
	cc.meta = prm.prmCommonMeta
	cc.req = &req
	cc.statusRes = &res
	cc.call = func() (responseV2, error) {
		return rpcapi.AnnounceLocalTrust(&c.c, &req, client.WithContext(ctx))
	}

	// process call
	if !cc.processCall() {
		return nil, cc.err
	}

	return &res, nil
}

// PrmAnnounceIntermediateTrust groups parameters of AnnounceIntermediateTrust operation.
type PrmAnnounceIntermediateTrust struct {
	prmCommonMeta

	epoch uint64

	iter uint32

	trustSet bool
	trust    reputation.PeerToPeerTrust
}

// SetEpoch sets number of FrostFS epoch with which client's calculation algorithm is initialized.
// Required parameter, must not be zero.
func (x *PrmAnnounceIntermediateTrust) SetEpoch(epoch uint64) {
	x.epoch = epoch
}

// SetIteration sets current sequence number of the client's calculation algorithm.
// By default, corresponds to initial (zero) iteration.
func (x *PrmAnnounceIntermediateTrust) SetIteration(iter uint32) {
	x.iter = iter
}

// SetCurrentValue sets current global trust value computed at the specified iteration
// of the client's calculation algorithm. Required parameter.
func (x *PrmAnnounceIntermediateTrust) SetCurrentValue(trust reputation.PeerToPeerTrust) {
	x.trust = trust
	x.trustSet = true
}

// ResAnnounceIntermediateTrust groups results of AnnounceIntermediateTrust operation.
type ResAnnounceIntermediateTrust struct {
	statusRes
}

// AnnounceIntermediateTrust sends global trust values calculated for the specified FrostFS network participants
// at some stage of client's calculation algorithm.
//
// Exactly one return value is non-nil. By default, server status is returned in res structure.
// Any client's internal or transport errors are returned as `error`.
// If PrmInit.ResolveFrostFSFailures has been called, unsuccessful
// FrostFS status codes are returned as `error`, otherwise, are included
// in the returned result structure.
//
// Returns an error if parameters are set incorrectly (see PrmAnnounceIntermediateTrust docs).
// Context is required and must not be nil. It is used for network communication.
//
// Return statuses:
//   - global (see Client docs).
func (c *Client) AnnounceIntermediateTrust(ctx context.Context, prm PrmAnnounceIntermediateTrust) (*ResAnnounceIntermediateTrust, error) {
	// check parameters
	switch {
	case ctx == nil:
		return nil, errorMissingContext
	case prm.epoch == 0:
		return nil, errorZeroEpoch
	case !prm.trustSet:
		return nil, errorTrustNotSet
	}

	var trust v2reputation.PeerToPeerTrust
	prm.trust.WriteToV2(&trust)

	// form request body
	reqBody := new(v2reputation.AnnounceIntermediateResultRequestBody)
	reqBody.SetEpoch(prm.epoch)
	reqBody.SetIteration(prm.iter)
	reqBody.SetTrust(&trust)

	// form request
	var req v2reputation.AnnounceIntermediateResultRequest

	req.SetBody(reqBody)

	// init call context

	var (
		cc  contextCall
		res ResAnnounceIntermediateTrust
	)

	c.initCallContext(&cc)
	cc.meta = prm.prmCommonMeta
	cc.req = &req
	cc.statusRes = &res
	cc.call = func() (responseV2, error) {
		return rpcapi.AnnounceIntermediateResult(&c.c, &req, client.WithContext(ctx))
	}

	// process call
	if !cc.processCall() {
		return nil, cc.err
	}

	return &res, nil
}
