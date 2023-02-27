package client

import (
	"context"
	"crypto/ecdsa"
	"crypto/tls"
	"errors"
	"time"

	v2accounting "github.com/TrueCloudLab/frostfs-api-go/v2/accounting"
	"github.com/TrueCloudLab/frostfs-api-go/v2/rpc"
	"github.com/TrueCloudLab/frostfs-api-go/v2/rpc/client"
)

// Client represents virtual connection to the FrostFS network to communicate
// with FrostFS server using FrostFS API protocol. It is designed to provide
// an abstraction interface from the protocol details of data transfer over
// a network in FrostFS.
//
// Client can be created using simple Go variable declaration. Before starting
// work with the Client, it SHOULD BE correctly initialized (see Init method).
// Before executing the FrostFS operations using the Client, connection to the
// server MUST BE correctly established (see Dial method and pay attention
// to the mandatory parameters). Using the Client before connecting have
// been established can lead to a panic. After the work, the Client SHOULD BE
// closed (see Close method): it frees internal and system resources which were
// allocated for the period of work of the Client. Calling Init/Dial/Close method
// during the communication process step strongly discouraged as it leads to
// undefined behavior.
//
// Each method which produces a FrostFS API call may return a server response.
// Status responses are returned in the result structure, and can be cast
// to built-in error instance (or in the returned error if the client is
// configured accordingly). Certain statuses can be checked using `apistatus`
// and standard `errors` packages. Note that package provides some helper
// functions to work with status returns (e.g. IsErrContainerNotFound).
// All possible responses are documented in methods, however, some may be
// returned from all of them (pay attention to the presence of the pointer sign):
//   - *apistatus.ServerInternal on internal server error;
//   - *apistatus.NodeUnderMaintenance if a server is under maintenance;
//   - *apistatus.SuccessDefaultV2 on default success.
//
// Client MUST NOT be copied by value: use pointer to Client instead.
//
// See client package overview to get some examples.
type Client struct {
	prm PrmInit

	c client.Client

	server frostFSAPIServer
}

// Init brings the Client instance to its initial state.
//
// One-time method call during application init stage (before Dial) is expected.
// Calling multiple times leads to undefined behavior.
//
// See docs of PrmInit methods for details. See also Dial / Close.
func (c *Client) Init(prm PrmInit) {
	c.prm = prm
}

// Dial establishes a connection to the server from the FrostFS network.
// Returns an error describing failure reason. If failed, the Client
// SHOULD NOT be used.
//
// Uses the context specified by SetContext if it was called with non-nil
// argument, otherwise context.Background() is used. Dial returns context
// errors, see context package docs for details.
//
// Returns an error if required parameters are set incorrectly, look carefully
// at the method documentation.
//
// One-time method call during application start-up stage (after Init ) is expected.
// Calling multiple times leads to undefined behavior.
//
// See also Init / Close.
func (c *Client) Dial(prm PrmDial) error {
	if prm.endpoint == "" {
		return errorServerAddrUnset
	}

	if prm.timeoutDialSet {
		if prm.timeoutDial <= 0 {
			return errorNonPositiveTimeout
		}
	} else {
		prm.timeoutDial = 5 * time.Second
	}

	if prm.streamTimeoutSet {
		if prm.streamTimeout <= 0 {
			return errorNonPositiveTimeout
		}
	} else {
		prm.streamTimeout = 10 * time.Second
	}

	c.c = *client.New(append(
		client.WithNetworkURIAddress(prm.endpoint, prm.tlsConfig),
		client.WithDialTimeout(prm.timeoutDial),
		client.WithRWTimeout(prm.streamTimeout),
	)...)

	c.setFrostFSAPIServer((*coreServer)(&c.c))

	if prm.parentCtx == nil {
		prm.parentCtx = context.Background()
	}

	// TODO: (neofs-api-go#382) perform generic dial stage of the client.Client
	_, err := rpc.Balance(&c.c, new(v2accounting.BalanceRequest),
		client.WithContext(prm.parentCtx),
	)
	// return context errors since they signal about dial problem
	if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
		return err
	}

	return nil
}

