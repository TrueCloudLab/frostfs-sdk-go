package pool

//go:generate mockgen -destination mock_test.go -package pool . Client

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/nspcc-dev/neo-go/pkg/crypto/keys"
	"github.com/nspcc-dev/neofs-sdk-go/client"
	"github.com/nspcc-dev/neofs-sdk-go/netmap"
	"github.com/nspcc-dev/neofs-sdk-go/object"
	"github.com/nspcc-dev/neofs-sdk-go/object/address"
	"github.com/nspcc-dev/neofs-sdk-go/owner"
	"github.com/nspcc-dev/neofs-sdk-go/session"
	sessiontest "github.com/nspcc-dev/neofs-sdk-go/session/test"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestBuildPoolClientFailed(t *testing.T) {
	clientBuilder := func(_ string) (Client, error) {
		return nil, fmt.Errorf("error")
	}

	pb := new(Builder)
	pb.AddNode("peer0", 1, 1)

	opts := &BuilderOptions{
		Key:           newPrivateKey(t),
		clientBuilder: clientBuilder,
	}

	_, err := pb.Build(context.TODO(), opts)
	require.Error(t, err)
}

func TestBuildPoolCreateSessionFailed(t *testing.T) {
	ctrl := gomock.NewController(t)

	ni := &netmap.NodeInfo{}
	ni.SetAddresses("addr1", "addr2")

	clientBuilder := func(_ string) (Client, error) {
		mockClient := NewMockClient(ctrl)
		mockClient.EXPECT().CreateSession(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("error session")).AnyTimes()
		mockClient.EXPECT().EndpointInfo(gomock.Any(), gomock.Any()).Return(&client.ResEndpointInfo{}, nil).AnyTimes()
		mockClient.EXPECT().NetworkInfo(gomock.Any(), gomock.Any()).Return(&client.ResNetworkInfo{}, nil).AnyTimes()
		return mockClient, nil
	}

	pb := new(Builder)
	pb.AddNode("peer0", 1, 1)

	opts := &BuilderOptions{
		Key:           newPrivateKey(t),
		clientBuilder: clientBuilder,
	}

	_, err := pb.Build(context.TODO(), opts)
	require.Error(t, err)
}

func newPrivateKey(t *testing.T) *ecdsa.PrivateKey {
	p, err := keys.NewPrivateKey()
	require.NoError(t, err)
	return &p.PrivateKey
}

func TestBuildPoolOneNodeFailed(t *testing.T) {
	t.Skip("NeoFS API client can't be mocked") // neofs-sdk-go#85

	ctrl := gomock.NewController(t)
	ctrl2 := gomock.NewController(t)

	ni := &netmap.NodeInfo{}
	ni.SetAddresses("addr1", "addr2")

	var expectedToken *session.Token
	clientCount := -1
	clientBuilder := func(_ string) (Client, error) {
		clientCount++
		mockClient := NewMockClient(ctrl)
		mockInvokes := 0
		mockClient.EXPECT().CreateSession(gomock.Any(), gomock.Any()).DoAndReturn(func(_, _ interface{}, _ ...interface{}) (*session.Token, error) {
			mockInvokes++
			if mockInvokes == 1 {
				expectedToken = newToken(t)
				return nil, fmt.Errorf("error session")
			}
			return expectedToken, nil
		}).AnyTimes()

		mockClient.EXPECT().EndpointInfo(gomock.Any(), gomock.Any()).Return(nil, nil).AnyTimes()
		mockClient.EXPECT().NetworkInfo(gomock.Any(), gomock.Any()).Return(&client.ResNetworkInfo{}, nil).AnyTimes()

		mockClient2 := NewMockClient(ctrl2)
		mockClient2.EXPECT().CreateSession(gomock.Any(), gomock.Any()).Return(nil, nil).AnyTimes()
		mockClient2.EXPECT().EndpointInfo(gomock.Any(), gomock.Any()).Return(nil, nil).AnyTimes()
		mockClient.EXPECT().NetworkInfo(gomock.Any(), gomock.Any()).Return(&client.ResNetworkInfo{}, nil).AnyTimes()

		if clientCount == 0 {
			return mockClient, nil
		}
		return mockClient2, nil
	}

	pb := new(Builder)
	pb.AddNode("peer0", 9, 1)
	pb.AddNode("peer1", 1, 1)

	log, err := zap.NewProduction()
	require.NoError(t, err)
	opts := &BuilderOptions{
		Key:                     newPrivateKey(t),
		clientBuilder:           clientBuilder,
		ClientRebalanceInterval: 1000 * time.Millisecond,
		Logger:                  log,
	}

	clientPool, err := pb.Build(context.TODO(), opts)
	require.NoError(t, err)
	t.Cleanup(clientPool.Close)

	condition := func() bool {
		_, st, err := clientPool.Connection()
		return err == nil && st == expectedToken
	}
	require.Never(t, condition, 900*time.Millisecond, 100*time.Millisecond)
	require.Eventually(t, condition, 3*time.Second, 300*time.Millisecond)
}

