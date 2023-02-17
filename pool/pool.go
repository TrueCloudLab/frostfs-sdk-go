package pool

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"io"
	"math"
	"math/rand"
	"sort"
	"sync"
	"time"

	"github.com/TrueCloudLab/frostfs-sdk-go/accounting"
	"github.com/TrueCloudLab/frostfs-sdk-go/bearer"
	sdkClient "github.com/TrueCloudLab/frostfs-sdk-go/client"
	apistatus "github.com/TrueCloudLab/frostfs-sdk-go/client/status"
	"github.com/TrueCloudLab/frostfs-sdk-go/container"
	cid "github.com/TrueCloudLab/frostfs-sdk-go/container/id"
	frostfsecdsa "github.com/TrueCloudLab/frostfs-sdk-go/crypto/ecdsa"
	"github.com/TrueCloudLab/frostfs-sdk-go/eacl"
	"github.com/TrueCloudLab/frostfs-sdk-go/netmap"
	"github.com/TrueCloudLab/frostfs-sdk-go/object"
	oid "github.com/TrueCloudLab/frostfs-sdk-go/object/id"
	"github.com/TrueCloudLab/frostfs-sdk-go/object/relations"
	"github.com/TrueCloudLab/frostfs-sdk-go/session"
	"github.com/TrueCloudLab/frostfs-sdk-go/user"
	"github.com/google/uuid"
	"github.com/nspcc-dev/neo-go/pkg/crypto/keys"
	"go.uber.org/atomic"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// client represents virtual connection to the single FrostFS network endpoint from which Pool is formed.
// This interface is expected to have exactly one production implementation - clientWrapper.
// Others are expected to be for test purposes only.
type client interface {
	// see clientWrapper.balanceGet.
	balanceGet(context.Context, PrmBalanceGet) (accounting.Decimal, error)
	// see clientWrapper.containerPut.
	containerPut(context.Context, PrmContainerPut) (cid.ID, error)
	// see clientWrapper.containerGet.
	containerGet(context.Context, PrmContainerGet) (container.Container, error)
	// see clientWrapper.containerList.
	containerList(context.Context, PrmContainerList) ([]cid.ID, error)
	// see clientWrapper.containerDelete.
	containerDelete(context.Context, PrmContainerDelete) error
	// see clientWrapper.containerEACL.
	containerEACL(context.Context, PrmContainerEACL) (eacl.Table, error)
	// see clientWrapper.containerSetEACL.
	containerSetEACL(context.Context, PrmContainerSetEACL) error
	// see clientWrapper.endpointInfo.
	endpointInfo(context.Context, prmEndpointInfo) (netmap.NodeInfo, error)
	// see clientWrapper.networkInfo.
	networkInfo(context.Context, prmNetworkInfo) (netmap.NetworkInfo, error)
	// see clientWrapper.objectPut.
	objectPut(context.Context, PrmObjectPut) (oid.ID, error)
	// see clientWrapper.objectDelete.
	objectDelete(context.Context, PrmObjectDelete) error
	// see clientWrapper.objectGet.
	objectGet(context.Context, PrmObjectGet) (ResGetObject, error)
	// see clientWrapper.objectHead.
	objectHead(context.Context, PrmObjectHead) (object.Object, error)
	// see clientWrapper.objectRange.
	objectRange(context.Context, PrmObjectRange) (ResObjectRange, error)
	// see clientWrapper.objectSearch.
	objectSearch(context.Context, PrmObjectSearch) (ResObjectSearch, error)
	// see clientWrapper.sessionCreate.
	sessionCreate(context.Context, prmCreateSession) (resCreateSession, error)

	clientStatus

	// see clientWrapper.dial.
	dial(ctx context.Context) error
	// see clientWrapper.restartIfUnhealthy.
	restartIfUnhealthy(ctx context.Context) (bool, bool)
}

// clientStatus provide access to some metrics for connection.
type clientStatus interface {
	// isHealthy checks if the connection can handle requests.
	isHealthy() bool
	// setUnhealthy marks client as unhealthy.
	setUnhealthy()
	// address return address of endpoint.
	address() string
	// currentErrorRate returns current errors rate.
	// After specific threshold connection is considered as unhealthy.
	// Pool.startRebalance routine can make this connection healthy again.
	currentErrorRate() uint32
	// overallErrorRate returns the number of all happened errors.
	overallErrorRate() uint64
	// methodsStatus returns statistic for all used methods.
	methodsStatus() []statusSnapshot
}

// errPoolClientUnhealthy is an error to indicate that client in pool is unhealthy.
var errPoolClientUnhealthy = errors.New("pool client unhealthy")

// clientStatusMonitor count error rate and other statistics for connection.
type clientStatusMonitor struct {
	logger         *zap.Logger
	addr           string
	healthy        *atomic.Bool
	errorThreshold uint32

	mu                sync.RWMutex // protect counters
	currentErrorCount uint32
	overallErrorCount uint64
	methods           []*methodStatus
}

// methodStatus provide statistic for specific method.
type methodStatus struct {
	name string
	mu   sync.RWMutex // protect counters
	statusSnapshot
}

// statusSnapshot is statistic for specific method.
type statusSnapshot struct {
	allTime     uint64
	allRequests uint64
}

// MethodIndex index of method in list of statuses in clientStatusMonitor.
type MethodIndex int

const (
	methodBalanceGet MethodIndex = iota
	methodContainerPut
	methodContainerGet
	methodContainerList
	methodContainerDelete
	methodContainerEACL
	methodContainerSetEACL
	methodEndpointInfo
	methodNetworkInfo
	methodObjectPut
	methodObjectDelete
	methodObjectGet
	methodObjectHead
	methodObjectRange
	methodSessionCreate
	methodLast
)

// String implements fmt.Stringer.
func (m MethodIndex) String() string {
	switch m {
	case methodBalanceGet:
		return "balanceGet"
	case methodContainerPut:
		return "containerPut"
	case methodContainerGet:
		return "containerGet"
	case methodContainerList:
		return "containerList"
	case methodContainerDelete:
		return "containerDelete"
	case methodContainerEACL:
		return "containerEACL"
	case methodContainerSetEACL:
		return "containerSetEACL"
	case methodEndpointInfo:
		return "endpointInfo"
	case methodNetworkInfo:
		return "networkInfo"
	case methodObjectPut:
		return "objectPut"
	case methodObjectDelete:
		return "objectDelete"
	case methodObjectGet:
		return "objectGet"
	case methodObjectHead:
		return "objectHead"
	case methodObjectRange:
		return "objectRange"
	case methodSessionCreate:
		return "sessionCreate"
	case methodLast:
		return "it's a system name rather than a method"
	default:
		return "unknown"
	}
}

func newClientStatusMonitor(logger *zap.Logger, addr string, errorThreshold uint32) clientStatusMonitor {
	methods := make([]*methodStatus, methodLast)
	for i := methodBalanceGet; i < methodLast; i++ {
		methods[i] = &methodStatus{name: i.String()}
	}

	return clientStatusMonitor{
		logger:         logger,
		addr:           addr,
		healthy:        atomic.NewBool(true),
		errorThreshold: errorThreshold,
		methods:        methods,
	}
}

func (m *methodStatus) snapshot() statusSnapshot {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.statusSnapshot
}

func (m *methodStatus) incRequests(elapsed time.Duration) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.allTime += uint64(elapsed)
	m.allRequests++
}

// clientWrapper is used by default, alternative implementations are intended for testing purposes only.
type clientWrapper struct {
	clientMutex sync.RWMutex
	client      *sdkClient.Client
	prm         wrapperPrm

	clientStatusMonitor
}

// wrapperPrm is params to create clientWrapper.
type wrapperPrm struct {
	logger                  *zap.Logger
	address                 string
	key                     ecdsa.PrivateKey
	dialTimeout             time.Duration
	streamTimeout           time.Duration
	errorThreshold          uint32
	responseInfoCallback    func(sdkClient.ResponseMetaInfo) error
	poolRequestInfoCallback func(RequestInfo)
}

// setAddress sets endpoint to connect in FrostFS network.
func (x *wrapperPrm) setAddress(address string) {
	x.address = address
}

// setKey sets sdkClient.Client private key to be used for the protocol communication by default.
func (x *wrapperPrm) setKey(key ecdsa.PrivateKey) {
	x.key = key
}

// setDialTimeout sets the timeout for connection to be established.
func (x *wrapperPrm) setDialTimeout(timeout time.Duration) {
	x.dialTimeout = timeout
}

// setStreamTimeout sets the timeout for individual operations in streaming RPC.
func (x *wrapperPrm) setStreamTimeout(timeout time.Duration) {
	x.streamTimeout = timeout
}

// setErrorThreshold sets threshold after reaching which connection is considered unhealthy
// until Pool.startRebalance routing updates its status.
func (x *wrapperPrm) setErrorThreshold(threshold uint32) {
	x.errorThreshold = threshold
}

// setPoolRequestCallback sets callback that will be invoked after every pool response.
func (x *wrapperPrm) setPoolRequestCallback(f func(RequestInfo)) {
	x.poolRequestInfoCallback = f
}

// setResponseInfoCallback sets callback that will be invoked after every response.
func (x *wrapperPrm) setResponseInfoCallback(f func(sdkClient.ResponseMetaInfo) error) {
	x.responseInfoCallback = f
}

// newWrapper creates a clientWrapper that implements the client interface.
func newWrapper(prm wrapperPrm) *clientWrapper {
	var cl sdkClient.Client
	var prmInit sdkClient.PrmInit
	prmInit.SetDefaultPrivateKey(prm.key)
	prmInit.SetResponseInfoCallback(prm.responseInfoCallback)

	cl.Init(prmInit)

	res := &clientWrapper{
		client:              &cl,
		clientStatusMonitor: newClientStatusMonitor(prm.logger, prm.address, prm.errorThreshold),
		prm:                 prm,
	}

	return res
}