// sets underlying provider of frostFSAPIServer. The method is used for testing as an approach
// to skip Dial stage and override FrostFS API server. MUST NOT be used outside test code.
// In real applications wrapper over github.com/TrueCloudLab/frostfs-api-go/v2/rpc/client
// is statically used.
func (c *Client) setFrostFSAPIServer(server frostFSAPIServer) {
	c.server = server
}

// Close closes underlying connection to the FrostFS server. Implements io.Closer.
// MUST NOT be called before successful Dial. Can be called concurrently
// with server operations processing on running goroutines: in this case
// they are likely to fail due to a connection error.
//
// One-time method call during application shutdown stage (after Init and Dial)
// is expected. Calling multiple times leads to undefined behavior.
//
// See also Init / Dial.
func (c *Client) Close() error {
	return c.c.Conn().Close()
}

// PrmInit groups initialization parameters of Client instances.
//
// See also Init.
type PrmInit struct {
	resolveFrostFSErrors bool

	key ecdsa.PrivateKey

	cbRespInfo func(ResponseMetaInfo) error

	netMagic uint64
}

// SetDefaultPrivateKey sets Client private key to be used for the protocol
// communication by default.
//
// Required for operations without custom key parametrization (see corresponding Prm* docs).
func (x *PrmInit) SetDefaultPrivateKey(key ecdsa.PrivateKey) {
	x.key = key
}

// ResolveFrostFSFailures makes the Client to resolve failure statuses of the
// FrostFS protocol into Go built-in errors. These errors are returned from
// each protocol operation. By default, statuses aren't resolved and written
// to the resulting structure (see corresponding Res* docs).
func (x *PrmInit) ResolveFrostFSFailures() {
	x.resolveFrostFSErrors = true
}

// SetResponseInfoCallback makes the Client to pass ResponseMetaInfo from each
// FrostFS server response to f. Nil (default) means ignore response meta info.
func (x *PrmInit) SetResponseInfoCallback(f func(ResponseMetaInfo) error) {
	x.cbRespInfo = f
}

// PrmDial groups connection parameters for the Client.
//
// See also Dial.
type PrmDial struct {
	endpoint string

	tlsConfig *tls.Config

	timeoutDialSet bool
	timeoutDial    time.Duration

	streamTimeoutSet bool
	streamTimeout    time.Duration

	parentCtx context.Context
}

// SetServerURI sets server URI in the FrostFS network.
// Required parameter.
//
// Format of the URI:
//
//	[scheme://]host:port
//
// Supported schemes:
//
//	grpc
//	grpcs
//
// See also SetTLSConfig.
func (x *PrmDial) SetServerURI(endpoint string) {
	x.endpoint = endpoint
}

// SetTLSConfig sets tls.Config to open TLS client connection
// to the FrostFS server. Nil (default) means insecure connection.
//
// See also SetServerURI.
func (x *PrmDial) SetTLSConfig(tlsConfig *tls.Config) {
	x.tlsConfig = tlsConfig
}

// SetTimeout sets the timeout for connection to be established.
// MUST BE positive. If not called, 5s timeout will be used by default.
func (x *PrmDial) SetTimeout(timeout time.Duration) {
	x.timeoutDialSet = true
	x.timeoutDial = timeout
}

// SetStreamTimeout sets the timeout for individual operations in streaming RPC.
// MUST BE positive. If not called, 10s timeout will be used by default.
func (x *PrmDial) SetStreamTimeout(timeout time.Duration) {
	x.streamTimeoutSet = true
	x.streamTimeout = timeout
}

// SetContext allows to specify optional base context within which connection
// should be established.
//
// Context SHOULD NOT be nil.
func (x *PrmDial) SetContext(ctx context.Context) {
	x.parentCtx = ctx
}
