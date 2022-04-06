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
	"strings"
	"sync"
	"time"

	"github.com/nspcc-dev/neo-go/pkg/crypto/keys"
	sessionv2 "github.com/nspcc-dev/neofs-api-go/v2/session"
	"github.com/nspcc-dev/neofs-sdk-go/accounting"
	"github.com/nspcc-dev/neofs-sdk-go/client"
	"github.com/nspcc-dev/neofs-sdk-go/container"
	cid "github.com/nspcc-dev/neofs-sdk-go/container/id"
	"github.com/nspcc-dev/neofs-sdk-go/eacl"
	"github.com/nspcc-dev/neofs-sdk-go/object"
	"github.com/nspcc-dev/neofs-sdk-go/object/address"
	oid "github.com/nspcc-dev/neofs-sdk-go/object/id"
	"github.com/nspcc-dev/neofs-sdk-go/owner"
	"github.com/nspcc-dev/neofs-sdk-go/session"
	"github.com/nspcc-dev/neofs-sdk-go/token"
	"go.uber.org/zap"
)

// Client is a wrapper for client.Client to generate mock.
type Client interface {
	BalanceGet(context.Context, client.PrmBalanceGet) (*client.ResBalanceGet, error)
	ContainerPut(context.Context, client.PrmContainerPut) (*client.ResContainerPut, error)
	ContainerGet(context.Context, client.PrmContainerGet) (*client.ResContainerGet, error)
	ContainerList(context.Context, client.PrmContainerList) (*client.ResContainerList, error)
	ContainerDelete(context.Context, client.PrmContainerDelete) (*client.ResContainerDelete, error)
	ContainerEACL(context.Context, client.PrmContainerEACL) (*client.ResContainerEACL, error)
	ContainerSetEACL(context.Context, client.PrmContainerSetEACL) (*client.ResContainerSetEACL, error)
	EndpointInfo(context.Context, client.PrmEndpointInfo) (*client.ResEndpointInfo, error)
	NetworkInfo(context.Context, client.PrmNetworkInfo) (*client.ResNetworkInfo, error)
	ObjectPutInit(context.Context, client.PrmObjectPutInit) (*client.ObjectWriter, error)
	ObjectDelete(context.Context, client.PrmObjectDelete) (*client.ResObjectDelete, error)
	ObjectGetInit(context.Context, client.PrmObjectGet) (*client.ObjectReader, error)
	ObjectHead(context.Context, client.PrmObjectHead) (*client.ResObjectHead, error)
	ObjectRangeInit(context.Context, client.PrmObjectRange) (*client.ObjectRangeReader, error)
	ObjectSearchInit(context.Context, client.PrmObjectSearch) (*client.ObjectListReader, error)
	SessionCreate(context.Context, client.PrmSessionCreate) (*client.ResSessionCreate, error)
}

// BuilderOptions contains options used to build connection Pool.
type BuilderOptions struct {
	Key                       *ecdsa.PrivateKey
	Logger                    *zap.Logger
	NodeConnectionTimeout     time.Duration
	NodeRequestTimeout        time.Duration
	ClientRebalanceInterval   time.Duration
	SessionTokenThreshold     time.Duration
	SessionExpirationDuration uint64
	nodesParams               []*NodesParam
	clientBuilder             func(endpoint string) (Client, error)
}

type NodesParam struct {
	priority  int
	addresses []string
	weights   []float64
}

type NodeParam struct {
	priority int
	address  string
	weight   float64
}

// Builder is an interim structure used to collect node addresses/weights and
// build connection Pool subsequently.
type Builder struct {
	nodeParams []NodeParam
}

// ContainerPollingParams contains parameters used in polling is a container created or not.
type ContainerPollingParams struct {
	CreationTimeout time.Duration
	PollInterval    time.Duration
}

// DefaultPollingParams creates ContainerPollingParams with default values.
func DefaultPollingParams() *ContainerPollingParams {
	return &ContainerPollingParams{
		CreationTimeout: 120 * time.Second,
		PollInterval:    5 * time.Second,
	}
}

// AddNode adds address/weight pair to node PoolBuilder list.
func (pb *Builder) AddNode(address string, priority int, weight float64) *Builder {
	pb.nodeParams = append(pb.nodeParams, NodeParam{
		address:  address,
		priority: priority,
		weight:   weight,
	})
	return pb
}

