package pool

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"strconv"
	"testing"
	"time"

	apistatus "github.com/TrueCloudLab/frostfs-sdk-go/client/status"
	cid "github.com/TrueCloudLab/frostfs-sdk-go/container/id"
	frostfsecdsa "github.com/TrueCloudLab/frostfs-sdk-go/crypto/ecdsa"
	"github.com/TrueCloudLab/frostfs-sdk-go/object"
	oid "github.com/TrueCloudLab/frostfs-sdk-go/object/id"
	"github.com/TrueCloudLab/frostfs-sdk-go/session"
	"github.com/TrueCloudLab/frostfs-sdk-go/user"
	"github.com/nspcc-dev/neo-go/pkg/crypto/keys"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestBuildPoolClientFailed(t *testing.T) {
	mockClientBuilder := func(addr string) client {
		mockCli := newMockClient(addr, *newPrivateKey(t))
		mockCli.errOnDial()
		return mockCli
	}

	opts := InitParameters{
		key:        newPrivateKey(t),
		nodeParams: []NodeParam{{1, "peer0", 1}},
	}
	opts.setClientBuilder(mockClientBuilder)

	pool, err := NewPool(opts)
	require.NoError(t, err)
	err = pool.Dial(context.Background())
	require.Error(t, err)
}

func TestBuildPoolCreateSessionFailed(t *testing.T) {
	clientMockBuilder := func(addr string) client {
		mockCli := newMockClient(addr, *newPrivateKey(t))
		mockCli.errOnCreateSession()
		return mockCli
	}

	opts := InitParameters{
		key:        newPrivateKey(t),
		nodeParams: []NodeParam{{1, "peer0", 1}},
	}
	opts.setClientBuilder(clientMockBuilder)

	pool, err := NewPool(opts)
	require.NoError(t, err)
	err = pool.Dial(context.Background())
	require.Error(t, err)
}

func newPrivateKey(t *testing.T) *ecdsa.PrivateKey {
	p, err := keys.NewPrivateKey()
	require.NoError(t, err)
	return &p.PrivateKey
}

func TestBuildPoolOneNodeFailed(t *testing.T) {
	nodes := []NodeParam{
		{1, "peer0", 1},
		{2, "peer1", 1},
	}

	var clientKeys []*ecdsa.PrivateKey
	mockClientBuilder := func(addr string) client {
		key := newPrivateKey(t)
		clientKeys = append(clientKeys, key)

		if addr == nodes[0].address {
			mockCli := newMockClient(addr, *key)
			mockCli.errOnEndpointInfo()
			return mockCli
		}

		return newMockClient(addr, *key)
	}

	log, err := zap.NewProduction()
	require.NoError(t, err)
	opts := InitParameters{
		key:                     newPrivateKey(t),
		clientRebalanceInterval: 1000 * time.Millisecond,
		logger:                  log,
		nodeParams:              nodes,
	}
	opts.setClientBuilder(mockClientBuilder)

	clientPool, err := NewPool(opts)
	require.NoError(t, err)
	err = clientPool.Dial(context.Background())
	require.NoError(t, err)
	t.Cleanup(clientPool.Close)

	expectedAuthKey := frostfsecdsa.PublicKey(clientKeys[1].PublicKey)
	condition := func() bool {
		cp, err := clientPool.connection()
		if err != nil {
			return false
		}
		st, _ := clientPool.cache.Get(formCacheKey(cp.address(), clientPool.key))
		return st.AssertAuthKey(&expectedAuthKey)
	}
	require.Never(t, condition, 900*time.Millisecond, 100*time.Millisecond)
	require.Eventually(t, condition, 3*time.Second, 300*time.Millisecond)
}

func TestBuildPoolZeroNodes(t *testing.T) {
	opts := InitParameters{
		key: newPrivateKey(t),
	}
	_, err := NewPool(opts)
	require.Error(t, err)
}

