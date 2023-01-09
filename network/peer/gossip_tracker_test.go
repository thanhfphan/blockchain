package peer

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/thanhfphan/blockchain/ids"
)

var (
	p1, _ = ids.NodeIDFromString(ids.NodeIDPrefix + "1QCaxc8hutpdZ62iKZsn1TCG3nh7uPZojqaaaaaaap1")
	p2, _ = ids.NodeIDFromString(ids.NodeIDPrefix + "1QCaxc8hutpdZ62iKZsn1TCG3nh7uPZojqaaaaaaap2")
	p3, _ = ids.NodeIDFromString(ids.NodeIDPrefix + "1QCaxc8hutpdZ62iKZsn1TCG3nh7uPZojqaaaaaaap3")

	pv1, _ = ids.NodeIDFromString(ids.NodeIDPrefix + "1QCaxc8hutpdZ62iKZsn1TCG3nh7uPZojqaaaaaapv1")
	pv2, _ = ids.NodeIDFromString(ids.NodeIDPrefix + "1QCaxc8hutpdZ62iKZsn1TCG3nh7uPZojqaaaaaapv2")
	pv3, _ = ids.NodeIDFromString(ids.NodeIDPrefix + "1QCaxc8hutpdZ62iKZsn1TCG3nh7uPZojqaaaaaapv3")

	t1, _ = ids.IDFromString("1QCaxc8hutpdZ62iKZsn1TCG3nh7uPZojqaaaaaaat1")
	t2, _ = ids.IDFromString("1QCaxc8hutpdZ62iKZsn1TCG3nh7uPZojqaaaaaaat2")
	t3, _ = ids.IDFromString("1QCaxc8hutpdZ62iKZsn1TCG3nh7uPZojqaaaaaaat3")

	v1 = ids.ValidatorID{
		NodeID: pv1,
		TxID:   t1,
	}

	v2 = ids.ValidatorID{
		NodeID: pv2,
		TxID:   t2,
	}

	v3 = ids.ValidatorID{
		NodeID: pv3,
		TxID:   t3,
	}
)

func TestGossip_E2E(t *testing.T) {
	r := require.New(t)

	g, err := NewGossipTracker()
	r.NoError(err)

	r.True(g.AddValidator(v1))
	r.True(g.AddValidator(v2))

	// we should get an empty unknown since we're not tracking anything
	unknown, ok := g.GetUnknown(p1)
	r.False(ok)
	r.Nil(unknown)

	// we should get a unknown of [v1, v2] since v1 and v2 are registered
	r.True(g.StartTrackingPeer(p1))

	// check p1's unknown
	unknown, ok = g.GetUnknown(p1)
	r.True(ok)
	r.Contains(unknown, v1)
	r.Contains(unknown, v2)
	r.Len(unknown, 2)

	// Check p2's unknown. We should get nothing since we're not tracking it
	// yet.
	unknown, ok = g.GetUnknown(p2)
	r.False(ok)
	r.Nil(unknown)

	// Start tracking p2
	r.True(g.StartTrackingPeer(p2))

	// check p2's unknown
	unknown, ok = g.GetUnknown(p2)
	r.True(ok)
	r.Contains(unknown, v1)
	r.Contains(unknown, v2)
	r.Len(unknown, 2)

	// p1 now knows about v1, but not v2, so it should see [v2] in its unknown
	// p2 still knows nothing, so it should see both
	txIDs, ok := g.AddKnown(p1, []ids.ID{v1.TxID})
	r.True(ok)
	r.Equal([]ids.ID{v1.TxID}, txIDs)

	// p1 should have an unknown of [v2], since it knows v1
	unknown, ok = g.GetUnknown(p1)
	r.True(ok)
	r.Contains(unknown, v2)
	r.Len(unknown, 1)

	// p2 should have a unknown of [v1, v2], since it knows nothing
	unknown, ok = g.GetUnknown(p2)
	r.True(ok)
	r.Contains(unknown, v1)
	r.Contains(unknown, v2)
	r.Len(unknown, 2)

	// Add v3
	r.True(g.AddValidator(v3))

	// track p3, who knows of v1, v2, and v3
	// p1 and p2 still don't know of v3
	r.True(g.StartTrackingPeer(p3))

	txIDs, ok = g.AddKnown(p3, []ids.ID{v1.TxID, v2.TxID, v3.TxID})
	r.True(ok)
	r.Equal([]ids.ID{v1.TxID, v2.TxID, v3.TxID}, txIDs)

	// p1 doesn't know about [v2, v3]
	unknown, ok = g.GetUnknown(p1)
	r.True(ok)
	r.Contains(unknown, v2)
	r.Contains(unknown, v3)
	r.Len(unknown, 2)

	// p2 doesn't know about [v1, v2, v3]
	unknown, ok = g.GetUnknown(p2)
	r.True(ok)
	r.Contains(unknown, v1)
	r.Contains(unknown, v2)
	r.Contains(unknown, v3)
	r.Len(unknown, 3)

	// p3 knows about everyone
	unknown, ok = g.GetUnknown(p3)
	r.True(ok)
	r.Empty(unknown)

	// stop tracking p2
	r.True(g.StopTrackingPeer(p2))
	unknown, ok = g.GetUnknown(p2)
	r.False(ok)
	r.Nil(unknown)

	// p1 doesn't know about [v2, v3] because v2 is still registered as
	// a validator
	unknown, ok = g.GetUnknown(p1)
	r.True(ok)
	r.Contains(unknown, v2)
	r.Contains(unknown, v3)
	r.Len(unknown, 2)

	// Remove p2 from the validator set
	r.True(g.RemoveValidator(v2.NodeID))

	// p1 doesn't know about [v3] since v2 left the validator set
	unknown, ok = g.GetUnknown(p1)
	r.True(ok)
	r.Contains(unknown, v3)
	r.Len(unknown, 1)

	// p3 knows about everyone since it learned about v1 and v3 earlier.
	unknown, ok = g.GetUnknown(p3)
	r.Empty(unknown)
	r.True(ok)
}
