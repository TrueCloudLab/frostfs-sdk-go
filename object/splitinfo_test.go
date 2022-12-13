package object_test

import (
	"crypto/rand"
	"encoding/json"
	"testing"

	objv2 "github.com/TrueCloudLab/frostfs-api-go/v2/object"
	"github.com/TrueCloudLab/frostfs-sdk-go/object"
	oid "github.com/TrueCloudLab/frostfs-sdk-go/object/id"
	"github.com/stretchr/testify/require"
)

func TestSplitInfo(t *testing.T) {
	s := object.NewSplitInfo()
	splitID := object.NewSplitID()
	lastPart := generateID()
	link := generateID()

	s.SetSplitID(splitID)
	require.Equal(t, splitID, s.SplitID())

	s.SetLastPart(lastPart)
	lp, set := s.LastPart()
	require.True(t, set)
	require.Equal(t, lastPart, lp)

	s.SetLink(link)
	l, set := s.Link()
	require.True(t, set)
	require.Equal(t, link, l)
}

func TestSplitInfoMarshal(t *testing.T) {
	testToV2 := func(t *testing.T, s *object.SplitInfo) {
		v2 := s.ToV2()
		newS := object.NewSplitInfoFromV2(v2)

		require.Equal(t, s, newS)
	}
	testMarshal := func(t *testing.T, s *object.SplitInfo) {
		data, err := s.Marshal()
		require.NoError(t, err)

		newS := object.NewSplitInfo()

		err = newS.Unmarshal(data)
		require.NoError(t, err)
		require.Equal(t, s, newS)
	}

	t.Run("good, both fields are set", func(t *testing.T) {
		s := object.NewSplitInfo()
		s.SetSplitID(object.NewSplitID())
		s.SetLink(generateID())
		s.SetLastPart(generateID())

		testToV2(t, s)
		testMarshal(t, s)
	})
	t.Run("good, only link is set", func(t *testing.T) {
		s := object.NewSplitInfo()
		s.SetSplitID(object.NewSplitID())
		s.SetLink(generateID())

		testToV2(t, s)
		testMarshal(t, s)
	})
	t.Run("good, only last part is set", func(t *testing.T) {
		s := object.NewSplitInfo()
		s.SetSplitID(object.NewSplitID())
		s.SetLastPart(generateID())

		testToV2(t, s)
		testMarshal(t, s)
	})
	t.Run("bad, no fields are set", func(t *testing.T) {
		s := object.NewSplitInfo()
		s.SetSplitID(object.NewSplitID())

		data, err := s.Marshal()
		require.NoError(t, err)
		require.Error(t, object.NewSplitInfo().Unmarshal(data))
	})
}

func generateID() oid.ID {
	var buf [32]byte
	_, _ = rand.Read(buf[:])

	var id oid.ID
	id.SetSHA256(buf)

	return id
}

func TestNewSplitInfoFromV2(t *testing.T) {
	t.Run("from nil", func(t *testing.T) {
		var x *objv2.SplitInfo

		require.Nil(t, object.NewSplitInfoFromV2(x))
	})
}

func TestSplitInfo_ToV2(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		var x *object.SplitInfo

		require.Nil(t, x.ToV2())
	})
}

func TestNewSplitInfo(t *testing.T) {
	t.Run("default values", func(t *testing.T) {
		si := object.NewSplitInfo()

		// check initial values
		require.Nil(t, si.SplitID())
		_, set := si.LastPart()
		require.False(t, set)
		_, set = si.Link()
		require.False(t, set)

		// convert to v2 message
		siV2 := si.ToV2()

		require.Nil(t, siV2.GetSplitID())
		require.Nil(t, siV2.GetLastPart())
		require.Nil(t, siV2.GetLink())
	})
}

func TestSplitInfoMarshalJSON(t *testing.T) {
	t.Run("good", func(t *testing.T) {
		s := object.NewSplitInfo()
		s.SetSplitID(object.NewSplitID())
		s.SetLastPart(generateID())
		s.SetLink(generateID())

		data, err := s.MarshalJSON()
		require.NoError(t, err)

		actual := object.NewSplitInfo()
		require.NoError(t, json.Unmarshal(data, actual))
		require.Equal(t, s, actual)
	})
	t.Run("bad link", func(t *testing.T) {
		data := `{"splitId":"Sn707289RrqDyJOrZMbMoQ==","lastPart":{"value":"Y7baWE0UdUOBr1ELKX3Q5v1LKRubQUbI81Q5UxCVeow="},"link":{"value":"bad"}}`
		require.Error(t, json.Unmarshal([]byte(data), object.NewSplitInfo()))
	})
	t.Run("bad last part", func(t *testing.T) {
		data := `{"splitId":"Sn707289RrqDyJOrZMbMoQ==","lastPart":{"value":"bad"},"link":{"value":"eRyPNCNNxHfxPcjijlv05HEcdoep/b7eHNLRSmDlnts="}}`
		require.Error(t, json.Unmarshal([]byte(data), object.NewSplitInfo()))
	})
}
