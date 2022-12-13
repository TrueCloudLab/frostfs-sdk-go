package netmap_test

import (
	"testing"

	"github.com/TrueCloudLab/frostfs-sdk-go/netmap"
	subnetid "github.com/TrueCloudLab/frostfs-sdk-go/subnet/id"
	"github.com/stretchr/testify/require"
)

func TestNodeInfoSubnets(t *testing.T) {
	t.Run("enter subnet", func(t *testing.T) {
		var id subnetid.ID

		id.SetNumeric(13)

		var node netmap.NodeInfo

		node.EnterSubnet(id)

		mIDs := make(map[string]struct{})

		err := node.IterateSubnets(func(id subnetid.ID) error {
			mIDs[id.String()] = struct{}{}
			return nil
		})

		require.NoError(t, err)

		_, ok := mIDs[id.String()]
		require.True(t, ok)
	})

	t.Run("iterate with removal", func(t *testing.T) {
		t.Run("not last", func(t *testing.T) {
			var id, idrm subnetid.ID

			id.SetNumeric(13)
			idrm.SetNumeric(23)

			var node netmap.NodeInfo

			node.EnterSubnet(id)
			node.EnterSubnet(idrm)

			err := node.IterateSubnets(func(id subnetid.ID) error {
				if subnetid.IsZero(id) || id.Equals(idrm) {
					return netmap.ErrRemoveSubnet
				}

				return nil
			})

			require.NoError(t, err)

			mIDs := make(map[string]struct{})

			err = node.IterateSubnets(func(id subnetid.ID) error {
				mIDs[id.String()] = struct{}{}
				return nil
			})

			require.NoError(t, err)

			var zeroID subnetid.ID

			_, ok := mIDs[zeroID.String()]
			require.False(t, ok)

			_, ok = mIDs[idrm.String()]
			require.False(t, ok)

			_, ok = mIDs[id.String()]
			require.True(t, ok)
		})

		t.Run("last", func(t *testing.T) {
			var node netmap.NodeInfo

			err := node.IterateSubnets(func(id subnetid.ID) error {
				return netmap.ErrRemoveSubnet
			})

			require.Error(t, err)
		})
	})
}

func TestEnterSubnet(t *testing.T) {
	var (
		id   subnetid.ID
		node netmap.NodeInfo
	)

	require.True(t, netmap.BelongsToSubnet(node, id))

	node.EnterSubnet(id)
	require.True(t, netmap.BelongsToSubnet(node, id))

	node.ExitSubnet(id)
	require.False(t, netmap.BelongsToSubnet(node, id))

	id.SetNumeric(10)
	node.EnterSubnet(id)
	require.True(t, netmap.BelongsToSubnet(node, id))
	require.False(t, netmap.BelongsToSubnet(node, subnetid.ID{}))

	node.ExitSubnet(id)
	require.False(t, netmap.BelongsToSubnet(node, id))
	require.False(t, netmap.BelongsToSubnet(node, subnetid.ID{}))
}

func TestBelongsToSubnet(t *testing.T) {
	var id, idMiss, idZero subnetid.ID

	id.SetNumeric(13)
	idMiss.SetNumeric(23)

	var node netmap.NodeInfo

	node.EnterSubnet(id)

	require.True(t, netmap.BelongsToSubnet(node, idZero))
	require.True(t, netmap.BelongsToSubnet(node, id))
	require.False(t, netmap.BelongsToSubnet(node, idMiss))
}