// Build creates new Pool based on current PoolBuilder state and options.
func (pb *Builder) Build(ctx context.Context, options *BuilderOptions) (*Pool, error) {
	if len(pb.nodeParams) == 0 {
		return nil, errors.New("no NeoFS peers configured")
	}

	nodesParams := make(map[int]*NodesParam)
	for _, param := range pb.nodeParams {
		nodes, ok := nodesParams[param.priority]
		if !ok {
			nodes = &NodesParam{priority: param.priority}
		}
		nodes.addresses = append(nodes.addresses, param.address)
		nodes.weights = append(nodes.weights, param.weight)
		nodesParams[param.priority] = nodes
	}

	for _, nodes := range nodesParams {
		nodes.weights = adjustWeights(nodes.weights)
		options.nodesParams = append(options.nodesParams, nodes)
	}

	sort.Slice(options.nodesParams, func(i, j int) bool {
		return options.nodesParams[i].priority < options.nodesParams[j].priority
	})

	return newPool(ctx, options)
}

type clientPack struct {
	client  Client
	healthy bool
	address string
}

type CallOption func(config *callConfig)

type callConfig struct {
	useDefaultSession bool
	verb              sessionv2.ObjectSessionVerb
	addr              *address.Address

	key    *ecdsa.PrivateKey
	btoken *token.BearerToken
	stoken *session.Token
}

func WithKey(key *ecdsa.PrivateKey) CallOption {
	return func(config *callConfig) {
		config.key = key
	}
}

func WithBearer(token *token.BearerToken) CallOption {
	return func(config *callConfig) {
		config.btoken = token
	}
}

func WithSession(token *session.Token) CallOption {
	return func(config *callConfig) {
		config.stoken = token
	}
}

func useDefaultSession() CallOption {
	return func(config *callConfig) {
		config.useDefaultSession = true
	}
}

func useAddress(addr *address.Address) CallOption {
	return func(config *callConfig) {
		config.addr = addr
	}
}

func useVerb(verb sessionv2.ObjectSessionVerb) CallOption {
	return func(config *callConfig) {
		config.verb = verb
	}
}

func cfgFromOpts(opts ...CallOption) *callConfig {
	var cfg = new(callConfig)
	for _, opt := range opts {
		opt(cfg)
	}
	return cfg
}

type Pool struct {
	innerPools      []*innerPool
	key             *ecdsa.PrivateKey
	owner           *owner.ID
	cancel          context.CancelFunc
	closedCh        chan struct{}
	cache           *SessionCache
	stokenDuration  uint64
	stokenThreshold time.Duration
}

type innerPool struct {
	lock        sync.RWMutex
	sampler     *Sampler
	clientPacks []*clientPack
}

const (
	defaultSessionTokenExpirationDuration = 100 // in blocks

	defaultSessionTokenThreshold = 5 * time.Second
)

func newPool(ctx context.Context, options *BuilderOptions) (*Pool, error) {
	cache, err := NewCache()
	if err != nil {
		return nil, fmt.Errorf("couldn't create cache: %w", err)
	}

	if options.SessionExpirationDuration == 0 {
		options.SessionExpirationDuration = defaultSessionTokenExpirationDuration
	}

	if options.SessionTokenThreshold <= 0 {
		options.SessionTokenThreshold = defaultSessionTokenThreshold
	}

	ownerID := owner.NewIDFromPublicKey(&options.Key.PublicKey)

	inner := make([]*innerPool, len(options.nodesParams))
	var atLeastOneHealthy bool

	if options.clientBuilder == nil {
		options.clientBuilder = func(addr string) (Client, error) {
			var c client.Client

			var prmInit client.PrmInit
			prmInit.ResolveNeoFSFailures()
			prmInit.SetDefaultPrivateKey(*options.Key)

			c.Init(prmInit)

			var prmDial client.PrmDial
			prmDial.SetServerURI(addr)
			prmDial.SetTimeout(options.NodeConnectionTimeout)

			return &c, c.Dial(prmDial)
		}
	}

	for i, params := range options.nodesParams {
		clientPacks := make([]*clientPack, len(params.weights))
		for j, addr := range params.addresses {
			c, err := options.clientBuilder(addr)
			if err != nil {
				return nil, err
			}
			var healthy bool
			cliRes, err := createSessionTokenForDuration(ctx, c, options.SessionExpirationDuration)
			if err != nil && options.Logger != nil {
				options.Logger.Warn("failed to create neofs session token for client",
					zap.String("address", addr),
					zap.Error(err))
			} else if err == nil {
				healthy, atLeastOneHealthy = true, true
				st := sessionTokenForOwner(ownerID, cliRes)
				_ = cache.Put(formCacheKey(addr, options.Key), st)
			}
			clientPacks[j] = &clientPack{client: c, healthy: healthy, address: addr}
		}
		source := rand.NewSource(time.Now().UnixNano())
		sampler := NewSampler(params.weights, source)

		inner[i] = &innerPool{
			sampler:     sampler,
			clientPacks: clientPacks,
		}
	}

	if !atLeastOneHealthy {
		return nil, fmt.Errorf("at least one node must be healthy")
	}

	ctx, cancel := context.WithCancel(ctx)
	pool := &Pool{
		innerPools:      inner,
		key:             options.Key,
		owner:           ownerID,
		cancel:          cancel,
		closedCh:        make(chan struct{}),
		cache:           cache,
		stokenDuration:  options.SessionExpirationDuration,
		stokenThreshold: options.SessionTokenThreshold,
	}
	go startRebalance(ctx, pool, options)
	return pool, nil
}

