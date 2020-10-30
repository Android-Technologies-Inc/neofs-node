package netmap

import (
	"crypto/ecdsa"
	"encoding/hex"
	"testing"

	"github.com/nspcc-dev/neofs-api-go/pkg/netmap"
	netmapv2 "github.com/nspcc-dev/neofs-api-go/v2/netmap"
	crypto "github.com/nspcc-dev/neofs-crypto"
	"github.com/nspcc-dev/neofs-node/pkg/util/test"
	"github.com/stretchr/testify/require"
)

func TestCleanupTable(t *testing.T) {
	infos := []netmapv2.NodeInfo{
		newNodeInfo(&test.DecodeKey(1).PublicKey),
		newNodeInfo(&test.DecodeKey(2).PublicKey),
		newNodeInfo(&test.DecodeKey(3).PublicKey),
	}

	networkMap, err := netmap.NewNetmap(netmap.NodesFromV2(infos))
	require.NoError(t, err)

	mapInfos := map[string]struct{}{
		hex.EncodeToString(infos[0].GetPublicKey()): {},
		hex.EncodeToString(infos[1].GetPublicKey()): {},
		hex.EncodeToString(infos[2].GetPublicKey()): {},
	}

	t.Run("update", func(t *testing.T) {
		c := newCleanupTable(true, 1)
		c.update(networkMap, 1)
		require.Len(t, c.lastAccess, len(infos))

		for k, v := range c.lastAccess {
			require.EqualValues(t, 1, v.epoch)
			require.False(t, v.removeFlag)

			_, ok := mapInfos[k]
			require.True(t, ok)
		}

		t.Run("update with flagged", func(t *testing.T) {
			key := hex.EncodeToString(infos[0].GetPublicKey())
			c.flag(key)

			c.update(networkMap, 2)
			require.EqualValues(t, 1, c.lastAccess[key].epoch)
			require.False(t, c.lastAccess[key].removeFlag)
		})
	})

	t.Run("touch", func(t *testing.T) {
		c := newCleanupTable(true, 1)
		c.update(networkMap, 1)

		key := hex.EncodeToString(infos[1].GetPublicKey())
		require.True(t, c.touch(key, 11))
		require.EqualValues(t, 11, c.lastAccess[key].epoch)

		require.False(t, c.touch(key+"x", 12))
		require.EqualValues(t, 12, c.lastAccess[key+"x"].epoch)
	})

	t.Run("flag", func(t *testing.T) {
		c := newCleanupTable(true, 1)
		c.update(networkMap, 1)

		key := hex.EncodeToString(infos[1].GetPublicKey())
		c.flag(key)
		require.True(t, c.lastAccess[key].removeFlag)

		require.False(t, c.touch(key, 2))
		require.False(t, c.lastAccess[key].removeFlag)
	})

	t.Run("iterator", func(t *testing.T) {
		c := newCleanupTable(true, 2)
		c.update(networkMap, 1)

		t.Run("no nodes to remove", func(t *testing.T) {
			cnt := 0
			require.NoError(t,
				c.forEachRemoveCandidate(2, func(_ string) error {
					cnt++
					return nil
				}))
			require.EqualValues(t, 0, cnt)
		})

		t.Run("all nodes to remove", func(t *testing.T) {
			cnt := 0
			require.NoError(t,
				c.forEachRemoveCandidate(4, func(s string) error {
					cnt++
					_, ok := mapInfos[s]
					require.True(t, ok)
					return nil
				}))
			require.EqualValues(t, len(infos), cnt)
		})

		t.Run("some nodes to remove", func(t *testing.T) {
			cnt := 0
			key := hex.EncodeToString(infos[1].GetPublicKey())

			require.False(t, c.touch(key, 4)) // one node was updated

			require.NoError(t,
				c.forEachRemoveCandidate(4, func(s string) error {
					cnt++
					require.NotEqual(t, s, key)
					return nil
				}))
			require.EqualValues(t, len(infos)-1, cnt)
		})
	})
}

func newNodeInfo(key *ecdsa.PublicKey) (n netmapv2.NodeInfo) {
	n.SetPublicKey(crypto.MarshalPublicKey(key))
	return n
}