// dial establishes a connection to the server from the FrostFS network.
// Returns an error describing failure reason. If failed, the client
// SHOULD NOT be used.
func (c *clientWrapper) dial(ctx context.Context) error {
	cl, err := c.getClient()
	if err != nil {
		return err
	}

	var prmDial sdkClient.PrmDial
	prmDial.SetServerURI(c.prm.address)
	prmDial.SetTimeout(c.prm.dialTimeout)
	prmDial.SetStreamTimeout(c.prm.streamTimeout)
	prmDial.SetContext(ctx)

	if err = cl.Dial(prmDial); err != nil {
		c.setUnhealthy()
		return err
	}

	return nil
}

// restartIfUnhealthy checks healthy status of client and recreate it if status is unhealthy.
// Return current healthy status and indicating if status was changed by this function call.
func (c *clientWrapper) restartIfUnhealthy(ctx context.Context) (healthy, changed bool) {
	var wasHealthy bool
	if _, err := c.endpointInfo(ctx, prmEndpointInfo{}); err == nil {
		return true, false
	} else if !errors.Is(err, errPoolClientUnhealthy) {
		wasHealthy = true
	}

	var cl sdkClient.Client
	var prmInit sdkClient.PrmInit
	prmInit.SetDefaultPrivateKey(c.prm.key)
	prmInit.SetResponseInfoCallback(c.prm.responseInfoCallback)

	cl.Init(prmInit)

	var prmDial sdkClient.PrmDial
	prmDial.SetServerURI(c.prm.address)
	prmDial.SetTimeout(c.prm.dialTimeout)
	prmDial.SetStreamTimeout(c.prm.streamTimeout)
	prmDial.SetContext(ctx)

	if err := cl.Dial(prmDial); err != nil {
		c.setUnhealthy()
		return false, wasHealthy
	}

	c.clientMutex.Lock()
	c.client = &cl
	c.clientMutex.Unlock()

	if _, err := cl.EndpointInfo(ctx, sdkClient.PrmEndpointInfo{}); err != nil {
		c.setUnhealthy()
		return false, wasHealthy
	}

	c.setHealthy()
	return true, !wasHealthy
}

func (c *clientWrapper) getClient() (*sdkClient.Client, error) {
	c.clientMutex.RLock()
	defer c.clientMutex.RUnlock()
	if c.isHealthy() {
		return c.client, nil
	}
	return nil, errPoolClientUnhealthy
}

// balanceGet invokes sdkClient.BalanceGet parse response status to error and return result as is.
func (c *clientWrapper) balanceGet(ctx context.Context, prm PrmBalanceGet) (accounting.Decimal, error) {
	cl, err := c.getClient()
	if err != nil {
		return accounting.Decimal{}, err
	}

	var cliPrm sdkClient.PrmBalanceGet
	cliPrm.SetAccount(prm.account)

	start := time.Now()
	res, err := cl.BalanceGet(ctx, cliPrm)
	c.incRequests(time.Since(start), methodBalanceGet)
	var st apistatus.Status
	if res != nil {
		st = res.Status()
	}
	if err = c.handleError(st, err); err != nil {
		return accounting.Decimal{}, fmt.Errorf("balance get on client: %w", err)
	}

	return res.Amount(), nil
}

// containerPut invokes sdkClient.ContainerPut parse response status to error and return result as is.
// It also waits for the container to appear on the network.
func (c *clientWrapper) containerPut(ctx context.Context, prm PrmContainerPut) (cid.ID, error) {
	cl, err := c.getClient()
	if err != nil {
		return cid.ID{}, err
	}

	start := time.Now()
	res, err := cl.ContainerPut(ctx, prm.prmClient)
	c.incRequests(time.Since(start), methodContainerPut)
	var st apistatus.Status
	if res != nil {
		st = res.Status()
	}
	if err = c.handleError(st, err); err != nil {
		return cid.ID{}, fmt.Errorf("container put on client: %w", err)
	}

	if !prm.waitParamsSet {
		prm.waitParams.setDefaults()
	}

	idCnr := res.ID()

	err = waitForContainerPresence(ctx, c, idCnr, &prm.waitParams)
	if err = c.handleError(nil, err); err != nil {
		return cid.ID{}, fmt.Errorf("wait container presence on client: %w", err)
	}

	return idCnr, nil
}

// containerGet invokes sdkClient.ContainerGet parse response status to error and return result as is.
func (c *clientWrapper) containerGet(ctx context.Context, prm PrmContainerGet) (container.Container, error) {
	cl, err := c.getClient()
	if err != nil {
		return container.Container{}, err
	}

	var cliPrm sdkClient.PrmContainerGet
	cliPrm.SetContainer(prm.cnrID)

	start := time.Now()
	res, err := cl.ContainerGet(ctx, cliPrm)
	c.incRequests(time.Since(start), methodContainerGet)
	var st apistatus.Status
	if res != nil {
		st = res.Status()
	}
	if err = c.handleError(st, err); err != nil {
		return container.Container{}, fmt.Errorf("container get on client: %w", err)
	}

	return res.Container(), nil
}

// containerList invokes sdkClient.ContainerList parse response status to error and return result as is.
func (c *clientWrapper) containerList(ctx context.Context, prm PrmContainerList) ([]cid.ID, error) {
	cl, err := c.getClient()
	if err != nil {
		return nil, err
	}

	var cliPrm sdkClient.PrmContainerList
	cliPrm.SetAccount(prm.ownerID)

	start := time.Now()
	res, err := cl.ContainerList(ctx, cliPrm)
	c.incRequests(time.Since(start), methodContainerList)
	var st apistatus.Status
	if res != nil {
		st = res.Status()
	}
	if err = c.handleError(st, err); err != nil {
		return nil, fmt.Errorf("container list on client: %w", err)
	}
	return res.Containers(), nil
}

// containerDelete invokes sdkClient.ContainerDelete parse response status to error.
// It also waits for the container to be removed from the network.
func (c *clientWrapper) containerDelete(ctx context.Context, prm PrmContainerDelete) error {
	cl, err := c.getClient()
	if err != nil {
		return err
	}

	var cliPrm sdkClient.PrmContainerDelete
	cliPrm.SetContainer(prm.cnrID)
	if prm.stokenSet {
		cliPrm.WithinSession(prm.stoken)
	}

	start := time.Now()
	res, err := cl.ContainerDelete(ctx, cliPrm)
	c.incRequests(time.Since(start), methodContainerDelete)
	var st apistatus.Status
	if res != nil {
		st = res.Status()
	}
	if err = c.handleError(st, err); err != nil {
		return fmt.Errorf("container delete on client: %w", err)
	}

	if !prm.waitParamsSet {
		prm.waitParams.setDefaults()
	}

	return waitForContainerRemoved(ctx, c, &prm.cnrID, &prm.waitParams)
}

// containerEACL invokes sdkClient.ContainerEACL parse response status to error and return result as is.
func (c *clientWrapper) containerEACL(ctx context.Context, prm PrmContainerEACL) (eacl.Table, error) {
	cl, err := c.getClient()
	if err != nil {
		return eacl.Table{}, err
	}

	var cliPrm sdkClient.PrmContainerEACL
	cliPrm.SetContainer(prm.cnrID)

	start := time.Now()
	res, err := cl.ContainerEACL(ctx, cliPrm)
	c.incRequests(time.Since(start), methodContainerEACL)
	var st apistatus.Status
	if res != nil {
		st = res.Status()
	}
	if err = c.handleError(st, err); err != nil {
		return eacl.Table{}, fmt.Errorf("get eacl on client: %w", err)
	}

	return res.Table(), nil
}

// containerSetEACL invokes sdkClient.ContainerSetEACL parse response status to error.
// It also waits for the EACL to appear on the network.
func (c *clientWrapper) containerSetEACL(ctx context.Context, prm PrmContainerSetEACL) error {
	cl, err := c.getClient()
	if err != nil {
		return err
	}

	var cliPrm sdkClient.PrmContainerSetEACL
	cliPrm.SetTable(prm.table)

	if prm.sessionSet {
		cliPrm.WithinSession(prm.session)
	}

	start := time.Now()
	res, err := cl.ContainerSetEACL(ctx, cliPrm)
	c.incRequests(time.Since(start), methodContainerSetEACL)
	var st apistatus.Status
	if res != nil {
		st = res.Status()
	}
	if err = c.handleError(st, err); err != nil {
		return fmt.Errorf("set eacl on client: %w", err)
	}

	if !prm.waitParamsSet {
		prm.waitParams.setDefaults()
	}

	var cIDp *cid.ID
	if cID, set := prm.table.CID(); set {
		cIDp = &cID
	}

	err = waitForEACLPresence(ctx, c, cIDp, &prm.table, &prm.waitParams)
	if err = c.handleError(nil, err); err != nil {
		return fmt.Errorf("wait eacl presence on client: %w", err)
	}

	return nil
}

// endpointInfo invokes sdkClient.EndpointInfo parse response status to error and return result as is.
func (c *clientWrapper) endpointInfo(ctx context.Context, _ prmEndpointInfo) (netmap.NodeInfo, error) {
	cl, err := c.getClient()
	if err != nil {
		return netmap.NodeInfo{}, err
	}

	start := time.Now()
	res, err := cl.EndpointInfo(ctx, sdkClient.PrmEndpointInfo{})
	c.incRequests(time.Since(start), methodEndpointInfo)
	var st apistatus.Status
	if res != nil {
		st = res.Status()
	}
	if err = c.handleError(st, err); err != nil {
		return netmap.NodeInfo{}, fmt.Errorf("endpoint info on client: %w", err)
	}

	return res.NodeInfo(), nil
}

// networkInfo invokes sdkClient.NetworkInfo parse response status to error and return result as is.
func (c *clientWrapper) networkInfo(ctx context.Context, _ prmNetworkInfo) (netmap.NetworkInfo, error) {
	cl, err := c.getClient()
	if err != nil {
		return netmap.NetworkInfo{}, err
	}

	start := time.Now()
	res, err := cl.NetworkInfo(ctx, sdkClient.PrmNetworkInfo{})
	c.incRequests(time.Since(start), methodNetworkInfo)
	var st apistatus.Status
	if res != nil {
		st = res.Status()
	}
	if err = c.handleError(st, err); err != nil {
		return netmap.NetworkInfo{}, fmt.Errorf("network info on client: %w", err)
	}

	return res.Info(), nil
}

