package apistatus_test

import (
	"testing"

	"github.com/TrueCloudLab/frostfs-api-go/v2/status"
	apistatus "github.com/TrueCloudLab/frostfs-sdk-go/client/status"
	"github.com/stretchr/testify/require"
)

func TestServerInternal_Message(t *testing.T) {
	const msg = "some message"

	var st apistatus.ServerInternal

	res := st.Message()
	resv2 := apistatus.ToStatusV2(st).Message()
	require.Empty(t, res)
	require.Empty(t, resv2)

	st.SetMessage(msg)

	res = st.Message()
	resv2 = apistatus.ToStatusV2(st).Message()
	require.Equal(t, msg, res)
	require.Equal(t, msg, resv2)
}

func TestWrongMagicNumber_CorrectMagic(t *testing.T) {
	const magic = 1337

	var st apistatus.WrongMagicNumber

	res, ok := st.CorrectMagic()
	require.Zero(t, res)
	require.Zero(t, ok)

	st.WriteCorrectMagic(magic)

	res, ok = st.CorrectMagic()
	require.EqualValues(t, magic, res)
	require.EqualValues(t, 1, ok)

	// corrupt the value
	apistatus.ToStatusV2(st).IterateDetails(func(d *status.Detail) bool {
		d.SetValue([]byte{1, 2, 3}) // any slice with len != 8
		return true
	})

	_, ok = st.CorrectMagic()
	require.EqualValues(t, -1, ok)
}

func TestSignatureVerification(t *testing.T) {
	t.Run("default", func(t *testing.T) {
		var st apistatus.SignatureVerification

		require.Empty(t, st.Message())
	})

	t.Run("custom message", func(t *testing.T) {
		var st apistatus.SignatureVerification
		msg := "some message"

		st.SetMessage(msg)

		stV2 := st.ToStatusV2()

		require.Equal(t, msg, st.Message())
		require.Equal(t, msg, stV2.Message())
	})

	t.Run("empty to V2", func(t *testing.T) {
		var st apistatus.SignatureVerification

		stV2 := st.ToStatusV2()

		require.Equal(t, "signature verification failed", stV2.Message())
	})

	t.Run("non-empty to V2", func(t *testing.T) {
		var st apistatus.SignatureVerification
		msg := "some other msg"

		st.SetMessage(msg)

		stV2 := st.ToStatusV2()

		require.Equal(t, msg, stV2.Message())
	})
}

func TestNodeUnderMaintenance(t *testing.T) {
	t.Run("default", func(t *testing.T) {
		var st apistatus.NodeUnderMaintenance

		require.Empty(t, st.Message())
	})

	t.Run("custom message", func(t *testing.T) {
		var st apistatus.NodeUnderMaintenance
		msg := "some message"

		st.SetMessage(msg)

		stV2 := st.ToStatusV2()

		require.Equal(t, msg, st.Message())
		require.Equal(t, msg, stV2.Message())
	})

	t.Run("empty to V2", func(t *testing.T) {
		var st apistatus.NodeUnderMaintenance

		stV2 := st.ToStatusV2()

		require.Empty(t, "", stV2.Message())
	})

	t.Run("non-empty to V2", func(t *testing.T) {
		var st apistatus.NodeUnderMaintenance
		msg := "some other msg"

		st.SetMessage(msg)

		stV2 := st.ToStatusV2()

		require.Equal(t, msg, stV2.Message())
	})
}
