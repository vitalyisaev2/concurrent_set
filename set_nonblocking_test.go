package set

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAtomicMarkableReference(t *testing.T) {
	t.Run("construction", func(t *testing.T) {
		testCases := []struct {
			mark bool
			val  int
		}{
			{
				mark: false,
				val:  10,
			},
			{
				mark: true,
				val:  20,
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
		node1 := &nonBlockingNode{value: 1}
		node2 := &nonBlockingNode{value: 2}
		mark1 := true
		mark2 := false

		amr := newAtomicMarkableReference(node1, mark1)

		require.True(t, amr.compareAndSet(node1, node2, mark1, mark2))

		require.Equal(t, node2, amr.getNode())
		require.Equal(t, mark2, amr.getMark())
	})
}