// objectPut writes object to FrostFS.
func (c *clientWrapper) objectPut(ctx context.Context, prm PrmObjectPut) (oid.ID, error) {
	cl, err := c.getClient()
	if err != nil {
		return oid.ID{}, err
	}

	var cliPrm sdkClient.PrmObjectPutInit
	cliPrm.SetCopiesNumber(prm.copiesNumber)
	if prm.stoken != nil {
		cliPrm.WithinSession(*prm.stoken)
	}
	if prm.key != nil {
		cliPrm.UseKey(*prm.key)
	}
	if prm.btoken != nil {
		cliPrm.WithBearerToken(*prm.btoken)
	}

	start := time.Now()
	wObj, err := cl.ObjectPutInit(ctx, cliPrm)
	c.incRequests(time.Since(start), methodObjectPut)
	if err = c.handleError(nil, err); err != nil {
		return oid.ID{}, fmt.Errorf("init writing on API client: %w", err)
	}

	if wObj.WriteHeader(prm.hdr) {
		sz := prm.hdr.PayloadSize()

		if data := prm.hdr.Payload(); len(data) > 0 {
			if prm.payload != nil {
				prm.payload = io.MultiReader(bytes.NewReader(data), prm.payload)
			} else {
				prm.payload = bytes.NewReader(data)
				sz = uint64(len(data))
			}
		}

		if prm.payload != nil {
			const defaultBufferSizePut = 3 << 20 // configure?

			if sz == 0 || sz > defaultBufferSizePut {
				sz = defaultBufferSizePut
			}

			buf := make([]byte, sz)

			var n int

			for {
				n, err = prm.payload.Read(buf)
				if n > 0 {
					start = time.Now()
					successWrite := wObj.WritePayloadChunk(buf[:n])
					c.incRequests(time.Since(start), methodObjectPut)
					if !successWrite {
						break
					}

					continue
				}

				if errors.Is(err, io.EOF) {
					break
				}

				return oid.ID{}, fmt.Errorf("read payload: %w", c.handleError(nil, err))
			}
		}
	}

	res, err := wObj.Close()
	var st apistatus.Status
	if res != nil {
		st = res.Status()
	}
	if err = c.handleError(st, err); err != nil { // here err already carries both status and client errors
		return oid.ID{}, fmt.Errorf("client failure: %w", err)
	}

	return res.StoredObjectID(), nil
}

// objectDelete invokes sdkClient.ObjectDelete parse response status to error.
func (c *clientWrapper) objectDelete(ctx context.Context, prm PrmObjectDelete) error {
	cl, err := c.getClient()
	if err != nil {
		return err
	}

	var cliPrm sdkClient.PrmObjectDelete
	cliPrm.FromContainer(prm.addr.Container())
	cliPrm.ByID(prm.addr.Object())

	if prm.stoken != nil {
		cliPrm.WithinSession(*prm.stoken)
	}

	if prm.btoken != nil {
		cliPrm.WithBearerToken(*prm.btoken)
	}

	if prm.key != nil {
		cliPrm.UseKey(*prm.key)
	}

	start := time.Now()
	res, err := cl.ObjectDelete(ctx, cliPrm)
	c.incRequests(time.Since(start), methodObjectDelete)
	var st apistatus.Status
	if res != nil {
		st = res.Status()
	}
	if err = c.handleError(st, err); err != nil {
		return fmt.Errorf("delete object on client: %w", err)
	}
	return nil
}

// objectGet returns reader for object.
func (c *clientWrapper) objectGet(ctx context.Context, prm PrmObjectGet) (ResGetObject, error) {
	cl, err := c.getClient()
	if err != nil {
		return ResGetObject{}, err
	}

	var cliPrm sdkClient.PrmObjectGet
	cliPrm.FromContainer(prm.addr.Container())
	cliPrm.ByID(prm.addr.Object())

	if prm.stoken != nil {
		cliPrm.WithinSession(*prm.stoken)
	}

	if prm.btoken != nil {
		cliPrm.WithBearerToken(*prm.btoken)
	}

	if prm.key != nil {
		cliPrm.UseKey(*prm.key)
	}

	var res ResGetObject

	rObj, err := cl.ObjectGetInit(ctx, cliPrm)
	if err = c.handleError(nil, err); err != nil {
		return ResGetObject{}, fmt.Errorf("init object reading on client: %w", err)
	}

	start := time.Now()
	successReadHeader := rObj.ReadHeader(&res.Header)
	c.incRequests(time.Since(start), methodObjectGet)
	if !successReadHeader {
		rObjRes, err := rObj.Close()
		var st apistatus.Status
		if rObjRes != nil {
			st = rObjRes.Status()
		}
		err = c.handleError(st, err)
		return res, fmt.Errorf("read header: %w", err)
	}

	res.Payload = &objectReadCloser{
		reader: rObj,
		elapsedTimeCallback: func(elapsed time.Duration) {
			c.incRequests(elapsed, methodObjectGet)
		},
	}

	return res, nil
}

// objectHead invokes sdkClient.ObjectHead parse response status to error and return result as is.
func (c *clientWrapper) objectHead(ctx context.Context, prm PrmObjectHead) (object.Object, error) {
	cl, err := c.getClient()
	if err != nil {
		return object.Object{}, err
	}

	var cliPrm sdkClient.PrmObjectHead
	cliPrm.FromContainer(prm.addr.Container())
	cliPrm.ByID(prm.addr.Object())
	if prm.raw {
		cliPrm.MarkRaw()
	}

	if prm.stoken != nil {
		cliPrm.WithinSession(*prm.stoken)
	}

	if prm.btoken != nil {
		cliPrm.WithBearerToken(*prm.btoken)
	}

	if prm.key != nil {
		cliPrm.UseKey(*prm.key)
	}

	var obj object.Object

	start := time.Now()
	res, err := cl.ObjectHead(ctx, cliPrm)
	c.incRequests(time.Since(start), methodObjectHead)
	var st apistatus.Status
	if res != nil {
		st = res.Status()
	}
	if err = c.handleError(st, err); err != nil {
		return obj, fmt.Errorf("read object header via client: %w", err)
	}
	if !res.ReadHeader(&obj) {
		return obj, errors.New("missing object header in response")
	}

	return obj, nil
}

// objectRange returns object range reader.
func (c *clientWrapper) objectRange(ctx context.Context, prm PrmObjectRange) (ResObjectRange, error) {
	cl, err := c.getClient()
	if err != nil {
		return ResObjectRange{}, err
	}

	var cliPrm sdkClient.PrmObjectRange
	cliPrm.FromContainer(prm.addr.Container())
	cliPrm.ByID(prm.addr.Object())
	cliPrm.SetOffset(prm.off)
	cliPrm.SetLength(prm.ln)

	if prm.stoken != nil {
		cliPrm.WithinSession(*prm.stoken)
	}

	if prm.btoken != nil {
		cliPrm.WithBearerToken(*prm.btoken)
	}

	if prm.key != nil {
		cliPrm.UseKey(*prm.key)
	}

	start := time.Now()
	res, err := cl.ObjectRangeInit(ctx, cliPrm)
	c.incRequests(time.Since(start), methodObjectRange)
	if err = c.handleError(nil, err); err != nil {
		return ResObjectRange{}, fmt.Errorf("init payload range reading on client: %w", err)
	}

	return ResObjectRange{
		payload: res,
		elapsedTimeCallback: func(elapsed time.Duration) {
			c.incRequests(elapsed, methodObjectRange)
		},
	}, nil
}

// objectSearch invokes sdkClient.ObjectSearchInit parse response status to error and return result as is.
func (c *clientWrapper) objectSearch(ctx context.Context, prm PrmObjectSearch) (ResObjectSearch, error) {
	cl, err := c.getClient()
	if err != nil {
		return ResObjectSearch{}, err
	}

	var cliPrm sdkClient.PrmObjectSearch

	cliPrm.InContainer(prm.cnrID)
	cliPrm.SetFilters(prm.filters)

	if prm.stoken != nil {
		cliPrm.WithinSession(*prm.stoken)
	}

	if prm.btoken != nil {
		cliPrm.WithBearerToken(*prm.btoken)
	}

	if prm.key != nil {
		cliPrm.UseKey(*prm.key)
	}

	res, err := cl.ObjectSearchInit(ctx, cliPrm)
	if err = c.handleError(nil, err); err != nil {
		return ResObjectSearch{}, fmt.Errorf("init object searching on client: %w", err)
	}

	return ResObjectSearch{r: res}, nil
}

// sessionCreate invokes sdkClient.SessionCreate parse response status to error and return result as is.
func (c *clientWrapper) sessionCreate(ctx context.Context, prm prmCreateSession) (resCreateSession, error) {
	cl, err := c.getClient()
	if err != nil {
		return resCreateSession{}, err
	}

	var cliPrm sdkClient.PrmSessionCreate
	cliPrm.SetExp(prm.exp)
	cliPrm.UseKey(prm.key)

	start := time.Now()
	res, err := cl.SessionCreate(ctx, cliPrm)
	c.incRequests(time.Since(start), methodSessionCreate)
	var st apistatus.Status
	if res != nil {
		st = res.Status()
	}
	if err = c.handleError(st, err); err != nil {
		return resCreateSession{}, fmt.Errorf("session creation on client: %w", err)
	}

	return resCreateSession{
		id:         res.ID(),
		sessionKey: res.PublicKey(),
	}, nil
}

func (c *clientStatusMonitor) isHealthy() bool {
	return c.healthy.Load()
}

func (c *clientStatusMonitor) setHealthy() {
	c.healthy.Store(true)
}

func (c *clientStatusMonitor) setUnhealthy() {
	c.healthy.Store(false)
}

func (c *clientStatusMonitor) address() string {
	return c.addr
}

func (c *clientStatusMonitor) incErrorRate() {
	c.mu.Lock()
	c.currentErrorCount++
	c.overallErrorCount++

	thresholdReached := c.currentErrorCount >= c.errorThreshold
	if thresholdReached {
		c.setUnhealthy()
		c.currentErrorCount = 0
	}
	c.mu.Unlock()

	if thresholdReached && c.logger != nil {
		c.logger.Warn("error threshold reached",
			zap.String("address", c.addr), zap.Uint32("threshold", c.errorThreshold))
	}
}