func TestOneNode(t *testing.T) {
	key1 := newPrivateKey(t)
	mockClientBuilder := func(addr string) client {
		return newMockClient(addr, *key1)
	}

	opts := InitParameters{
		key:        newPrivateKey(t),
		nodeParams: []NodeParam{{1, "peer0", 1}},
	}
	opts.setClientBuilder(mockClientBuilder)

	pool, err := NewPool(opts)
	require.NoError(t, err)
	err = pool.Dial(context.Background())
	require.NoError(t, err)
	t.Cleanup(pool.Close)

	cp, err := pool.connection()
	require.NoError(t, err)
	st, _ := pool.cache.Get(formCacheKey(cp.address(), pool.key))
	expectedAuthKey := frostfsecdsa.PublicKey(key1.PublicKey)
	require.True(t, st.AssertAuthKey(&expectedAuthKey))
}

func TestTwoNodes(t *testing.T) {
	var clientKeys []*ecdsa.PrivateKey
	mockClientBuilder := func(addr string) client {
		key := newPrivateKey(t)
		clientKeys = append(clientKeys, key)
		return newMockClient(addr, *key)
	}

	opts := InitParameters{
		key: newPrivateKey(t),
		nodeParams: []NodeParam{
			{1, "peer0", 1},
			{1, "peer1", 1},
		},
	}
	opts.setClientBuilder(mockClientBuilder)

	pool, err := NewPool(opts)
	require.NoError(t, err)
	err = pool.Dial(context.Background())
	require.NoError(t, err)
	t.Cleanup(pool.Close)

	cp, err := pool.connection()
	require.NoError(t, err)
	st, _ := pool.cache.Get(formCacheKey(cp.address(), pool.key))
	require.True(t, assertAuthKeyForAny(st, clientKeys))
}

func assertAuthKeyForAny(st session.Object, clientKeys []*ecdsa.PrivateKey) bool {
	for _, key := range clientKeys {
		expectedAuthKey := frostfsecdsa.PublicKey(key.PublicKey)
		if st.AssertAuthKey(&expectedAuthKey) {
			return true
		}
	}
	return false
}

func TestOneOfTwoFailed(t *testing.T) {
	nodes := []NodeParam{
		{1, "peer0", 1},
		{9, "peer1", 1},
	}

	var clientKeys []*ecdsa.PrivateKey
	mockClientBuilder := func(addr string) client {
		key := newPrivateKey(t)
		clientKeys = append(clientKeys, key)

		if addr == nodes[0].address {
			return newMockClient(addr, *key)
		}

		mockCli := newMockClient(addr, *key)
		mockCli.errOnEndpointInfo()
		mockCli.errOnNetworkInfo()
		return mockCli
	}

	opts := InitParameters{
		key:                     newPrivateKey(t),
		nodeParams:              nodes,
		clientRebalanceInterval: 200 * time.Millisecond,
	}
	opts.setClientBuilder(mockClientBuilder)

	pool, err := NewPool(opts)
	require.NoError(t, err)
	err = pool.Dial(context.Background())
	require.NoError(t, err)

	require.NoError(t, err)
	t.Cleanup(pool.Close)

	time.Sleep(2 * time.Second)

	for i := 0; i < 5; i++ {
		cp, err := pool.connection()
		require.NoError(t, err)
		st, _ := pool.cache.Get(formCacheKey(cp.address(), pool.key))
		require.True(t, assertAuthKeyForAny(st, clientKeys))
	}
}

func TestTwoFailed(t *testing.T) {
	var clientKeys []*ecdsa.PrivateKey
	mockClientBuilder := func(addr string) client {
		key := newPrivateKey(t)
		clientKeys = append(clientKeys, key)
		mockCli := newMockClient(addr, *key)
		mockCli.errOnEndpointInfo()
		return mockCli
	}

	opts := InitParameters{
		key: newPrivateKey(t),
		nodeParams: []NodeParam{
			{1, "peer0", 1},
			{1, "peer1", 1},
		},
		clientRebalanceInterval: 200 * time.Millisecond,
	}
	opts.setClientBuilder(mockClientBuilder)

	pool, err := NewPool(opts)
	require.NoError(t, err)
	err = pool.Dial(context.Background())
	require.NoError(t, err)

	t.Cleanup(pool.Close)

	time.Sleep(2 * time.Second)

	_, err = pool.connection()
	require.Error(t, err)
	require.Contains(t, err.Error(), "no healthy")
}