func TestBuildPoolZeroNodes(t *testing.T) {
	pb := new(Builder)
	opts := &BuilderOptions{
		Key: newPrivateKey(t),
	}
	_, err := pb.Build(context.TODO(), opts)
	require.Error(t, err)
}

func TestOneNode(t *testing.T) {
	t.Skip("NeoFS API client can't be mocked") // neofs-sdk-go#85

	ctrl := gomock.NewController(t)

	tok := session.NewToken()
	uid, err := uuid.New().MarshalBinary()
	require.NoError(t, err)
	tok.SetID(uid)

	clientBuilder := func(_ string) (Client, error) {
		mockClient := NewMockClient(ctrl)
		mockClient.EXPECT().CreateSession(gomock.Any(), gomock.Any()).Return(tok, nil)
		mockClient.EXPECT().EndpointInfo(gomock.Any(), gomock.Any()).Return(&client.ResEndpointInfo{}, nil).AnyTimes()
		mockClient.EXPECT().NetworkInfo(gomock.Any(), gomock.Any()).Return(&client.ResNetworkInfo{}, nil).AnyTimes()
		return mockClient, nil
	}

	pb := new(Builder)
	pb.AddNode("peer0", 1, 1)

	opts := &BuilderOptions{
		Key:           newPrivateKey(t),
		clientBuilder: clientBuilder,
	}

	pool, err := pb.Build(context.Background(), opts)
	require.NoError(t, err)
	t.Cleanup(pool.Close)

	_, st, err := pool.Connection()
	require.NoError(t, err)
	require.Equal(t, tok, st)
}

func TestTwoNodes(t *testing.T) {
	t.Skip("NeoFS API client can't be mocked") // neofs-sdk-go#85

	ctrl := gomock.NewController(t)

	var tokens []*session.Token
	clientBuilder := func(_ string) (Client, error) {
		mockClient := NewMockClient(ctrl)
		mockClient.EXPECT().CreateSession(gomock.Any(), gomock.Any()).DoAndReturn(func(_, _ interface{}, _ ...interface{}) (*session.Token, error) {
			tok := session.NewToken()
			uid, err := uuid.New().MarshalBinary()
			require.NoError(t, err)
			tok.SetID(uid)
			tokens = append(tokens, tok)
			return tok, err
		})
		mockClient.EXPECT().EndpointInfo(gomock.Any(), gomock.Any()).Return(&client.ResEndpointInfo{}, nil).AnyTimes()
		mockClient.EXPECT().NetworkInfo(gomock.Any(), gomock.Any()).Return(&client.ResNetworkInfo{}, nil).AnyTimes()
		return mockClient, nil
	}

	pb := new(Builder)
	pb.AddNode("peer0", 1, 1)
	pb.AddNode("peer1", 1, 1)

	opts := &BuilderOptions{
		Key:           newPrivateKey(t),
		clientBuilder: clientBuilder,
	}

	pool, err := pb.Build(context.Background(), opts)
	require.NoError(t, err)
	t.Cleanup(pool.Close)

	_, st, err := pool.Connection()
	require.NoError(t, err)
	require.Contains(t, tokens, st)
}