func (c *clientStatusMonitor) currentErrorRate() uint32 {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.currentErrorCount
}

func (c *clientStatusMonitor) overallErrorRate() uint64 {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.overallErrorCount
}

func (c *clientStatusMonitor) methodsStatus() []statusSnapshot {
	result := make([]statusSnapshot, len(c.methods))
	for i, val := range c.methods {
		result[i] = val.snapshot()
	}

	return result
}

func (c *clientWrapper) incRequests(elapsed time.Duration, method MethodIndex) {
	methodStat := c.methods[method]
	methodStat.incRequests(elapsed)
	if c.prm.poolRequestInfoCallback != nil {
		c.prm.poolRequestInfoCallback(RequestInfo{
			Address: c.prm.address,
			Method:  method,
			Elapsed: elapsed,
		})
	}
}

func (c *clientStatusMonitor) handleError(st apistatus.Status, err error) error {
	if err != nil {
		// non-status logic error that could be returned
		// from the SDK client; should not be considered
		// as a connection error
		var siErr *object.SplitInfoError
		if !errors.As(err, &siErr) {
			c.incErrorRate()
		}

		return err
	}

	err = apistatus.ErrFromStatus(st)
	switch err.(type) {
	case apistatus.ServerInternal, *apistatus.ServerInternal,
		apistatus.WrongMagicNumber, *apistatus.WrongMagicNumber,
		apistatus.SignatureVerification, *apistatus.SignatureVerification,
		apistatus.NodeUnderMaintenance, *apistatus.NodeUnderMaintenance:
		c.incErrorRate()
	}

	return err
}

// clientBuilder is a type alias of client constructors which open connection
// to the given endpoint.
type clientBuilder = func(endpoint string) client

// RequestInfo groups info about pool request.
type RequestInfo struct {
	Address string
	Method  MethodIndex
	Elapsed time.Duration
}

// InitParameters contains values used to initialize connection Pool.
type InitParameters struct {
	key                       *ecdsa.PrivateKey
	logger                    *zap.Logger
	nodeDialTimeout           time.Duration
	nodeStreamTimeout         time.Duration
	healthcheckTimeout        time.Duration
	clientRebalanceInterval   time.Duration
	sessionExpirationDuration uint64
	errorThreshold            uint32
	nodeParams                []NodeParam
	requestCallback           func(RequestInfo)

	clientBuilder clientBuilder
}

// SetKey specifies default key to be used for the protocol communication by default.
func (x *InitParameters) SetKey(key *ecdsa.PrivateKey) {
	x.key = key
}

// SetLogger specifies logger.
func (x *InitParameters) SetLogger(logger *zap.Logger) {
	x.logger = logger
}

// SetNodeDialTimeout specifies the timeout for connection to be established.
func (x *InitParameters) SetNodeDialTimeout(timeout time.Duration) {
	x.nodeDialTimeout = timeout
}

// SetNodeStreamTimeout specifies the timeout for individual operations in streaming RPC.
func (x *InitParameters) SetNodeStreamTimeout(timeout time.Duration) {
	x.nodeStreamTimeout = timeout
}

// SetHealthcheckTimeout specifies the timeout for request to node to decide if it is alive.
//
// See also Pool.Dial.
func (x *InitParameters) SetHealthcheckTimeout(timeout time.Duration) {
	x.healthcheckTimeout = timeout
}

// SetClientRebalanceInterval specifies the interval for updating nodes health status.
//
// See also Pool.Dial.
func (x *InitParameters) SetClientRebalanceInterval(interval time.Duration) {
	x.clientRebalanceInterval = interval
}

// SetSessionExpirationDuration specifies the session token lifetime in epochs.
func (x *InitParameters) SetSessionExpirationDuration(expirationDuration uint64) {
	x.sessionExpirationDuration = expirationDuration
}

// SetErrorThreshold specifies the number of errors on connection after which node is considered as unhealthy.
func (x *InitParameters) SetErrorThreshold(threshold uint32) {
	x.errorThreshold = threshold
}

// SetRequestCallback makes the pool client to pass RequestInfo for each
// request to f. Nil (default) means ignore RequestInfo.
func (x *InitParameters) SetRequestCallback(f func(RequestInfo)) {
	x.requestCallback = f
}

// AddNode append information about the node to which you want to connect.
func (x *InitParameters) AddNode(nodeParam NodeParam) {
	x.nodeParams = append(x.nodeParams, nodeParam)
}

// setClientBuilder sets clientBuilder used for client construction.
// Wraps setClientBuilderContext without a context.
func (x *InitParameters) setClientBuilder(builder clientBuilder) {
	x.clientBuilder = builder
}

// isMissingClientBuilder checks if client constructor was not specified.
func (x *InitParameters) isMissingClientBuilder() bool {
	return x.clientBuilder == nil
}

type rebalanceParameters struct {
	nodesParams               []*nodesParam
	nodeRequestTimeout        time.Duration
	clientRebalanceInterval   time.Duration
	sessionExpirationDuration uint64
}

type nodesParam struct {
	priority  int
	addresses []string
	weights   []float64
}

// NodeParam groups parameters of remote node.
type NodeParam struct {
	priority int
	address  string
	weight   float64
}

// NewNodeParam creates NodeParam using parameters.
func NewNodeParam(priority int, address string, weight float64) (prm NodeParam) {
	prm.SetPriority(priority)
	prm.SetAddress(address)
	prm.SetWeight(weight)

	return
}

// SetPriority specifies priority of the node.
// Negative value is allowed. In the result node groups
// with the same priority will be sorted by descent.
func (x *NodeParam) SetPriority(priority int) {
	x.priority = priority
}

// SetAddress specifies address of the node.
func (x *NodeParam) SetAddress(address string) {
	x.address = address
}

// SetWeight specifies weight of the node.
func (x *NodeParam) SetWeight(weight float64) {
	x.weight = weight
}

// WaitParams contains parameters used in polling is a something applied on FrostFS network.
type WaitParams struct {
	timeout      time.Duration
	pollInterval time.Duration
}

// SetTimeout specifies the time to wait for the operation to complete.
func (x *WaitParams) SetTimeout(timeout time.Duration) {
	x.timeout = timeout
}

// SetPollInterval specifies the interval, once it will check the completion of the operation.
func (x *WaitParams) SetPollInterval(tick time.Duration) {
	x.pollInterval = tick
}

func (x *WaitParams) setDefaults() {
	x.timeout = 120 * time.Second
	x.pollInterval = 5 * time.Second
}

// checkForPositive panics if any of the wait params isn't positive.
func (x *WaitParams) checkForPositive() {
	if x.timeout <= 0 || x.pollInterval <= 0 {
		panic("all wait params must be positive")
	}
}

type prmContext struct {
	defaultSession bool
	verb           session.ObjectVerb
	cnr            cid.ID

	objSet bool
	objs   []oid.ID
}

func (x *prmContext) useDefaultSession() {
	x.defaultSession = true
}

func (x *prmContext) useContainer(cnr cid.ID) {
	x.cnr = cnr
}

func (x *prmContext) useObjects(ids []oid.ID) {
	x.objs = ids
	x.objSet = true
}

func (x *prmContext) useAddress(addr oid.Address) {
	x.cnr = addr.Container()
	x.objs = []oid.ID{addr.Object()}
	x.objSet = true
}

func (x *prmContext) useVerb(verb session.ObjectVerb) {
	x.verb = verb
}

type prmCommon struct {
	key    *ecdsa.PrivateKey
	btoken *bearer.Token
	stoken *session.Object
}

// UseKey specifies private key to sign the requests.
// If key is not provided, then Pool default key is used.
func (x *prmCommon) UseKey(key *ecdsa.PrivateKey) {
	x.key = key
}

// UseBearer attaches bearer token to be used for the operation.
func (x *prmCommon) UseBearer(token bearer.Token) {
	x.btoken = &token
}

// UseSession specifies session within which operation should be performed.
func (x *prmCommon) UseSession(token session.Object) {
	x.stoken = &token
}

// PrmObjectPut groups parameters of PutObject operation.
type PrmObjectPut struct {
	prmCommon

	hdr object.Object

	payload io.Reader

	copiesNumber uint32
}

// SetHeader specifies header of the object.
func (x *PrmObjectPut) SetHeader(hdr object.Object) {
	x.hdr = hdr
}

// SetPayload specifies payload of the object.
func (x *PrmObjectPut) SetPayload(payload io.Reader) {
	x.payload = payload
}

// SetCopiesNumber sets number of object copies that is enough to consider put successful.
// Zero means using default behavior.
func (x *PrmObjectPut) SetCopiesNumber(copiesNumber uint32) {
	x.copiesNumber = copiesNumber
}

// PrmObjectDelete groups parameters of DeleteObject operation.
type PrmObjectDelete struct {
	prmCommon

	addr oid.Address
}

// SetAddress specifies FrostFS address of the object.
func (x *PrmObjectDelete) SetAddress(addr oid.Address) {
	x.addr = addr
}

// PrmObjectGet groups parameters of GetObject operation.
type PrmObjectGet struct {
	prmCommon

	addr oid.Address
}

// SetAddress specifies FrostFS address of the object.
func (x *PrmObjectGet) SetAddress(addr oid.Address) {
	x.addr = addr
}

// PrmObjectHead groups parameters of HeadObject operation.
type PrmObjectHead struct {
	prmCommon

	addr oid.Address
	raw  bool
}

// SetAddress specifies FrostFS address of the object.
func (x *PrmObjectHead) SetAddress(addr oid.Address) {
	x.addr = addr
}

// MarkRaw marks an intent to read physically stored object.
func (x *PrmObjectHead) MarkRaw() {
	x.raw = true
}

// PrmObjectRange groups parameters of RangeObject operation.
type PrmObjectRange struct {
	prmCommon

	addr    oid.Address
	off, ln uint64
}

// SetAddress specifies FrostFS address of the object.
func (x *PrmObjectRange) SetAddress(addr oid.Address) {
	x.addr = addr
}

// SetOffset sets offset of the payload range to be read.
func (x *PrmObjectRange) SetOffset(offset uint64) {
	x.off = offset
}