func startRebalance(ctx context.Context, p *Pool, options *BuilderOptions) {
	ticker := time.NewTimer(options.ClientRebalanceInterval)
	buffers := make([][]float64, len(options.nodesParams))
	for i, params := range options.nodesParams {
		buffers[i] = make([]float64, len(params.weights))
	}

	for {
		select {
		case <-ctx.Done():
			close(p.closedCh)
			return
		case <-ticker.C:
			updateNodesHealth(ctx, p, options, buffers)
			ticker.Reset(options.ClientRebalanceInterval)
		}
	}
}

func updateNodesHealth(ctx context.Context, p *Pool, options *BuilderOptions, buffers [][]float64) {
	wg := sync.WaitGroup{}
	for i, inner := range p.innerPools {
		wg.Add(1)

		bufferWeights := buffers[i]
		go func(i int, innerPool *innerPool) {
			defer wg.Done()
			updateInnerNodesHealth(ctx, p, i, options, bufferWeights)
		}(i, inner)
	}
	wg.Wait()
}

func updateInnerNodesHealth(ctx context.Context, pool *Pool, i int, options *BuilderOptions, bufferWeights []float64) {
	if i > len(pool.innerPools)-1 {
		return
	}
	p := pool.innerPools[i]

	healthyChanged := false
	wg := sync.WaitGroup{}

	var prmEndpoint client.PrmEndpointInfo

	for j, cPack := range p.clientPacks {
		wg.Add(1)
		go func(j int, cli Client) {
			defer wg.Done()
			ok := true
			tctx, c := context.WithTimeout(ctx, options.NodeRequestTimeout)
			defer c()

			if _, err := cli.EndpointInfo(tctx, prmEndpoint); err != nil {
				ok = false
				bufferWeights[j] = 0
			}
			p.lock.RLock()
			cp := *p.clientPacks[j]
			p.lock.RUnlock()

			if ok {
				bufferWeights[j] = options.nodesParams[i].weights[j]
				if !cp.healthy {
					cliRes, err := createSessionTokenForDuration(ctx, cli, options.SessionExpirationDuration)
					if err != nil {
						ok = false
						bufferWeights[j] = 0
					} else {
						tkn := pool.newSessionToken(cliRes)

						_ = pool.cache.Put(formCacheKey(cp.address, pool.key), tkn)
					}
				}
			} else {
				pool.cache.DeleteByPrefix(cp.address)
			}

			p.lock.Lock()
			if p.clientPacks[j].healthy != ok {
				p.clientPacks[j].healthy = ok
				healthyChanged = true
			}
			p.lock.Unlock()
		}(j, cPack.client)
	}
	wg.Wait()

	if healthyChanged {
		probabilities := adjustWeights(bufferWeights)
		source := rand.NewSource(time.Now().UnixNano())
		p.lock.Lock()
		p.sampler = NewSampler(probabilities, source)
		p.lock.Unlock()
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

func (p *Pool) Connection() (Client, *session.Token, error) {
	cp, err := p.connection()
	if err != nil {
		return nil, nil, err
	}

	tok := p.cache.Get(formCacheKey(cp.address, p.key))
	return cp.client, tok, nil
}

func (p *Pool) connection() (*clientPack, error) {
	for _, inner := range p.innerPools {
		cp, err := inner.connection()
		if err == nil {
			return cp, nil
		}
	}

	return nil, errors.New("no healthy client")
}

func (p *innerPool) connection() (*clientPack, error) {
	p.lock.RLock()
	defer p.lock.RUnlock()
	if len(p.clientPacks) == 1 {
		cp := p.clientPacks[0]
		if cp.healthy {
			return cp, nil
		}
		return nil, errors.New("no healthy client")
	}
	attempts := 3 * len(p.clientPacks)
	for k := 0; k < attempts; k++ {
		i := p.sampler.Next()
		if cp := p.clientPacks[i]; cp.healthy {
			return cp, nil
		}
	}

	return nil, errors.New("no healthy client")
}

func (p *Pool) OwnerID() *owner.ID {
	return p.owner
}

func formCacheKey(address string, key *ecdsa.PrivateKey) string {
	k := keys.PrivateKey{PrivateKey: *key}
	return address + k.String()
}

func (p *Pool) conn(ctx context.Context, cfg *callConfig) (*clientPack, error) {
	cp, err := p.connection()
	if err != nil {
		return nil, err
	}

	key := p.key
	if cfg.key != nil {
		key = cfg.key
	}

	sessionToken := cfg.stoken
	if sessionToken == nil && cfg.useDefaultSession {
		cacheKey := formCacheKey(cp.address, key)
		sessionToken = p.cache.Get(cacheKey)
		if sessionToken == nil {
			cliRes, err := createSessionTokenForDuration(ctx, cp.client, p.stokenDuration)
			if err != nil {
				return nil, err
			}

			ownerID := owner.NewIDFromPublicKey(&key.PublicKey)
			sessionToken = sessionTokenForOwner(ownerID, cliRes)

			cfg.stoken = sessionToken

			_ = p.cache.Put(cacheKey, sessionToken)
		}
	}

	return cp, nil
}

func (p *Pool) checkSessionTokenErr(err error, address string) bool {
	if err == nil {
		return false
	}

	if strings.Contains(err.Error(), "session token does not exist") ||
		strings.Contains(err.Error(), "session token has been expired") {
		p.cache.DeleteByPrefix(address)
		return true
	}

	return false
}

func createSessionTokenForDuration(ctx context.Context, c Client, dur uint64) (*client.ResSessionCreate, error) {
	ni, err := c.NetworkInfo(ctx, client.PrmNetworkInfo{})
	if err != nil {
		return nil, err
	}

	epoch := ni.Info().CurrentEpoch()

	var prm client.PrmSessionCreate
	if math.MaxUint64-epoch < dur {
		prm.SetExp(math.MaxUint64)
	} else {
		prm.SetExp(epoch + dur)
	}

	return c.SessionCreate(ctx, prm)
}

func (p *Pool) removeSessionTokenAfterThreshold(cfg *callConfig) error {
	cp, err := p.connection()
	if err != nil {
		return err
	}

	key := p.key
	if cfg.key != nil {
		key = cfg.key
	}

	ts, ok := p.cache.GetAccessTime(formCacheKey(cp.address, key))
	if ok && time.Since(ts) > p.stokenThreshold {
		p.cache.DeleteByPrefix(cp.address)
	}

	return nil
}

type callContext struct {
	// base context for RPC
	context.Context

	client Client

	// client endpoint
	endpoint string

	// request signer
	key *ecdsa.PrivateKey

	// flag to open default session if session token is missing
	sessionDefault bool
	sessionTarget  func(session.Token)
	sessionContext *session.ObjectContext
}

func (p *Pool) initCallContext(ctx *callContext, cfg *callConfig) error {
	cp, err := p.connection()
	if err != nil {
		return err
	}

	ctx.key = cfg.key
	if ctx.key == nil {
		// use pool key if caller didn't specify its own
		ctx.key = p.key
	}

	ctx.endpoint = cp.address
	ctx.client = cp.client

	if ctx.sessionTarget != nil && cfg.stoken != nil {
		ctx.sessionTarget(*cfg.stoken)
	}

	// note that we don't override session provided by the caller
	ctx.sessionDefault = cfg.stoken == nil && cfg.useDefaultSession
	if ctx.sessionDefault {
		ctx.sessionContext = session.NewObjectContext()
		ctx.sessionContext.ToV2().SetVerb(cfg.verb)
		ctx.sessionContext.ApplyTo(cfg.addr)
	}

	return err
}

type callContextWithRetry struct {
	callContext

	noRetry bool
}

func (p *Pool) initCallContextWithRetry(ctx *callContextWithRetry, cfg *callConfig) error {
	err := p.initCallContext(&ctx.callContext, cfg)
	if err != nil {
		return err
	}

	// don't retry if session was specified by the caller
	ctx.noRetry = cfg.stoken != nil

	return nil
}

// opens new session or uses cached one.
// Must be called only on initialized callContext with set sessionTarget.
func (p *Pool) openDefaultSession(ctx *callContext) error {
	cacheKey := formCacheKey(ctx.endpoint, ctx.key)

	tok := p.cache.Get(cacheKey)
	if tok == nil {
		// open new session
		cliRes, err := createSessionTokenForDuration(ctx, ctx.client, p.stokenDuration)
		if err != nil {
			return fmt.Errorf("session API client: %w", err)
		}

		tok = sessionTokenForOwner(owner.NewIDFromPublicKey(&ctx.key.PublicKey), cliRes)
		// cache the opened session
		p.cache.Put(cacheKey, tok)
	}

	tokToSign := *tok
	tokToSign.SetContext(ctx.sessionContext)

	// sign the token
	if err := tokToSign.Sign(ctx.key); err != nil {
		return fmt.Errorf("sign token of the opened session: %w", err)
	}

	ctx.sessionTarget(tokToSign)

	return nil
}

// opens default session (if sessionDefault is set), and calls f. If f returns
// session-related error (*), and retrying is enabled, then f is called once more.
//
// (*) in this case cached token is removed.
func (p *Pool) callWithRetry(ctx *callContextWithRetry, f func() error) error {
	var err error

	if ctx.sessionDefault {
		err = p.openDefaultSession(&ctx.callContext)
		if err != nil {
			return fmt.Errorf("open default session: %w", err)
		}
	}

	err = f()

	if p.checkSessionTokenErr(err, ctx.endpoint) && !ctx.noRetry {
		// don't retry anymore
		ctx.noRetry = true
		return p.callWithRetry(ctx, f)
	}

	return err
}

func (p *Pool) PutObject(ctx context.Context, hdr object.Object, payload io.Reader, opts ...CallOption) (*oid.ID, error) {
	cfg := cfgFromOpts(append(opts,
		useDefaultSession(),
		useVerb(sessionv2.ObjectVerbPut),
		useAddress(newAddressFromCnrID(hdr.ContainerID())))...)

	// Put object is different from other object service methods. Put request
	// can't be resent in case of session token failures (i.e. session token is
	// invalid due to lifetime expiration or server restart). The reason is that
	// object's payload can be provided as a stream that should be read only once.
	//
	// To solve this issue, pool regenerates session tokens upon each request.
	// In case of subsequent requests, pool avoids session token initialization
	// by checking when the session token was accessed for the last time. If it
	// hits a threshold, session token is removed from cache for a new one to be
	// issued.
	err := p.removeSessionTokenAfterThreshold(cfg)
	if err != nil {
		return nil, err
	}

	var ctxCall callContext

	ctxCall.Context = ctx

	err = p.initCallContext(&ctxCall, cfg)
	if err != nil {
		return nil, fmt.Errorf("init call context")
	}

	var prm client.PrmObjectPutInit

	wObj, err := ctxCall.client.ObjectPutInit(ctx, prm)
	if err != nil {
		return nil, fmt.Errorf("init writing on API client: %w", err)
	}

	if ctxCall.sessionDefault {
		ctxCall.sessionTarget = wObj.WithinSession
		err = p.openDefaultSession(&ctxCall)
		if err != nil {
			return nil, fmt.Errorf("open default session: %w", err)
		}
	}

	wObj.UseKey(*ctxCall.key)

	if cfg.btoken != nil {
		wObj.WithBearerToken(*cfg.btoken)
	}

	if wObj.WriteHeader(hdr) {
		sz := hdr.PayloadSize()

		if data := hdr.Payload(); len(data) > 0 {
			if payload != nil {
				payload = io.MultiReader(bytes.NewReader(data), payload)
			} else {
				payload = bytes.NewReader(data)
				sz = uint64(len(data))
			}
		}

		if payload != nil {
			const defaultBufferSizePut = 3 << 20 // configure?

			if sz == 0 || sz > defaultBufferSizePut {
				sz = defaultBufferSizePut
			}

			buf := make([]byte, sz)

			var n int

			for {
				n, err = payload.Read(buf)
				if n > 0 {
					if !wObj.WritePayloadChunk(buf[:n]) {
						break
					}

					continue
				}

				if errors.Is(err, io.EOF) {
					break
				}

				return nil, fmt.Errorf("read payload: %w", err)
			}
		}
	}

	res, err := wObj.Close()
	if err != nil { // here err already carries both status and client errors
		// removes session token from cache in case of token error
		p.checkSessionTokenErr(err, ctxCall.endpoint)
		return nil, fmt.Errorf("client failure: %w", err)
	}

	var id oid.ID

	if !res.ReadStoredObjectID(&id) {
		return nil, errors.New("missing ID of the stored object")
	}

	return &id, nil
}

func (p *Pool) DeleteObject(ctx context.Context, addr address.Address, opts ...CallOption) error {
	cfg := cfgFromOpts(append(opts,
		useDefaultSession(),
		useVerb(sessionv2.ObjectVerbDelete),
		useAddress(&addr))...)

	var prm client.PrmObjectDelete

	var cc callContextWithRetry

	cc.Context = ctx
	cc.sessionTarget = prm.WithinSession

	err := p.initCallContextWithRetry(&cc, cfg)
	if err != nil {
		return err
	}

	if cnr := addr.ContainerID(); cnr != nil {
		prm.FromContainer(*cnr)
	}

	if obj := addr.ObjectID(); obj != nil {
		prm.ByID(*obj)
	}

	if cfg.btoken != nil {
		prm.WithBearerToken(*cfg.btoken)
	}

	prm.UseKey(*cc.key)

	return p.callWithRetry(&cc, func() error {
		_, err := cc.client.ObjectDelete(ctx, prm)
		if err != nil {
			return fmt.Errorf("remove object via client: %w", err)
		}

		return nil
	})
}

type objectReadCloser client.ObjectReader

func (x *objectReadCloser) Read(p []byte) (int, error) {
	return (*client.ObjectReader)(x).Read(p)
}

func (x *objectReadCloser) Close() error {
	_, err := (*client.ObjectReader)(x).Close()
	return err
}

type ResGetObject struct {
	Header object.Object

	Payload io.ReadCloser
}

func (p *Pool) GetObject(ctx context.Context, addr address.Address, opts ...CallOption) (*ResGetObject, error) {
	cfg := cfgFromOpts(append(opts,
		useDefaultSession(),
		useVerb(sessionv2.ObjectVerbGet),
		useAddress(&addr))...)

	var prm client.PrmObjectGet

	var cc callContextWithRetry

	cc.Context = ctx
	cc.sessionTarget = prm.WithinSession

	err := p.initCallContextWithRetry(&cc, cfg)
	if err != nil {
		return nil, err
	}

	if cnr := addr.ContainerID(); cnr != nil {
		prm.FromContainer(*cnr)
	}

	if obj := addr.ObjectID(); obj != nil {
		prm.ByID(*obj)
	}

	if cfg.btoken != nil {
		prm.WithBearerToken(*cfg.btoken)
	}

	var res ResGetObject

	err = p.callWithRetry(&cc, func() error {
		rObj, err := cc.client.ObjectGetInit(ctx, prm)
		if err != nil {
			return fmt.Errorf("init object reading on client: %w", err)
		}

		rObj.UseKey(*cc.key)

		if !rObj.ReadHeader(&res.Header) {
			_, err = rObj.Close()
			return fmt.Errorf("read header: %w", err)
		}

		res.Payload = (*objectReadCloser)(rObj)

		return nil
	})
	if err != nil {
		return nil, err
	}

	return &res, nil
}

func (p *Pool) HeadObject(ctx context.Context, addr address.Address, opts ...CallOption) (*object.Object, error) {
	cfg := cfgFromOpts(append(opts,
		useDefaultSession(),
		useVerb(sessionv2.ObjectVerbHead),
		useAddress(&addr))...)

	var prm client.PrmObjectHead

	var cc callContextWithRetry

	cc.Context = ctx
	cc.sessionTarget = prm.WithinSession

	err := p.initCallContextWithRetry(&cc, cfg)
	if err != nil {
		return nil, err
	}

	if cnr := addr.ContainerID(); cnr != nil {
		prm.FromContainer(*cnr)
	}

	if obj := addr.ObjectID(); obj != nil {
		prm.ByID(*obj)
	}

	if cfg.btoken != nil {
		prm.WithBearerToken(*cfg.btoken)
	}

	prm.UseKey(*cc.key)

	var obj object.Object

	err = p.callWithRetry(&cc, func() error {
		res, err := cc.client.ObjectHead(ctx, prm)
		if err != nil {
			return fmt.Errorf("read object header via client: %w", err)
		}

		if !res.ReadHeader(&obj) {
			return errors.New("missing object header in response")
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return &obj, nil
}

type ResObjectRange struct {
	payload *client.ObjectRangeReader
}

func (x *ResObjectRange) Read(p []byte) (int, error) {
	return x.payload.Read(p)
}

func (x *ResObjectRange) Close() error {
	_, err := x.payload.Close()
	return err
}

func (p *Pool) ObjectRange(ctx context.Context, addr address.Address, off, ln uint64, opts ...CallOption) (*ResObjectRange, error) {
	cfg := cfgFromOpts(append(opts,
		useDefaultSession(),
		useVerb(sessionv2.ObjectVerbRange),
		useAddress(&addr))...)

	var prm client.PrmObjectRange

	prm.SetOffset(off)
	prm.SetLength(ln)

	var cc callContextWithRetry

	cc.Context = ctx
	cc.sessionTarget = prm.WithinSession

	err := p.initCallContextWithRetry(&cc, cfg)
	if err != nil {
		return nil, err
	}

	if cnr := addr.ContainerID(); cnr != nil {
		prm.FromContainer(*cnr)
	}

	if obj := addr.ObjectID(); obj != nil {
		prm.ByID(*obj)
	}

	if cfg.btoken != nil {
		prm.WithBearerToken(*cfg.btoken)
	}

	var res ResObjectRange

	err = p.callWithRetry(&cc, func() error {
		var err error

		res.payload, err = cc.client.ObjectRangeInit(ctx, prm)
		if err != nil {
			return fmt.Errorf("init payload range reading on client: %w", err)
		}

		res.payload.UseKey(*cc.key)

		return nil
	})
	if err != nil {
		return nil, err
	}

	return &res, nil
}

type ResObjectSearch struct {
	r *client.ObjectListReader
}

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

func (x *ResObjectSearch) Iterate(f func(oid.ID) bool) error {
	return x.r.Iterate(f)
}

func (x *ResObjectSearch) Close() {
	_, _ = x.r.Close()
}

func (p *Pool) SearchObjects(ctx context.Context, idCnr cid.ID, filters object.SearchFilters, opts ...CallOption) (*ResObjectSearch, error) {
	cfg := cfgFromOpts(append(opts,
		useDefaultSession(),
		useVerb(sessionv2.ObjectVerbSearch),
		useAddress(newAddressFromCnrID(&idCnr)))...)

	var prm client.PrmObjectSearch

	prm.InContainer(idCnr)
	prm.SetFilters(filters)

	if cfg.btoken != nil {
		prm.WithBearerToken(*cfg.btoken)
	}

	var cc callContextWithRetry

	cc.Context = ctx
	cc.sessionTarget = prm.WithinSession

	err := p.initCallContextWithRetry(&cc, cfg)
	if err != nil {
		return nil, err
	}

	var res ResObjectSearch

	err = p.callWithRetry(&cc, func() error {
		var err error

		res.r, err = cc.client.ObjectSearchInit(ctx, prm)
		if err != nil {
			return fmt.Errorf("init object searching on client: %w", err)
		}

		res.r.UseKey(*cc.key)

		return nil
	})
	if err != nil {
		return nil, err
	}

	return &res, nil
}

func (p *Pool) PutContainer(ctx context.Context, cnr *container.Container, opts ...CallOption) (*cid.ID, error) {
	cfg := cfgFromOpts(opts...)
	cp, err := p.conn(ctx, cfg)
	if err != nil {
		return nil, err
	}

	var cliPrm client.PrmContainerPut

	if cnr != nil {
		cliPrm.SetContainer(*cnr)
	}

	res, err := cp.client.ContainerPut(ctx, cliPrm)
	if err != nil { // here err already carries both status and client errors
		return nil, err
	}

	return res.ID(), nil
}

func (p *Pool) GetContainer(ctx context.Context, cid *cid.ID, opts ...CallOption) (*container.Container, error) {
	cfg := cfgFromOpts(opts...)
	cp, err := p.conn(ctx, cfg)
	if err != nil {
		return nil, err
	}

	var cliPrm client.PrmContainerGet

	if cid != nil {
		cliPrm.SetContainer(*cid)
	}

	res, err := cp.client.ContainerGet(ctx, cliPrm)
	if err != nil { // here err already carries both status and client errors
		return nil, err
	}

	return res.Container(), nil
}

func (p *Pool) ListContainers(ctx context.Context, ownerID *owner.ID, opts ...CallOption) ([]cid.ID, error) {
	cfg := cfgFromOpts(opts...)
	cp, err := p.conn(ctx, cfg)
	if err != nil {
		return nil, err
	}

	var cliPrm client.PrmContainerList

	if ownerID != nil {
		cliPrm.SetAccount(*ownerID)
	}

	res, err := cp.client.ContainerList(ctx, cliPrm)
	if err != nil { // here err already carries both status and client errors
		return nil, err
	}

	return res.Containers(), nil
}

func (p *Pool) DeleteContainer(ctx context.Context, cid *cid.ID, opts ...CallOption) error {
	cfg := cfgFromOpts(opts...)
	cp, err := p.conn(ctx, cfg)
	if err != nil {
		return err
	}

	var cliPrm client.PrmContainerDelete

	if cid != nil {
		cliPrm.SetContainer(*cid)
	}

	if cfg.stoken != nil {
		cliPrm.SetSessionToken(*cfg.stoken)
	}

	_, err = cp.client.ContainerDelete(ctx, cliPrm)

	// here err already carries both status and client errors

	return err
}

func (p *Pool) GetEACL(ctx context.Context, cid *cid.ID, opts ...CallOption) (*eacl.Table, error) {
	cfg := cfgFromOpts(opts...)
	cp, err := p.conn(ctx, cfg)
	if err != nil {
		return nil, err
	}

	var cliPrm client.PrmContainerEACL

	if cid != nil {
		cliPrm.SetContainer(*cid)
	}

	res, err := cp.client.ContainerEACL(ctx, cliPrm)
	if err != nil { // here err already carries both status and client errors
		return nil, err
	}

	return res.Table(), nil
}

func (p *Pool) SetEACL(ctx context.Context, table *eacl.Table, opts ...CallOption) error {
	cfg := cfgFromOpts(opts...)
	cp, err := p.conn(ctx, cfg)
	if err != nil {
		return err
	}

	var cliPrm client.PrmContainerSetEACL

	if table != nil {
		cliPrm.SetTable(*table)
	}

	_, err = cp.client.ContainerSetEACL(ctx, cliPrm)

	// here err already carries both status and client errors

	return err
}

func (p *Pool) Balance(ctx context.Context, o *owner.ID, opts ...CallOption) (*accounting.Decimal, error) {
	cfg := cfgFromOpts(opts...)
	cp, err := p.conn(ctx, cfg)
	if err != nil {
		return nil, err
	}

	var cliPrm client.PrmBalanceGet

	if o != nil {
		cliPrm.SetAccount(*o)
	}

	res, err := cp.client.BalanceGet(ctx, cliPrm)
	if err != nil { // here err already carries both status and client errors
		return nil, err
	}

	return res.Amount(), nil
}

func (p *Pool) WaitForContainerPresence(ctx context.Context, cid *cid.ID, pollParams *ContainerPollingParams) error {
	conn, _, err := p.Connection()
	if err != nil {
		return err
	}
	wctx, cancel := context.WithTimeout(ctx, pollParams.CreationTimeout)
	defer cancel()
	ticker := time.NewTimer(pollParams.PollInterval)
	defer ticker.Stop()
	wdone := wctx.Done()
	done := ctx.Done()

	var cliPrm client.PrmContainerGet

	if cid != nil {
		cliPrm.SetContainer(*cid)
	}

	for {
		select {
		case <-done:
			return ctx.Err()
		case <-wdone:
			return wctx.Err()
		case <-ticker.C:
			_, err = conn.ContainerGet(ctx, cliPrm)
			if err == nil {
				return nil
			}
			ticker.Reset(pollParams.PollInterval)
		}
	}
}

// Close closes the Pool and releases all the associated resources.
func (p *Pool) Close() {
	p.cancel()
	<-p.closedCh
}

// creates new session token from SessionCreate call result.
func (p *Pool) newSessionToken(cliRes *client.ResSessionCreate) *session.Token {
	return sessionTokenForOwner(p.owner, cliRes)
}

// creates new session token with specified owner from SessionCreate call result.
func sessionTokenForOwner(id *owner.ID, cliRes *client.ResSessionCreate) *session.Token {
	st := session.NewToken()
	st.SetOwnerID(id)
	st.SetID(cliRes.ID())
	st.SetSessionKey(cliRes.PublicKey())

	return st
}

func newAddressFromCnrID(cnrID *cid.ID) *address.Address {
	addr := address.NewAddress()
	addr.SetContainerID(cnrID)
	return addr
}

func copySessionTokenWithoutSignatureAndContext(from session.Token) (to session.Token) {
	to.SetIat(from.Iat())
	to.SetExp(from.Exp())
	to.SetNbf(from.Nbf())

	sessionTokenID := make([]byte, len(from.ID()))
	copy(sessionTokenID, from.ID())
	to.SetID(sessionTokenID)

	sessionTokenKey := make([]byte, len(from.SessionKey()))
	copy(sessionTokenKey, from.SessionKey())
	to.SetSessionKey(sessionTokenKey)

	var sessionTokenOwner owner.ID
	buf, err := from.OwnerID().Marshal()
	if err != nil {
		panic(err) // should never happen
	}
	err = sessionTokenOwner.Unmarshal(buf)
	if err != nil {
		panic(err) // should never happen
	}
	to.SetOwnerID(&sessionTokenOwner)

	return to
}