func TestOneOfTwoFailed(t *testing.T) {
	t.Skip("NeoFS API client can't be mocked") // neofs-sdk-go#85

	ctrl := gomock.NewController(t)
	ctrl2 := gomock.NewController(t)

	var tokens []*session.Token
	clientCount := -1
	clientBuilder := func(_ string) (Client, error) {
		clientCount++
		mockClient := NewMockClient(ctrl)
		mockClient.EXPECT().CreateSession(gomock.Any(), gomock.Any()).DoAndReturn(func(_, _ interface{}, _ ...interface{}) (*session.Token, error) {
			tok := newToken(t)
			tokens = append(tokens, tok)
			return tok, nil
		}).AnyTimes()
		mockClient.EXPECT().EndpointInfo(gomock.Any(), gomock.Any()).Return(nil, nil).AnyTimes()
		mockClient.EXPECT().NetworkInfo(gomock.Any(), gomock.Any()).Return(&client.ResNetworkInfo{}, nil).AnyTimes()

		mockClient2 := NewMockClient(ctrl2)
		mockClient2.EXPECT().CreateSession(gomock.Any(), gomock.Any()).DoAndReturn(func(_, _ interface{}, _ ...interface{}) (*session.Token, error) {
			tok := newToken(t)
			tokens = append(tokens, tok)
			return tok, nil
		}).AnyTimes()
		mockClient2.EXPECT().EndpointInfo(gomock.Any(), gomock.Any()).DoAndReturn(func(_ interface{}, _ ...interface{}) (*client.ResEndpointInfo, error) {
			return nil, fmt.Errorf("error")
		}).AnyTimes()
		mockClient2.EXPECT().NetworkInfo(gomock.Any(), gomock.Any()).DoAndReturn(func(_ interface{}, _ ...interface{}) (*client.ResNetworkInfo, error) {
			return nil, fmt.Errorf("error")
		}).AnyTimes()

		if clientCount == 0 {
			return mockClient, nil
		}
		return mockClient2, nil
	}

	pb := new(Builder)
	pb.AddNode("peer0", 1, 1)
	pb.AddNode("peer1", 9, 1)

	opts := &BuilderOptions{
		Key:                     newPrivateKey(t),
		clientBuilder:           clientBuilder,
		ClientRebalanceInterval: 200 * time.Millisecond,
	}

	pool, err := pb.Build(context.Background(), opts)
	require.NoError(t, err)
	t.Cleanup(pool.Close)

	time.Sleep(2 * time.Second)

	for i := 0; i < 5; i++ {
		_, st, err := pool.Connection()
		require.NoError(t, err)
		require.Equal(t, tokens[0], st)
	}
}

func TestTwoFailed(t *testing.T) {
	t.Skip("NeoFS API client can't be mocked") // neofs-sdk-go#85

	ctrl := gomock.NewController(t)

	clientBuilder := func(_ string) (Client, error) {
		mockClient := NewMockClient(ctrl)
		mockClient.EXPECT().CreateSession(gomock.Any(), gomock.Any()).Return(nil, nil).AnyTimes()
		mockClient.EXPECT().EndpointInfo(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("error")).AnyTimes()
		mockClient.EXPECT().NetworkInfo(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("error")).AnyTimes()
		return mockClient, nil
	}

	pb := new(Builder)
	pb.AddNode("peer0", 1, 1)
	pb.AddNode("peer1", 1, 1)

	opts := &BuilderOptions{
		Key:                     newPrivateKey(t),
		clientBuilder:           clientBuilder,
		ClientRebalanceInterval: 200 * time.Millisecond,
	}

	pool, err := pb.Build(context.Background(), opts)
	require.NoError(t, err)
	t.Cleanup(pool.Close)

	time.Sleep(2 * time.Second)

	_, _, err = pool.Connection()
	require.Error(t, err)
	require.Contains(t, err.Error(), "no healthy")
}