// SetLength sets length of the payload range to be read.
func (x *PrmObjectRange) SetLength(length uint64) {
	x.ln = length
}

// PrmObjectSearch groups parameters of SearchObjects operation.
type PrmObjectSearch struct {
	prmCommon

	cnrID   cid.ID
	filters object.SearchFilters
}

// SetContainerID specifies the container in which to look for objects.
func (x *PrmObjectSearch) SetContainerID(cnrID cid.ID) {
	x.cnrID = cnrID
}

// SetFilters specifies filters by which to select objects.
func (x *PrmObjectSearch) SetFilters(filters object.SearchFilters) {
	x.filters = filters
}

// PrmContainerPut groups parameters of PutContainer operation.
type PrmContainerPut struct {
	prmClient sdkClient.PrmContainerPut

	waitParams    WaitParams
	waitParamsSet bool
}

// SetContainer container structure to be used as a parameter of the base
// client's operation.
//
// See github.com/TrueCloudLab/frostfs-sdk-go/client.PrmContainerPut.SetContainer.
func (x *PrmContainerPut) SetContainer(cnr container.Container) {
	x.prmClient.SetContainer(cnr)
}

// WithinSession specifies session to be used as a parameter of the base
// client's operation.
//
// See github.com/TrueCloudLab/frostfs-sdk-go/client.PrmContainerPut.WithinSession.
func (x *PrmContainerPut) WithinSession(s session.Container) {
	x.prmClient.WithinSession(s)
}

// SetWaitParams specifies timeout params to complete operation.
// If not provided the default one will be used.
// Panics if any of the wait params isn't positive.
func (x *PrmContainerPut) SetWaitParams(waitParams WaitParams) {
	waitParams.checkForPositive()
	x.waitParams = waitParams
	x.waitParamsSet = true
}

// PrmContainerGet groups parameters of GetContainer operation.
type PrmContainerGet struct {
	cnrID cid.ID
}

// SetContainerID specifies identifier of the container to be read.
func (x *PrmContainerGet) SetContainerID(cnrID cid.ID) {
	x.cnrID = cnrID
}

// PrmContainerList groups parameters of ListContainers operation.
type PrmContainerList struct {
	ownerID user.ID
}

// SetOwnerID specifies identifier of the FrostFS account to list the containers.
func (x *PrmContainerList) SetOwnerID(ownerID user.ID) {
	x.ownerID = ownerID
}

// PrmContainerDelete groups parameters of DeleteContainer operation.
type PrmContainerDelete struct {
	cnrID cid.ID

	stoken    session.Container
	stokenSet bool

	waitParams    WaitParams
	waitParamsSet bool
}

// SetContainerID specifies identifier of the FrostFS container to be removed.
func (x *PrmContainerDelete) SetContainerID(cnrID cid.ID) {
	x.cnrID = cnrID
}

// SetSessionToken specifies session within which operation should be performed.
func (x *PrmContainerDelete) SetSessionToken(token session.Container) {
	x.stoken = token
	x.stokenSet = true
}

// SetWaitParams specifies timeout params to complete operation.
// If not provided the default one will be used.
// Panics if any of the wait params isn't positive.
func (x *PrmContainerDelete) SetWaitParams(waitParams WaitParams) {
	waitParams.checkForPositive()
	x.waitParams = waitParams
	x.waitParamsSet = true
}

// PrmContainerEACL groups parameters of GetEACL operation.
type PrmContainerEACL struct {
	cnrID cid.ID
}

// SetContainerID specifies identifier of the FrostFS container to read the eACL table.
func (x *PrmContainerEACL) SetContainerID(cnrID cid.ID) {
	x.cnrID = cnrID
}

// PrmContainerSetEACL groups parameters of SetEACL operation.
type PrmContainerSetEACL struct {
	table eacl.Table

	sessionSet bool
	session    session.Container

	waitParams    WaitParams
	waitParamsSet bool
}

// SetTable sets structure of container's extended ACL to be used as a
// parameter of the base client's operation.
//
// See github.com/TrueCloudLab/frostfs-sdk-go/client.PrmContainerSetEACL.SetTable.
func (x *PrmContainerSetEACL) SetTable(table eacl.Table) {
	x.table = table
}

// WithinSession specifies session to be used as a parameter of the base
// client's operation.
//
// See github.com/TrueCloudLab/frostfs-sdk-go/client.PrmContainerSetEACL.WithinSession.
func (x *PrmContainerSetEACL) WithinSession(s session.Container) {
	x.session = s
	x.sessionSet = true
}

// SetWaitParams specifies timeout params to complete operation.
// If not provided the default one will be used.
// Panics if any of the wait params isn't positive.
func (x *PrmContainerSetEACL) SetWaitParams(waitParams WaitParams) {
	waitParams.checkForPositive()
	x.waitParams = waitParams
	x.waitParamsSet = true
}

// PrmBalanceGet groups parameters of Balance operation.
type PrmBalanceGet struct {
	account user.ID
}

// SetAccount specifies identifier of the FrostFS account for which the balance is requested.
func (x *PrmBalanceGet) SetAccount(id user.ID) {
	x.account = id
}

// prmEndpointInfo groups parameters of sessionCreate operation.
type prmCreateSession struct {
	exp uint64
	key ecdsa.PrivateKey
}

// setExp sets number of the last FrostFS epoch in the lifetime of the session after which it will be expired.
func (x *prmCreateSession) setExp(exp uint64) {
	x.exp = exp
}

// useKey specifies owner private key for session token.
// If key is not provided, then Pool default key is used.
func (x *prmCreateSession) useKey(key ecdsa.PrivateKey) {
	x.key = key
}

// prmEndpointInfo groups parameters of endpointInfo operation.
type prmEndpointInfo struct{}

// prmNetworkInfo groups parameters of networkInfo operation.
type prmNetworkInfo struct{}

// resCreateSession groups resulting values of sessionCreate operation.
type resCreateSession struct {
	id []byte

	sessionKey []byte
}

// Pool represents virtual connection to the FrostFS network to communicate
// with multiple FrostFS servers without thinking about switching between servers
// due to load balancing proportions or their unavailability.
// It is designed to provide a convenient abstraction from the multiple sdkClient.client types.
//
// Pool can be created and initialized using NewPool function.
// Before executing the FrostFS operations using the Pool, connection to the
// servers MUST BE correctly established (see Dial method).
// Using the Pool before connecting have been established can lead to a panic.
// After the work, the Pool SHOULD BE closed (see Close method): it frees internal
// and system resources which were allocated for the period of work of the Pool.
// Calling Dial/Close methods during the communication process step strongly discouraged
// as it leads to undefined behavior.
//
// Each method which produces a FrostFS API call may return an error.
// Status of underlying server response is casted to built-in error instance.
// Certain statuses can be checked using `sdkClient` and standard `errors` packages.
// Note that package provides some helper functions to work with status returns
// (e.g. sdkClient.IsErrContainerNotFound, sdkClient.IsErrObjectNotFound).
//
// See pool package overview to get some examples.
type Pool struct {
	innerPools      []*innerPool
	key             *ecdsa.PrivateKey
	cancel          context.CancelFunc
	closedCh        chan struct{}
	cache           *sessionCache
	stokenDuration  uint64
	rebalanceParams rebalanceParameters
	clientBuilder   clientBuilder
	logger          *zap.Logger
}

type innerPool struct {
	lock    sync.RWMutex
	sampler *sampler
	clients []client
}

const (
	defaultSessionTokenExpirationDuration = 100 // in blocks
	defaultErrorThreshold                 = 100

	defaultRebalanceInterval  = 15 * time.Second
	defaultHealthcheckTimeout = 4 * time.Second
	defaultDialTimeout        = 5 * time.Second
	defaultStreamTimeout      = 10 * time.Second
)

// NewPool creates connection pool using parameters.
func NewPool(options InitParameters) (*Pool, error) {
	if options.key == nil {
		return nil, fmt.Errorf("missed required parameter 'Key'")
	}

	nodesParams, err := adjustNodeParams(options.nodeParams)
	if err != nil {
		return nil, err
	}

	cache, err := newCache()
	if err != nil {
		return nil, fmt.Errorf("couldn't create cache: %w", err)
	}

	fillDefaultInitParams(&options, cache)

	pool := &Pool{
		key:            options.key,
		cache:          cache,
		logger:         options.logger,
		stokenDuration: options.sessionExpirationDuration,
		rebalanceParams: rebalanceParameters{
			nodesParams:               nodesParams,
			nodeRequestTimeout:        options.healthcheckTimeout,
			clientRebalanceInterval:   options.clientRebalanceInterval,
			sessionExpirationDuration: options.sessionExpirationDuration,
		},
		clientBuilder: options.clientBuilder,
	}

	return pool, nil
}

// Dial establishes a connection to the servers from the FrostFS network.
// It also starts a routine that checks the health of the nodes and
// updates the weights of the nodes for balancing.
// Returns an error describing failure reason.
//
// If failed, the Pool SHOULD NOT be used.
//
// See also InitParameters.SetClientRebalanceInterval.
func (p *Pool) Dial(ctx context.Context) error {
	inner := make([]*innerPool, len(p.rebalanceParams.nodesParams))
	var atLeastOneHealthy bool

	for i, params := range p.rebalanceParams.nodesParams {
		clients := make([]client, len(params.weights))
		for j, addr := range params.addresses {
			clients[j] = p.clientBuilder(addr)
			if err := clients[j].dial(ctx); err != nil {
				p.log(zap.WarnLevel, "failed to build client", zap.String("address", addr), zap.Error(err))
				continue
			}

			var st session.Object
			err := initSessionForDuration(ctx, &st, clients[j], p.rebalanceParams.sessionExpirationDuration, *p.key)
			if err != nil {
				clients[j].setUnhealthy()
				p.log(zap.WarnLevel, "failed to create frostfs session token for client",
					zap.String("address", addr), zap.Error(err))
				continue
			}

			_ = p.cache.Put(formCacheKey(addr, p.key), st)
			atLeastOneHealthy = true
		}
		source := rand.NewSource(time.Now().UnixNano())
		sampl := newSampler(params.weights, source)

		inner[i] = &innerPool{
			sampler: sampl,
			clients: clients,
		}
	}

	if !atLeastOneHealthy {
		return fmt.Errorf("at least one node must be healthy")
	}

	ctx, cancel := context.WithCancel(ctx)
	p.cancel = cancel
	p.closedCh = make(chan struct{})
	p.innerPools = inner

	go p.startRebalance(ctx)
	return nil
}

