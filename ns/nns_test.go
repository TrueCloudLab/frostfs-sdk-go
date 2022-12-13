package ns

import (
	"errors"
	"fmt"
	"math/big"
	"math/rand"
	"strings"
	"testing"

	"github.com/TrueCloudLab/frostfs-sdk-go/container"
	cidtest "github.com/TrueCloudLab/frostfs-sdk-go/container/id/test"
	"github.com/nspcc-dev/neo-go/pkg/neorpc/result"
	"github.com/nspcc-dev/neo-go/pkg/util"
	"github.com/nspcc-dev/neo-go/pkg/vm/stackitem"
	"github.com/nspcc-dev/neo-go/pkg/vm/vmstate"
	"github.com/stretchr/testify/require"
)

// testNeoClient represents test Neo client which checks invocation arguments
// and returns predefined result.
type testNeoClient struct {
	t *testing.T

	expectedContract util.Uint160

	res result.Invoke

	err error
}

func (x *testNeoClient) Call(contract util.Uint160, operation string, params ...interface{}) (*result.Invoke, error) {
	var domain string

	require.Equal(x.t, x.expectedContract, contract)
	require.Equal(x.t, "resolve", operation)
	require.Len(x.t, params, 2)
	require.NotPanics(x.t, func() { domain = params[0].(string) })
	require.NotPanics(x.t, func() { _ = params[1].(int64) })
	require.True(x.t, strings.HasSuffix(domain, ".container"))
	require.NotEmpty(x.t, strings.TrimSuffix(domain, ".container"))

	return &x.res, x.err
}

// implements test stackitem.Item which is obviously incorrect:
// it returns itself on Convert(stackitem.ArrayT), but returns integer from Value.
type brokenArrayStackItem struct {
	stackitem.Item
}

func (x brokenArrayStackItem) Value() interface{} {
	return 1
}

func (x brokenArrayStackItem) Convert(t stackitem.Type) (stackitem.Item, error) {
	if t != stackitem.ArrayT {
		panic(fmt.Sprintf("unexpected stack item type %s", t))
	}

	return x, nil
}

func TestNNS_ResolveContainerDomain(t *testing.T) {
	var testContainerDomain container.Domain
	testContainerDomain.SetName("some_container")

	var nnsContract util.Uint160

	rand.Read(nnsContract[:])

	testC := &testNeoClient{
		t:                t,
		expectedContract: nnsContract,
	}

	n := NNS{
		nnsContract: nnsContract,
		invoker:     testC,
	}

	t.Run("invocation failure", func(t *testing.T) {
		err1 := errors.New("invoke err")
		testC.err = err1

		_, err2 := n.ResolveContainerDomain(testContainerDomain)
		require.ErrorIs(t, err2, err1)
	})

	testC.err = nil

	t.Run("fault exception", func(t *testing.T) {
		_, err := n.ResolveContainerDomain(testContainerDomain)
		require.Error(t, err)
	})

	testC.res.State = vmstate.Halt.String()

	t.Run("empty stack", func(t *testing.T) {
		_, err := n.ResolveContainerDomain(testContainerDomain)
		require.Error(t, err)
	})

	testC.res.Stack = make([]stackitem.Item, 1)

	t.Run("non-array last stack item", func(t *testing.T) {
		testC.res.Stack[0] = stackitem.NewBigInteger(big.NewInt(11))

		_, err := n.ResolveContainerDomain(testContainerDomain)
		require.Error(t, err)
	})

	t.Run("null array", func(t *testing.T) {
		testC.res.Stack[0] = stackitem.Null{}

		_, err := n.ResolveContainerDomain(testContainerDomain)
		require.ErrorIs(t, err, errNotFound)
	})

	t.Run("array stack item with non-slice value", func(t *testing.T) {
		testC.res.Stack[0] = brokenArrayStackItem{}

		_, err := n.ResolveContainerDomain(testContainerDomain)
		require.Error(t, err)
	})

	arr := make([]stackitem.Item, 2)
	testC.res.Stack[0] = stackitem.NewArray(arr)

	t.Run("non-bytes array element", func(t *testing.T) {
		arr[0] = stackitem.NewArray(nil)

		_, err := n.ResolveContainerDomain(testContainerDomain)
		require.Error(t, err)
	})

	arr[0] = stackitem.NewByteArray([]byte("some byte array 1"))

	t.Run("non-container array elements", func(t *testing.T) {
		arr[1] = stackitem.NewByteArray([]byte("some byte array 2"))

		_, err := n.ResolveContainerDomain(testContainerDomain)
		require.Error(t, err)
	})

	t.Run("with container array element", func(t *testing.T) {
		id := cidtest.ID()

		arr[1] = stackitem.NewByteArray([]byte(id.EncodeToString()))

		res, err := n.ResolveContainerDomain(testContainerDomain)
		require.NoError(t, err)

		require.Equal(t, id, res)
	})
}