func TestSessionCache(t *testing.T) {
	t.Skip("NeoFS API client can't be mocked") // neofs-sdk-go#85

	ctrl := gomock.NewController(t)

	var tokens []*session.Token
	clientBuilder := func(_ string) (Client, error) {
		mockClient := NewMockClient(ctrl)
		mockClient.EXPECT().CreateSession(gomock.Any(), gomock.Any()).DoAndReturn(func(_, _ interface{}, _ ...interface{}) (*session.Token, error) {
			tok := session.NewToken()
			uid, err := uuid.New().MarshalBinary()
			require.NoError(t, err)
			tok.SetID(uid)
			tokens = append(tokens, tok)
			return tok, err
		}).MaxTimes(3)

		mockClient.EXPECT().GetObject(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("session token does not exist"))
		mockClient.EXPECT().PutObject(gomock.Any(), gomock.Any()).Return(nil, nil)

		return mockClient, nil
	}

	pb := new(Builder)
	pb.AddNode("peer0", 1, 1)

	opts := &BuilderOptions{
		Key:                     newPrivateKey(t),
		clientBuilder:           clientBuilder,
		ClientRebalanceInterval: 30 * time.Second,
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	pool, err := pb.Build(ctx, opts)
	require.NoError(t, err)
	t.Cleanup(pool.Close)

	// cache must contain session token
	_, st, err := pool.Connection()
	require.NoError(t, err)
	require.Contains(t, tokens, st)

	_, err = pool.GetObject(ctx, address.Address{})
	require.Error(t, err)

	// cache must not contain session token
	_, st, err = pool.Connection()
	require.NoError(t, err)
	require.Nil(t, st)

	_, err = pool.PutObject(ctx, object.Object{}, nil)
	require.NoError(t, err)

	// cache must contain session token
	_, st, err = pool.Connection()
	require.NoError(t, err)
	require.Contains(t, tokens, st)
}

func TestPriority(t *testing.T) {
	t.Skip("NeoFS API client can't be mocked") // neofs-sdk-go#85

	ctrl := gomock.NewController(t)
	ctrl2 := gomock.NewController(t)

	var tokens []*session.Token
	clientCount := -1
	clientBuilder := func(_ string) (Client, error) {
		clientCount++
		mockClient := NewMockClient(ctrl)
		mockClient.EXPECT().CreateSession(gomock.Any(), gomock.Any()).DoAndReturn(func(_, _ interface{}, _ ...interface{}) (*session.Token, error) {
			tok := newToken(t)
			tokens = append(tokens, tok)
			return tok, nil
		}).AnyTimes()
		mockClient.EXPECT().EndpointInfo(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("error")).AnyTimes()
		mockClient.EXPECT().NetworkInfo(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("error")).AnyTimes()

		mockClient2 := NewMockClient(ctrl2)
		mockClient2.EXPECT().CreateSession(gomock.Any(), gomock.Any()).DoAndReturn(func(_, _ interface{}, _ ...interface{}) (*session.Token, error) {
			tok := newToken(t)
			tokens = append(tokens, tok)
			return tok, nil
		}).AnyTimes()
		mockClient2.EXPECT().EndpointInfo(gomock.Any(), gomock.Any()).Return(nil, nil).AnyTimes()
		mockClient2.EXPECT().NetworkInfo(gomock.Any(), gomock.Any()).Return(nil, nil).AnyTimes()

		if clientCount == 0 {
			return mockClient, nil
		}
		return mockClient2, nil
	}

	pb := new(Builder)
	pb.AddNode("peer0", 1, 1)
	pb.AddNode("peer1", 2, 100)

	opts := &BuilderOptions{
		Key:                     newPrivateKey(t),
		clientBuilder:           clientBuilder,
		ClientRebalanceInterval: 1500 * time.Millisecond,
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	pool, err := pb.Build(ctx, opts)
	require.NoError(t, err)
	t.Cleanup(pool.Close)

	firstNode := func() bool {
		_, st, err := pool.Connection()
		require.NoError(t, err)
		return st == tokens[0]
	}
	secondNode := func() bool {
		_, st, err := pool.Connection()
		require.NoError(t, err)
		return st == tokens[1]
	}
	require.Never(t, secondNode, time.Second, 200*time.Millisecond)

	require.Eventually(t, secondNode, time.Second, 200*time.Millisecond)
	require.Never(t, firstNode, time.Second, 200*time.Millisecond)
}

func TestSessionCacheWithKey(t *testing.T) {
	t.Skip("NeoFS API client can't be mocked") // neofs-sdk-go#85

	ctrl := gomock.NewController(t)

	var tokens []*session.Token
	clientBuilder := func(_ string) (Client, error) {
		mockClient := NewMockClient(ctrl)
		mockClient.EXPECT().CreateSession(gomock.Any(), gomock.Any()).DoAndReturn(func(_, _ interface{}, _ ...interface{}) (*session.Token, error) {
			tok := session.NewToken()
			uid, err := uuid.New().MarshalBinary()
			require.NoError(t, err)
			tok.SetID(uid)
			tokens = append(tokens, tok)
			return tok, err
		}).MaxTimes(2)

		mockClient.EXPECT().GetObject(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil)

		return mockClient, nil
	}

	pb := new(Builder)
	pb.AddNode("peer0", 1, 1)

	opts := &BuilderOptions{
		Key:                     newPrivateKey(t),
		clientBuilder:           clientBuilder,
		ClientRebalanceInterval: 30 * time.Second,
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	pool, err := pb.Build(ctx, opts)
	require.NoError(t, err)

	// cache must contain session token
	_, st, err := pool.Connection()
	require.NoError(t, err)
	require.Contains(t, tokens, st)

	_, err = pool.GetObject(ctx, address.Address{}, WithKey(newPrivateKey(t)))
	require.NoError(t, err)
	require.Len(t, tokens, 2)
}

func newToken(t *testing.T) *session.Token {
	tok := session.NewToken()
	uid, err := uuid.New().MarshalBinary()
	require.NoError(t, err)
	tok.SetID(uid)

	return tok
}

func TestSessionTokenOwner(t *testing.T) {
	ctrl := gomock.NewController(t)
	clientBuilder := func(_ string) (Client, error) {
		mockClient := NewMockClient(ctrl)
		mockClient.EXPECT().CreateSession(gomock.Any(), gomock.Any()).Return(&client.ResSessionCreate{}, nil).AnyTimes()
		mockClient.EXPECT().EndpointInfo(gomock.Any(), gomock.Any()).Return(&client.ResEndpointInfo{}, nil).AnyTimes()
		mockClient.EXPECT().NetworkInfo(gomock.Any(), gomock.Any()).Return(&client.ResNetworkInfo{}, nil).AnyTimes()
		return mockClient, nil
	}

	pb := new(Builder)
	pb.AddNode("peer0", 1, 1)

	opts := &BuilderOptions{
		Key:           newPrivateKey(t),
		clientBuilder: clientBuilder,
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	p, err := pb.Build(ctx, opts)
	require.NoError(t, err)
	t.Cleanup(p.Close)

	anonKey := newPrivateKey(t)
	anonOwner := owner.NewIDFromPublicKey(&anonKey.PublicKey)

	cfg := cfgFromOpts(WithKey(anonKey), useDefaultSession())
	cp, err := p.conn(ctx, cfg)
	require.NoError(t, err)

	tkn := p.cache.Get(formCacheKey(cp.address, anonKey))
	require.True(t, anonOwner.Equal(tkn.OwnerID()))
}

func TestWaitPresence(t *testing.T) {
	ctrl := gomock.NewController(t)

	mockClient := NewMockClient(ctrl)
	mockClient.EXPECT().CreateSession(gomock.Any(), gomock.Any()).Return(nil, nil).AnyTimes()
	mockClient.EXPECT().EndpointInfo(gomock.Any(), gomock.Any()).Return(nil, nil).AnyTimes()
	mockClient.EXPECT().NetworkInfo(gomock.Any(), gomock.Any()).Return(&client.ResNetworkInfo{}, nil).AnyTimes()
	mockClient.EXPECT().GetContainer(gomock.Any(), gomock.Any()).Return(nil, nil).AnyTimes()

	cache, err := NewCache()
	require.NoError(t, err)

	inner := &innerPool{
		sampler: NewSampler([]float64{1}, rand.NewSource(0)),
		clientPacks: []*clientPack{{
			client:  mockClient,
			healthy: true,
		}},
	}

	p := &Pool{
		innerPools: []*innerPool{inner},
		key:        newPrivateKey(t),
		cache:      cache,
	}

	t.Run("context canceled", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		go func() {
			time.Sleep(500 * time.Millisecond)
			cancel()
		}()

		err := p.WaitForContainerPresence(ctx, nil, DefaultPollingParams())
		require.Error(t, err)
		require.Contains(t, err.Error(), "context canceled")
	})

	t.Run("context deadline exceeded", func(t *testing.T) {
		ctx := context.Background()
		err := p.WaitForContainerPresence(ctx, nil, &ContainerPollingParams{
			CreationTimeout: 500 * time.Millisecond,
			PollInterval:    5 * time.Second,
		})
		require.Error(t, err)
		require.Contains(t, err.Error(), "context deadline exceeded")
	})

	t.Run("ok", func(t *testing.T) {
		ctx := context.Background()
		err := p.WaitForContainerPresence(ctx, nil, &ContainerPollingParams{
			CreationTimeout: 10 * time.Second,
			PollInterval:    500 * time.Millisecond,
		})
		require.NoError(t, err)
	})
}

func TestCopySessionTokenWithoutSignatureAndContext(t *testing.T) {
	from := sessiontest.SignedToken()
	to := copySessionTokenWithoutSignatureAndContext(*from)

	require.Equal(t, from.Nbf(), to.Nbf())
	require.Equal(t, from.Exp(), to.Exp())
	require.Equal(t, from.Iat(), to.Iat())
	require.Equal(t, from.ID(), to.ID())
	require.Equal(t, from.OwnerID().String(), to.OwnerID().String())
	require.Equal(t, from.SessionKey(), to.SessionKey())

	require.Empty(t, to.Signature().Sign())
	require.Empty(t, to.Signature().Key())

	t.Run("empty object context", func(t *testing.T) {
		octx := sessiontest.ObjectContext()
		from.SetContext(octx)
		to = copySessionTokenWithoutSignatureAndContext(*from)
		require.Nil(t, to.Context())
	})

	t.Run("empty container context", func(t *testing.T) {
		cctx := sessiontest.ContainerContext()
		from.SetContext(cctx)
		to = copySessionTokenWithoutSignatureAndContext(*from)
		require.Nil(t, to.Context())
	})
}
