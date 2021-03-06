package set

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestAtomicMarkableReference(t *testing.T) {
	rand.Seed(time.Now().Unix())

	t.Run("construction", func(t *testing.T) {
		testCases := []struct {
			mark bool
			val  int
		}{
			{
				mark: false,
				val:  rand.Int(),
			},
			{
				mark: true,
				val:  rand.Int(),
			},
		}

		for i, tc := range testCases {
			tc := tc

			t.Run(fmt.Sprint(i), func(t *testing.T) {
				node := &nonBlockingNode{value: tc.val}

				amr := newAtomicMarkableReference(node, tc.mark)
				amr.getNode()

				require.Equal(t, node, amr.getNode())
				require.Equal(t, tc.val, amr.getNode().value)
				require.Equal(t, tc.mark, amr.getMark())

				nodeRead, markRead := amr.getBoth()
				require.Equal(t, node, nodeRead)
				require.Equal(t, tc.val, nodeRead.value)
				require.Equal(t, tc.mark, markRead)
			})
		}
	})

	t.Run("mutation", func(t *testing.T) {
		node1 := &nonBlockingNode{value: rand.Int()}
		node2 := &nonBlockingNode{value: rand.Int()}
		mark1 := true
		mark2 := false

		amr := newAtomicMarkableReference(node1, mark1)

		require.True(t, amr.compareAndSet(node1, node2, mark1, mark2))

		require.Equal(t, node2, amr.getNode())
		require.Equal(t, mark2, amr.getMark())
	})
}