func TestSessionCache(t *testing.T) {
	key := newPrivateKey(t)
	expectedAuthKey := frostfsecdsa.PublicKey(key.PublicKey)

	mockClientBuilder := func(addr string) client {
		mockCli := newMockClient(addr, *key)
		mockCli.statusOnGetObject(apistatus.SessionTokenNotFound{})
		return mockCli
	}

	opts := InitParameters{
		key: newPrivateKey(t),
		nodeParams: []NodeParam{
			{1, "peer0", 1},
		},
		clientRebalanceInterval: 30 * time.Second,
	}
	opts.setClientBuilder(mockClientBuilder)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	pool, err := NewPool(opts)
	require.NoError(t, err)
	err = pool.Dial(ctx)
	require.NoError(t, err)
	t.Cleanup(pool.Close)

	// cache must contain session token
	cp, err := pool.connection()
	require.NoError(t, err)
	st, _ := pool.cache.Get(formCacheKey(cp.address(), pool.key))
	require.True(t, st.AssertAuthKey(&expectedAuthKey))

	var prm PrmObjectGet
	prm.SetAddress(oid.Address{})
	prm.UseSession(session.Object{})

	_, err = pool.GetObject(ctx, prm)
	require.Error(t, err)

	// cache must not contain session token
	cp, err = pool.connection()
	require.NoError(t, err)
	_, ok := pool.cache.Get(formCacheKey(cp.address(), pool.key))
	require.False(t, ok)

	var prm2 PrmObjectPut
	prm2.SetHeader(object.Object{})

	_, err = pool.PutObject(ctx, prm2)
	require.NoError(t, err)

	// cache must contain session token
	cp, err = pool.connection()
	require.NoError(t, err)
	st, _ = pool.cache.Get(formCacheKey(cp.address(), pool.key))
	require.True(t, st.AssertAuthKey(&expectedAuthKey))
}

func TestPriority(t *testing.T) {
	nodes := []NodeParam{
		{1, "peer0", 1},
		{2, "peer1", 100},
	}

	var clientKeys []*ecdsa.PrivateKey
	mockClientBuilder := func(addr string) client {
		key := newPrivateKey(t)
		clientKeys = append(clientKeys, key)

		if addr == nodes[0].address {
			mockCli := newMockClient(addr, *key)
			mockCli.errOnEndpointInfo()
			return mockCli
		}

		return newMockClient(addr, *key)
	}

	opts := InitParameters{
		key:                     newPrivateKey(t),
		nodeParams:              nodes,
		clientRebalanceInterval: 1500 * time.Millisecond,
	}
	opts.setClientBuilder(mockClientBuilder)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	pool, err := NewPool(opts)
	require.NoError(t, err)
	err = pool.Dial(ctx)
	require.NoError(t, err)
	t.Cleanup(pool.Close)

	expectedAuthKey1 := frostfsecdsa.PublicKey(clientKeys[0].PublicKey)
	firstNode := func() bool {
		cp, err := pool.connection()
		require.NoError(t, err)
		st, _ := pool.cache.Get(formCacheKey(cp.address(), pool.key))
		return st.AssertAuthKey(&expectedAuthKey1)
	}

	expectedAuthKey2 := frostfsecdsa.PublicKey(clientKeys[1].PublicKey)
	secondNode := func() bool {
		cp, err := pool.connection()
		require.NoError(t, err)
		st, _ := pool.cache.Get(formCacheKey(cp.address(), pool.key))
		return st.AssertAuthKey(&expectedAuthKey2)
	}
	require.Never(t, secondNode, time.Second, 200*time.Millisecond)

	require.Eventually(t, secondNode, time.Second, 200*time.Millisecond)
	require.Never(t, firstNode, time.Second, 200*time.Millisecond)
}