func (p *Pool) log(level zapcore.Level, msg string, fields ...zap.Field) {
	if p.logger == nil {
		return
	}

	p.logger.Log(level, msg, fields...)
}

func fillDefaultInitParams(params *InitParameters, cache *sessionCache) {
	if params.sessionExpirationDuration == 0 {
		params.sessionExpirationDuration = defaultSessionTokenExpirationDuration
	}

	if params.errorThreshold == 0 {
		params.errorThreshold = defaultErrorThreshold
	}

	if params.clientRebalanceInterval <= 0 {
		params.clientRebalanceInterval = defaultRebalanceInterval
	}

	if params.healthcheckTimeout <= 0 {
		params.healthcheckTimeout = defaultHealthcheckTimeout
	}

	if params.nodeDialTimeout <= 0 {
		params.nodeDialTimeout = defaultDialTimeout
	}

	if params.nodeStreamTimeout <= 0 {
		params.nodeStreamTimeout = defaultStreamTimeout
	}

	if params.isMissingClientBuilder() {
		params.setClientBuilder(func(addr string) client {
			var prm wrapperPrm
			prm.setAddress(addr)
			prm.setKey(*params.key)
			prm.setDialTimeout(params.nodeDialTimeout)
			prm.setStreamTimeout(params.nodeStreamTimeout)
			prm.setErrorThreshold(params.errorThreshold)
			prm.setPoolRequestCallback(params.requestCallback)
			prm.setResponseInfoCallback(func(info sdkClient.ResponseMetaInfo) error {
				cache.updateEpoch(info.Epoch())
				return nil
			})
			return newWrapper(prm)
		})
	}
}

func adjustNodeParams(nodeParams []NodeParam) ([]*nodesParam, error) {
	if len(nodeParams) == 0 {
		return nil, errors.New("no FrostFS peers configured")
	}

	nodesParamsMap := make(map[int]*nodesParam)
	for _, param := range nodeParams {
		nodes, ok := nodesParamsMap[param.priority]
		if !ok {
			nodes = &nodesParam{priority: param.priority}
		}
		nodes.addresses = append(nodes.addresses, param.address)
		nodes.weights = append(nodes.weights, param.weight)
		nodesParamsMap[param.priority] = nodes
	}

	nodesParams := make([]*nodesParam, 0, len(nodesParamsMap))
	for _, nodes := range nodesParamsMap {
		nodes.weights = adjustWeights(nodes.weights)
		nodesParams = append(nodesParams, nodes)
	}

	sort.Slice(nodesParams, func(i, j int) bool {
		return nodesParams[i].priority < nodesParams[j].priority
	})

	return nodesParams, nil
}

// startRebalance runs loop to monitor connection healthy status.
func (p *Pool) startRebalance(ctx context.Context) {
	ticker := time.NewTimer(p.rebalanceParams.clientRebalanceInterval)
	buffers := make([][]float64, len(p.rebalanceParams.nodesParams))
	for i, params := range p.rebalanceParams.nodesParams {
		buffers[i] = make([]float64, len(params.weights))
	}

	for {
		select {
		case <-ctx.Done():
			close(p.closedCh)
			return
		case <-ticker.C:
			p.updateNodesHealth(ctx, buffers)
			ticker.Reset(p.rebalanceParams.clientRebalanceInterval)
		}
	}
}

func (p *Pool) updateNodesHealth(ctx context.Context, buffers [][]float64) {
	wg := sync.WaitGroup{}
	for i, inner := range p.innerPools {
		wg.Add(1)

		bufferWeights := buffers[i]
		go func(i int, innerPool *innerPool) {
			defer wg.Done()
			p.updateInnerNodesHealth(ctx, i, bufferWeights)
		}(i, inner)
	}
	wg.Wait()
}

func (p *Pool) updateInnerNodesHealth(ctx context.Context, i int, bufferWeights []float64) {
	if i > len(p.innerPools)-1 {
		return
	}
	pool := p.innerPools[i]
	options := p.rebalanceParams

	healthyChanged := atomic.NewBool(false)
	wg := sync.WaitGroup{}

	for j, cli := range pool.clients {
		wg.Add(1)
		go func(j int, cli client) {
			defer wg.Done()

			tctx, c := context.WithTimeout(ctx, options.nodeRequestTimeout)
			defer c()

			healthy, changed := cli.restartIfUnhealthy(tctx)
			if healthy {
				bufferWeights[j] = options.nodesParams[i].weights[j]
			} else {
				bufferWeights[j] = 0
				p.cache.DeleteByPrefix(cli.address())
			}

			if changed {
				p.log(zap.DebugLevel, "health has changed",
					zap.String("address", cli.address()), zap.Bool("healthy", healthy))
				healthyChanged.Store(true)
			}
		}(j, cli)
	}
	wg.Wait()

	if healthyChanged.Load() {
		probabilities := adjustWeights(bufferWeights)
		source := rand.NewSource(time.Now().UnixNano())
		pool.lock.Lock()
		pool.sampler = newSampler(probabilities, source)
		pool.lock.Unlock()
	}
}

func adjustWeights(weights []float64) []float64 {
	adjusted := make([]float64, len(weights))
	sum := 0.0
	for _, weight := range weights {
		sum += weight
	}
	if sum > 0 {
		for i, weight := range weights {
			adjusted[i] = weight / sum
		}
	}

	return adjusted
}

func (p *Pool) connection() (client, error) {
	for _, inner := range p.innerPools {
		cp, err := inner.connection()
		if err == nil {
			return cp, nil
		}
	}

	return nil, errors.New("no healthy client")
}

func (p *innerPool) connection() (client, error) {
	p.lock.RLock() // need lock because of using p.sampler
	defer p.lock.RUnlock()
	if len(p.clients) == 1 {
		cp := p.clients[0]
		if cp.isHealthy() {
			return cp, nil
		}
		return nil, errors.New("no healthy client")
	}
	attempts := 3 * len(p.clients)
	for k := 0; k < attempts; k++ {
		i := p.sampler.Next()
		if cp := p.clients[i]; cp.isHealthy() {
			return cp, nil
		}
	}

	return nil, errors.New("no healthy client")
}

func formCacheKey(address string, key *ecdsa.PrivateKey) string {
	k := keys.PrivateKey{PrivateKey: *key}
	return address + k.String()
}

func (p *Pool) checkSessionTokenErr(err error, address string) bool {
	if err == nil {
		return false
	}

	if sdkClient.IsErrSessionNotFound(err) || sdkClient.IsErrSessionExpired(err) {
		p.cache.DeleteByPrefix(address)
		return true
	}

	return false
}

func initSessionForDuration(ctx context.Context, dst *session.Object, c client, dur uint64, ownerKey ecdsa.PrivateKey) error {
	ni, err := c.networkInfo(ctx, prmNetworkInfo{})
	if err != nil {
		return err
	}

	epoch := ni.CurrentEpoch()

	var exp uint64
	if math.MaxUint64-epoch < dur {
		exp = math.MaxUint64
	} else {
		exp = epoch + dur
	}
	var prm prmCreateSession
	prm.setExp(exp)
	prm.useKey(ownerKey)

	res, err := c.sessionCreate(ctx, prm)
	if err != nil {
		return err
	}

	var id uuid.UUID

	err = id.UnmarshalBinary(res.id)
	if err != nil {
		return fmt.Errorf("invalid session token ID: %w", err)
	}

	var key frostfsecdsa.PublicKey

	err = key.Decode(res.sessionKey)
	if err != nil {
		return fmt.Errorf("invalid public session key: %w", err)
	}

	dst.SetID(id)
	dst.SetAuthKey(&key)
	dst.SetExp(exp)

	return nil
}

type callContext struct {
	// base context for RPC
	context.Context

	client client

	// client endpoint
	endpoint string

	// request signer
	key *ecdsa.PrivateKey

	// flag to open default session if session token is missing
	sessionDefault bool
	sessionTarget  func(session.Object)
	sessionVerb    session.ObjectVerb
	sessionCnr     cid.ID
	sessionObjSet  bool
	sessionObjs    []oid.ID
}

func (p *Pool) initCallContext(ctx *callContext, cfg prmCommon, prmCtx prmContext) error {
	cp, err := p.connection()
	if err != nil {
		return err
	}

	ctx.key = cfg.key
	if ctx.key == nil {
		// use pool key if caller didn't specify its own
		ctx.key = p.key
	}

	ctx.endpoint = cp.address()
	ctx.client = cp

	if ctx.sessionTarget != nil && cfg.stoken != nil {
		ctx.sessionTarget(*cfg.stoken)
	}

	// note that we don't override session provided by the caller
	ctx.sessionDefault = cfg.stoken == nil && prmCtx.defaultSession
	if ctx.sessionDefault {
		ctx.sessionVerb = prmCtx.verb
		ctx.sessionCnr = prmCtx.cnr
		ctx.sessionObjSet = prmCtx.objSet
		ctx.sessionObjs = prmCtx.objs
	}

	return err
}

// opens new session or uses cached one.
// Must be called only on initialized callContext with set sessionTarget.
func (p *Pool) openDefaultSession(ctx *callContext) error {
	cacheKey := formCacheKey(ctx.endpoint, ctx.key)

	tok, ok := p.cache.Get(cacheKey)
	if !ok {
		// init new session
		err := initSessionForDuration(ctx, &tok, ctx.client, p.stokenDuration, *ctx.key)
		if err != nil {
			return fmt.Errorf("session API client: %w", err)
		}

		// cache the opened session
		p.cache.Put(cacheKey, tok)
	}

	tok.ForVerb(ctx.sessionVerb)
	tok.BindContainer(ctx.sessionCnr)

	if ctx.sessionObjSet {
		tok.LimitByObjects(ctx.sessionObjs...)
	}

	// sign the token
	if err := tok.Sign(*ctx.key); err != nil {
		return fmt.Errorf("sign token of the opened session: %w", err)
	}

	ctx.sessionTarget(tok)

	return nil
}