func TestSessionCacheWithKey(t *testing.T) {
	key := newPrivateKey(t)
	expectedAuthKey := frostfsecdsa.PublicKey(key.PublicKey)

	mockClientBuilder := func(addr string) client {
		return newMockClient(addr, *key)
	}

	opts := InitParameters{
		key: newPrivateKey(t),
		nodeParams: []NodeParam{
			{1, "peer0", 1},
		},
		clientRebalanceInterval: 30 * time.Second,
	}
	opts.setClientBuilder(mockClientBuilder)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	pool, err := NewPool(opts)
	require.NoError(t, err)
	err = pool.Dial(ctx)
	require.NoError(t, err)

	// cache must contain session token
	cp, err := pool.connection()
	require.NoError(t, err)
	st, _ := pool.cache.Get(formCacheKey(cp.address(), pool.key))
	require.True(t, st.AssertAuthKey(&expectedAuthKey))

	var prm PrmObjectGet
	prm.SetAddress(oid.Address{})
	anonKey := newPrivateKey(t)
	prm.UseKey(anonKey)

	_, err = pool.GetObject(ctx, prm)
	require.NoError(t, err)
	st, _ = pool.cache.Get(formCacheKey(cp.address(), anonKey))
	require.True(t, st.AssertAuthKey(&expectedAuthKey))
}

func TestSessionTokenOwner(t *testing.T) {
	mockClientBuilder := func(addr string) client {
		key := newPrivateKey(t)
		return newMockClient(addr, *key)
	}

	opts := InitParameters{
		key: newPrivateKey(t),
		nodeParams: []NodeParam{
			{1, "peer0", 1},
		},
	}
	opts.setClientBuilder(mockClientBuilder)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	p, err := NewPool(opts)
	require.NoError(t, err)
	err = p.Dial(ctx)
	require.NoError(t, err)
	t.Cleanup(p.Close)

	anonKey := newPrivateKey(t)
	var anonOwner user.ID
	user.IDFromKey(&anonOwner, anonKey.PublicKey)

	var prm prmCommon
	prm.UseKey(anonKey)
	var prmCtx prmContext
	prmCtx.useDefaultSession()

	var tkn session.Object
	var cc callContext
	cc.Context = ctx
	cc.sessionTarget = func(tok session.Object) {
		tkn = tok
	}
	err = p.initCallContext(&cc, prm, prmCtx)
	require.NoError(t, err)

	err = p.openDefaultSession(&cc)
	require.NoError(t, err)
	require.True(t, tkn.VerifySignature())
	require.True(t, tkn.Issuer().Equals(anonOwner))
}

func TestWaitPresence(t *testing.T) {
	mockCli := newMockClient("", *newPrivateKey(t))

	t.Run("context canceled", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		go func() {
			time.Sleep(500 * time.Millisecond)
			cancel()
		}()

		var idCnr cid.ID

		err := waitForContainerPresence(ctx, mockCli, idCnr, &WaitParams{
			timeout:      120 * time.Second,
			pollInterval: 5 * time.Second,
		})
		require.Error(t, err)
		require.Contains(t, err.Error(), "context canceled")
	})

	t.Run("context deadline exceeded", func(t *testing.T) {
		ctx := context.Background()
		var idCnr cid.ID
		err := waitForContainerPresence(ctx, mockCli, idCnr, &WaitParams{
			timeout:      500 * time.Millisecond,
			pollInterval: 5 * time.Second,
		})
		require.Error(t, err)
		require.Contains(t, err.Error(), "context deadline exceeded")
	})

	t.Run("ok", func(t *testing.T) {
		ctx := context.Background()
		var idCnr cid.ID
		err := waitForContainerPresence(ctx, mockCli, idCnr, &WaitParams{
			timeout:      10 * time.Second,
			pollInterval: 500 * time.Millisecond,
		})
		require.NoError(t, err)
	})
}