// opens default session (if sessionDefault is set), and calls f. If f returns
// session-related error then cached token is removed.
func (p *Pool) call(ctx *callContext, f func() error) error {
	var err error

	if ctx.sessionDefault {
		err = p.openDefaultSession(ctx)
		if err != nil {
			return fmt.Errorf("open default session: %w", err)
		}
	}

	err = f()
	_ = p.checkSessionTokenErr(err, ctx.endpoint)

	return err
}

// fillAppropriateKey use pool key if caller didn't specify its own.
func (p *Pool) fillAppropriateKey(prm *prmCommon) {
	if prm.key == nil {
		prm.key = p.key
	}
}

// PutObject writes an object through a remote server using FrostFS API protocol.
//
// Main return value MUST NOT be processed on an erroneous return.
func (p *Pool) PutObject(ctx context.Context, prm PrmObjectPut) (oid.ID, error) {
	cnr, _ := prm.hdr.ContainerID()

	var prmCtx prmContext
	prmCtx.useDefaultSession()
	prmCtx.useVerb(session.VerbObjectPut)
	prmCtx.useContainer(cnr)

	p.fillAppropriateKey(&prm.prmCommon)

	var ctxCall callContext

	ctxCall.Context = ctx

	if err := p.initCallContext(&ctxCall, prm.prmCommon, prmCtx); err != nil {
		return oid.ID{}, fmt.Errorf("init call context: %w", err)
	}

	if ctxCall.sessionDefault {
		ctxCall.sessionTarget = prm.UseSession
		if err := p.openDefaultSession(&ctxCall); err != nil {
			return oid.ID{}, fmt.Errorf("open default session: %w", err)
		}
	}

	id, err := ctxCall.client.objectPut(ctx, prm)
	if err != nil {
		// removes session token from cache in case of token error
		p.checkSessionTokenErr(err, ctxCall.endpoint)
		return id, fmt.Errorf("init writing on API client: %w", err)
	}

	return id, nil
}

// DeleteObject marks an object for deletion from the container using FrostFS API protocol.
// As a marker, a special unit called a tombstone is placed in the container.
// It confirms the user's intent to delete the object, and is itself a container object.
// Explicit deletion is done asynchronously, and is generally not guaranteed.
func (p *Pool) DeleteObject(ctx context.Context, prm PrmObjectDelete) error {
	var prmCtx prmContext
	prmCtx.useDefaultSession()
	prmCtx.useVerb(session.VerbObjectDelete)
	prmCtx.useAddress(prm.addr)

	if prm.stoken == nil { // collect phy objects only if we are about to open default session
		var tokens relations.Tokens
		tokens.Bearer = prm.btoken

		relatives, err := relations.ListAllRelations(ctx, p, prm.addr.Container(), prm.addr.Object(), tokens)
		if err != nil {
			return fmt.Errorf("failed to collect relatives: %w", err)
		}

		if len(relatives) != 0 {
			prmCtx.useContainer(prm.addr.Container())
			prmCtx.useObjects(append(relatives, prm.addr.Object()))
		}
	}

	p.fillAppropriateKey(&prm.prmCommon)

	var cc callContext

	cc.Context = ctx
	cc.sessionTarget = prm.UseSession

	err := p.initCallContext(&cc, prm.prmCommon, prmCtx)
	if err != nil {
		return err
	}

	return p.call(&cc, func() error {
		if err = cc.client.objectDelete(ctx, prm); err != nil {
			return fmt.Errorf("remove object via client: %w", err)
		}

		return nil
	})
}

type objectReadCloser struct {
	reader              *sdkClient.ObjectReader
	elapsedTimeCallback func(time.Duration)
}

// Read implements io.Reader of the object payload.
func (x *objectReadCloser) Read(p []byte) (int, error) {
	start := time.Now()
	n, err := x.reader.Read(p)
	x.elapsedTimeCallback(time.Since(start))
	return n, err
}

// Close implements io.Closer of the object payload.
func (x *objectReadCloser) Close() error {
	_, err := x.reader.Close()
	return err
}

// ResGetObject is designed to provide object header nad read one object payload from FrostFS system.
type ResGetObject struct {
	Header object.Object

	Payload io.ReadCloser
}

// GetObject reads object header and initiates reading an object payload through a remote server using FrostFS API protocol.
//
// Main return value MUST NOT be processed on an erroneous return.
func (p *Pool) GetObject(ctx context.Context, prm PrmObjectGet) (ResGetObject, error) {
	p.fillAppropriateKey(&prm.prmCommon)

	var cc callContext
	cc.Context = ctx
	cc.sessionTarget = prm.UseSession

	var res ResGetObject

	err := p.initCallContext(&cc, prm.prmCommon, prmContext{})
	if err != nil {
		return res, err
	}

	return res, p.call(&cc, func() error {
		res, err = cc.client.objectGet(ctx, prm)
		return err
	})
}

// HeadObject reads object header through a remote server using FrostFS API protocol.
//
// Main return value MUST NOT be processed on an erroneous return.
func (p *Pool) HeadObject(ctx context.Context, prm PrmObjectHead) (object.Object, error) {
	p.fillAppropriateKey(&prm.prmCommon)

	var cc callContext

	cc.Context = ctx
	cc.sessionTarget = prm.UseSession

	var obj object.Object

	err := p.initCallContext(&cc, prm.prmCommon, prmContext{})
	if err != nil {
		return obj, err
	}

	return obj, p.call(&cc, func() error {
		obj, err = cc.client.objectHead(ctx, prm)
		return err
	})
}

// ResObjectRange is designed to read payload range of one object
// from FrostFS system.
//
// Must be initialized using Pool.ObjectRange, any other
// usage is unsafe.
type ResObjectRange struct {
	payload             *sdkClient.ObjectRangeReader
	elapsedTimeCallback func(time.Duration)
}

// Read implements io.Reader of the object payload.
func (x *ResObjectRange) Read(p []byte) (int, error) {
	start := time.Now()
	n, err := x.payload.Read(p)
	x.elapsedTimeCallback(time.Since(start))
	return n, err
}

// Close ends reading the payload range and returns the result of the operation
// along with the final results. Must be called after using the ResObjectRange.
func (x *ResObjectRange) Close() error {
	_, err := x.payload.Close()
	return err
}

// ObjectRange initiates reading an object's payload range through a remote
// server using FrostFS API protocol.
//
// Main return value MUST NOT be processed on an erroneous return.
func (p *Pool) ObjectRange(ctx context.Context, prm PrmObjectRange) (ResObjectRange, error) {
	p.fillAppropriateKey(&prm.prmCommon)

	var cc callContext
	cc.Context = ctx
	cc.sessionTarget = prm.UseSession

	var res ResObjectRange

	err := p.initCallContext(&cc, prm.prmCommon, prmContext{})
	if err != nil {
		return res, err
	}

	return res, p.call(&cc, func() error {
		res, err = cc.client.objectRange(ctx, prm)
		return err
	})
}

// ResObjectSearch is designed to read list of object identifiers from FrostFS system.
//
// Must be initialized using Pool.SearchObjects, any other usage is unsafe.
type ResObjectSearch struct {
	r *sdkClient.ObjectListReader
}

// Read reads another list of the object identifiers.
func (x *ResObjectSearch) Read(buf []oid.ID) (int, error) {
	n, ok := x.r.Read(buf)
	if !ok {
		_, err := x.r.Close()
		if err == nil {
			return n, io.EOF
		}

		return n, err
	}

	return n, nil
}

// Iterate iterates over the list of found object identifiers.
// f can return true to stop iteration earlier.
//
// Returns an error if object can't be read.
func (x *ResObjectSearch) Iterate(f func(oid.ID) bool) error {
	return x.r.Iterate(f)
}

// Close ends reading list of the matched objects and returns the result of the operation
// along with the final results. Must be called after using the ResObjectSearch.
func (x *ResObjectSearch) Close() {
	_, _ = x.r.Close()
}

// SearchObjects initiates object selection through a remote server using FrostFS API protocol.
//
// The call only opens the transmission channel, explicit fetching of matched objects
// is done using the ResObjectSearch. Resulting reader must be finally closed.
//
// Main return value MUST NOT be processed on an erroneous return.
func (p *Pool) SearchObjects(ctx context.Context, prm PrmObjectSearch) (ResObjectSearch, error) {
	p.fillAppropriateKey(&prm.prmCommon)

	var cc callContext

	cc.Context = ctx
	cc.sessionTarget = prm.UseSession

	var res ResObjectSearch

	err := p.initCallContext(&cc, prm.prmCommon, prmContext{})
	if err != nil {
		return res, err
	}

	return res, p.call(&cc, func() error {
		res, err = cc.client.objectSearch(ctx, prm)
		return err
	})
}

// PutContainer sends request to save container in FrostFS and waits for the operation to complete.
//
// Waiting parameters can be specified using SetWaitParams. If not called, defaults are used:
//
//	polling interval: 5s
//	waiting timeout: 120s
//
// Success can be verified by reading by identifier (see GetContainer).
//
// Main return value MUST NOT be processed on an erroneous return.
func (p *Pool) PutContainer(ctx context.Context, prm PrmContainerPut) (cid.ID, error) {
	cp, err := p.connection()
	if err != nil {
		return cid.ID{}, err
	}

	return cp.containerPut(ctx, prm)
}

// GetContainer reads FrostFS container by ID.
//
// Main return value MUST NOT be processed on an erroneous return.
func (p *Pool) GetContainer(ctx context.Context, prm PrmContainerGet) (container.Container, error) {
	cp, err := p.connection()
	if err != nil {
		return container.Container{}, err
	}

	return cp.containerGet(ctx, prm)
}

// ListContainers requests identifiers of the account-owned containers.
func (p *Pool) ListContainers(ctx context.Context, prm PrmContainerList) ([]cid.ID, error) {
	cp, err := p.connection()
	if err != nil {
		return nil, err
	}

	return cp.containerList(ctx, prm)
}