func TestStatusMonitor(t *testing.T) {
	monitor := newClientStatusMonitor(zap.NewExample(), "", 10)
	monitor.errorThreshold = 3

	count := 10
	for i := 0; i < count; i++ {
		monitor.incErrorRate()
	}

	require.Equal(t, uint64(count), monitor.overallErrorRate())
	require.Equal(t, uint32(1), monitor.currentErrorRate())
}

func TestHandleError(t *testing.T) {
	monitor := newClientStatusMonitor(zap.NewExample(), "", 10)

	for i, tc := range []struct {
		status        apistatus.Status
		err           error
		expectedError bool
		countError    bool
	}{
		{
			status:        nil,
			err:           nil,
			expectedError: false,
			countError:    false,
		},
		{
			status:        apistatus.SuccessDefaultV2{},
			err:           nil,
			expectedError: false,
			countError:    false,
		},
		{
			status:        apistatus.SuccessDefaultV2{},
			err:           errors.New("error"),
			expectedError: true,
			countError:    true,
		},
		{
			status:        nil,
			err:           errors.New("error"),
			expectedError: true,
			countError:    true,
		},
		{
			status:        apistatus.ObjectNotFound{},
			err:           nil,
			expectedError: true,
			countError:    false,
		},
		{
			status:        apistatus.ServerInternal{},
			err:           nil,
			expectedError: true,
			countError:    true,
		},
		{
			status:        apistatus.WrongMagicNumber{},
			err:           nil,
			expectedError: true,
			countError:    true,
		},
		{
			status:        apistatus.SignatureVerification{},
			err:           nil,
			expectedError: true,
			countError:    true,
		},
		{
			status:        &apistatus.SignatureVerification{},
			err:           nil,
			expectedError: true,
			countError:    true,
		},
		{
			status:        apistatus.NodeUnderMaintenance{},
			err:           nil,
			expectedError: true,
			countError:    true,
		},
	} {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			errCount := monitor.currentErrorRate()
			err := monitor.handleError(tc.status, tc.err)
			if tc.expectedError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			if tc.countError {
				errCount++
			}
			require.Equal(t, errCount, monitor.currentErrorRate())
		})
	}
}

func TestSwitchAfterErrorThreshold(t *testing.T) {
	nodes := []NodeParam{
		{1, "peer0", 1},
		{2, "peer1", 100},
	}

	errorThreshold := 5

	var clientKeys []*ecdsa.PrivateKey
	mockClientBuilder := func(addr string) client {
		key := newPrivateKey(t)
		clientKeys = append(clientKeys, key)

		if addr == nodes[0].address {
			mockCli := newMockClient(addr, *key)
			mockCli.setThreshold(uint32(errorThreshold))
			mockCli.statusOnGetObject(apistatus.ServerInternal{})
			return mockCli
		}

		return newMockClient(addr, *key)
	}

	opts := InitParameters{
		key:                     newPrivateKey(t),
		nodeParams:              nodes,
		clientRebalanceInterval: 30 * time.Second,
	}
	opts.setClientBuilder(mockClientBuilder)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	pool, err := NewPool(opts)
	require.NoError(t, err)
	err = pool.Dial(ctx)
	require.NoError(t, err)
	t.Cleanup(pool.Close)

	for i := 0; i < errorThreshold; i++ {
		conn, err := pool.connection()
		require.NoError(t, err)
		require.Equal(t, nodes[0].address, conn.address())
		_, err = conn.objectGet(ctx, PrmObjectGet{})
		require.Error(t, err)
	}

	conn, err := pool.connection()
	require.NoError(t, err)
	require.Equal(t, nodes[1].address, conn.address())
	_, err = conn.objectGet(ctx, PrmObjectGet{})
	require.NoError(t, err)
}