// DeleteContainer sends request to remove the FrostFS container and waits for the operation to complete.
//
// Waiting parameters can be specified using SetWaitParams. If not called, defaults are used:
//
//	polling interval: 5s
//	waiting timeout: 120s
//
// Success can be verified by reading by identifier (see GetContainer).
func (p *Pool) DeleteContainer(ctx context.Context, prm PrmContainerDelete) error {
	cp, err := p.connection()
	if err != nil {
		return err
	}

	return cp.containerDelete(ctx, prm)
}

// GetEACL reads eACL table of the FrostFS container.
//
// Main return value MUST NOT be processed on an erroneous return.
func (p *Pool) GetEACL(ctx context.Context, prm PrmContainerEACL) (eacl.Table, error) {
	cp, err := p.connection()
	if err != nil {
		return eacl.Table{}, err
	}

	return cp.containerEACL(ctx, prm)
}

// SetEACL sends request to update eACL table of the FrostFS container and waits for the operation to complete.
//
// Waiting parameters can be specified using SetWaitParams. If not called, defaults are used:
//
//	polling interval: 5s
//	waiting timeout: 120s
//
// Success can be verified by reading by identifier (see GetEACL).
func (p *Pool) SetEACL(ctx context.Context, prm PrmContainerSetEACL) error {
	cp, err := p.connection()
	if err != nil {
		return err
	}

	return cp.containerSetEACL(ctx, prm)
}

// Balance requests current balance of the FrostFS account.
//
// Main return value MUST NOT be processed on an erroneous return.
func (p *Pool) Balance(ctx context.Context, prm PrmBalanceGet) (accounting.Decimal, error) {
	cp, err := p.connection()
	if err != nil {
		return accounting.Decimal{}, err
	}

	return cp.balanceGet(ctx, prm)
}

// Statistic returns connection statistics.
func (p Pool) Statistic() Statistic {
	stat := Statistic{}
	for _, inner := range p.innerPools {
		inner.lock.RLock()
		for _, cl := range inner.clients {
			node := NodeStatistic{
				address:       cl.address(),
				methods:       cl.methodsStatus(),
				overallErrors: cl.overallErrorRate(),
				currentErrors: cl.currentErrorRate(),
			}
			stat.nodes = append(stat.nodes, node)
			stat.overallErrors += node.overallErrors
		}
		inner.lock.RUnlock()
	}

	return stat
}

// waitForContainerPresence waits until the container is found on the FrostFS network.
func waitForContainerPresence(ctx context.Context, cli client, cnrID cid.ID, waitParams *WaitParams) error {
	var prm PrmContainerGet
	prm.SetContainerID(cnrID)

	return waitFor(ctx, waitParams, func(ctx context.Context) bool {
		_, err := cli.containerGet(ctx, prm)
		return err == nil
	})
}

// waitForEACLPresence waits until the container eacl is applied on the FrostFS network.
func waitForEACLPresence(ctx context.Context, cli client, cnrID *cid.ID, table *eacl.Table, waitParams *WaitParams) error {
	var prm PrmContainerEACL
	if cnrID != nil {
		prm.SetContainerID(*cnrID)
	}

	return waitFor(ctx, waitParams, func(ctx context.Context) bool {
		eaclTable, err := cli.containerEACL(ctx, prm)
		if err == nil {
			return eacl.EqualTables(*table, eaclTable)
		}
		return false
	})
}

// waitForContainerRemoved waits until the container is removed from the FrostFS network.
func waitForContainerRemoved(ctx context.Context, cli client, cnrID *cid.ID, waitParams *WaitParams) error {
	var prm PrmContainerGet
	if cnrID != nil {
		prm.SetContainerID(*cnrID)
	}

	return waitFor(ctx, waitParams, func(ctx context.Context) bool {
		_, err := cli.containerGet(ctx, prm)
		return sdkClient.IsErrContainerNotFound(err)
	})
}

// waitFor await that given condition will be met in waitParams time.
func waitFor(ctx context.Context, params *WaitParams, condition func(context.Context) bool) error {
	wctx, cancel := context.WithTimeout(ctx, params.timeout)
	defer cancel()
	ticker := time.NewTimer(params.pollInterval)
	defer ticker.Stop()
	wdone := wctx.Done()
	done := ctx.Done()
	for {
		select {
		case <-done:
			return ctx.Err()
		case <-wdone:
			return wctx.Err()
		case <-ticker.C:
			if condition(ctx) {
				return nil
			}
			ticker.Reset(params.pollInterval)
		}
	}
}

// NetworkInfo requests information about the FrostFS network of which the remote server is a part.
//
// Main return value MUST NOT be processed on an erroneous return.
func (p *Pool) NetworkInfo(ctx context.Context) (netmap.NetworkInfo, error) {
	cp, err := p.connection()
	if err != nil {
		return netmap.NetworkInfo{}, err
	}

	return cp.networkInfo(ctx, prmNetworkInfo{})
}

// Close closes the Pool and releases all the associated resources.
func (p *Pool) Close() {
	p.cancel()
	<-p.closedCh
}

// SyncContainerWithNetwork applies network configuration received via
// the Pool to the container. Changes the container if it does not satisfy
// network configuration.
//
// Pool and container MUST not be nil.
//
// Returns any error that does not allow reading configuration
// from the network.
func SyncContainerWithNetwork(ctx context.Context, cnr *container.Container, p *Pool) error {
	ni, err := p.NetworkInfo(ctx)
	if err != nil {
		return fmt.Errorf("network info: %w", err)
	}

	container.ApplyNetworkConfig(cnr, ni)

	return nil
}

// GetSplitInfo implements relations.Relations.
func (p *Pool) GetSplitInfo(ctx context.Context, cnrID cid.ID, objID oid.ID, tokens relations.Tokens) (*object.SplitInfo, error) {
	var addr oid.Address
	addr.SetContainer(cnrID)
	addr.SetObject(objID)

	var prm PrmObjectHead
	prm.SetAddress(addr)
	if tokens.Bearer != nil {
		prm.UseBearer(*tokens.Bearer)
	}
	if tokens.Session != nil {
		prm.UseSession(*tokens.Session)
	}
	prm.MarkRaw()

	_, err := p.HeadObject(ctx, prm)

	var errSplit *object.SplitInfoError

	switch {
	case errors.As(err, &errSplit):
		return errSplit.SplitInfo(), nil
	case err == nil:
		return nil, relations.ErrNoSplitInfo
	default:
		return nil, fmt.Errorf("failed to get raw object header: %w", err)
	}
}

// ListChildrenByLinker implements relations.Relations.
func (p *Pool) ListChildrenByLinker(ctx context.Context, cnrID cid.ID, objID oid.ID, tokens relations.Tokens) ([]oid.ID, error) {
	var addr oid.Address
	addr.SetContainer(cnrID)
	addr.SetObject(objID)

	var prm PrmObjectHead
	prm.SetAddress(addr)
	if tokens.Bearer != nil {
		prm.UseBearer(*tokens.Bearer)
	}
	if tokens.Session != nil {
		prm.UseSession(*tokens.Session)
	}

	res, err := p.HeadObject(ctx, prm)
	if err != nil {
		return nil, fmt.Errorf("failed to get linking object's header: %w", err)
	}

	return res.Children(), nil
}

// GetLeftSibling implements relations.Relations.
func (p *Pool) GetLeftSibling(ctx context.Context, cnrID cid.ID, objID oid.ID, tokens relations.Tokens) (oid.ID, error) {
	var addr oid.Address
	addr.SetContainer(cnrID)
	addr.SetObject(objID)

	var prm PrmObjectHead
	prm.SetAddress(addr)
	if tokens.Bearer != nil {
		prm.UseBearer(*tokens.Bearer)
	}
	if tokens.Session != nil {
		prm.UseSession(*tokens.Session)
	}

	res, err := p.HeadObject(ctx, prm)
	if err != nil {
		return oid.ID{}, fmt.Errorf("failed to read split chain member's header: %w", err)
	}

	idMember, ok := res.PreviousID()
	if !ok {
		return oid.ID{}, relations.ErrNoLeftSibling
	}
	return idMember, nil
}

// FindSiblingBySplitID implements relations.Relations.
func (p *Pool) FindSiblingBySplitID(ctx context.Context, cnrID cid.ID, splitID *object.SplitID, tokens relations.Tokens) ([]oid.ID, error) {
	var query object.SearchFilters
	query.AddSplitIDFilter(object.MatchStringEqual, splitID)

	var prm PrmObjectSearch
	prm.SetContainerID(cnrID)
	prm.SetFilters(query)
	if tokens.Bearer != nil {
		prm.UseBearer(*tokens.Bearer)
	}
	if tokens.Session != nil {
		prm.UseSession(*tokens.Session)
	}

	res, err := p.SearchObjects(ctx, prm)
	if err != nil {
		return nil, fmt.Errorf("failed to search objects by split ID: %w", err)
	}

	var members []oid.ID
	err = res.Iterate(func(id oid.ID) bool {
		members = append(members, id)
		return false
	})
	if err != nil {
		return nil, fmt.Errorf("failed to iterate found objects: %w", err)
	}

	return members, nil
}

// FindSiblingByParentID implements relations.Relations.
func (p *Pool) FindSiblingByParentID(ctx context.Context, cnrID cid.ID, objID oid.ID, tokens relations.Tokens) ([]oid.ID, error) {
	var query object.SearchFilters
	query.AddParentIDFilter(object.MatchStringEqual, objID)

	var prm PrmObjectSearch
	prm.SetContainerID(cnrID)
	prm.SetFilters(query)
	if tokens.Bearer != nil {
		prm.UseBearer(*tokens.Bearer)
	}
	if tokens.Session != nil {
		prm.UseSession(*tokens.Session)
	}

	resSearch, err := p.SearchObjects(ctx, prm)
	if err != nil {
		return nil, fmt.Errorf("failed to find object children: %w", err)
	}

	var res []oid.ID
	err = resSearch.Iterate(func(id oid.ID) bool {
		res = append(res, id)
		return false
	})
	if err != nil {
		return nil, fmt.Errorf("failed to iterate found objects: %w", err)
	}

	return res, nil
}
