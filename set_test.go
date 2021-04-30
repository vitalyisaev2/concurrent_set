package set

import (
	"fmt"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

type setKind int8

const (
	sequential setKind = iota + 1
	coarseGrained
	fineGrained
	optimistic
	lazy
)

func (k setKind) String() string {
	switch k {
	case sequential:
		return "sequential"
	case coarseGrained:
		return "coarse_grained"
	case fineGrained:
		return "fine_grained"
	case optimistic:
		return "optimistic"
	case lazy:
		return "lazy"
	default:
		panic("unknown setKind")
	}
}

type factory struct{}

func (factory) new(k setKind) Set {
	switch k {
	case sequential:
		return NewSequentialSet()
	case coarseGrained:
		return NewCoarseGrainedSyncSet()
	case fineGrained:
		return NewFineGrainedSyncSet()
	case optimistic:
		return NewOptimisticSyncSet()
	case lazy:
		return NewLazySyncSet()
	default:
		panic("unknown setKind")
	}
}

// TestSequential verifies sequential CRUD operations of various set implementations.
func TestSequential(t *testing.T) {
	f := factory{}

	kinds := []setKind{
		sequential,
		coarseGrained,
		fineGrained,
		optimistic,
		lazy,
	}

	for _, k := range kinds {
		k := k

		t.Run(k.String(), func(t *testing.T) {
			t.Run("ascending insertion", func(t *testing.T) {
				set := f.new(k)

				// add some values
				require.True(t, set.Insert(1))
				require.True(t, set.Insert(2))
				require.True(t, set.Insert(3))

				// verify their availability
				require.True(t, set.Contains(1))
				require.True(t, set.Contains(2))
				require.True(t, set.Contains(3))

				// drop values
				require.True(t, set.Remove(1))
				require.True(t, set.Remove(2))
				require.True(t, set.Remove(3))

				// check that they no longer exist
				require.False(t, set.Contains(1))
				require.False(t, set.Contains(2))
				require.False(t, set.Contains(3))
			})

			t.Run("descending insertion", func(t *testing.T) {
				set := f.new(k)

				// add some values
				require.True(t, set.Insert(3))
				require.True(t, set.Insert(2))
				require.True(t, set.Insert(1))

				// verify their availability
				require.True(t, set.Contains(3))
				require.True(t, set.Contains(2))
				require.True(t, set.Contains(1))

				// drop values
				require.True(t, set.Remove(3))
				require.True(t, set.Remove(2))
				require.True(t, set.Remove(1))

				// check that they no longer exist
				require.False(t, set.Contains(3))
				require.False(t, set.Contains(2))
				require.False(t, set.Contains(1))
			})
		})

		t.Run("cannot insert the same value twice", func(t *testing.T) {
			set := f.new(k)

			require.True(t, set.Insert(1))
			require.False(t, set.Insert(1))

			require.True(t, set.Insert(2))
			require.False(t, set.Insert(2))
		})

		t.Run("cannot remove the same value twice", func(t *testing.T) {
			set := f.new(k)

			require.True(t, set.Insert(1))
			require.True(t, set.Insert(2))

			require.True(t, set.Remove(2))
			require.False(t, set.Remove(2))

			require.True(t, set.Remove(1))
			require.False(t, set.Remove(1))
		})
	}
}

// TestConcurrent verifies concurrent CRUD operations of various set implementations.
func TestConcurrent(t *testing.T) {
	f := factory{}

	kinds := []setKind{
		coarseGrained,
		fineGrained,
		optimistic,
		lazy,
	}

	for _, k := range kinds {
		k := k

		const (
			threads = 8
			items   = 1000
		)

		t.Run(k.String(), func(t *testing.T) {
			t.Run("concurrent operations", func(t *testing.T) {
				set := f.new(k)

				wg := sync.WaitGroup{}
				wg.Add(threads)

				// every thread tries to run concurrent insertions
				for i := 0; i < threads; i++ {
					go func() {
						defer wg.Done()

						for j := 0; j < items; j++ {
							// result is not important
							set.Insert(j)
						}
					}()
				}

				wg.Wait()

				// verify set contents
				for j := 0; j < items; j++ {
					require.True(t, set.Contains(j), j)
				}

				wg.Add(threads)

				// every thread tries to run concurrent removals
				for i := 0; i < threads; i++ {
					go func() {
						defer wg.Done()

						for j := 0; j < items; j++ {
							// result is not important
							set.Remove(j)
						}
					}()
				}

				wg.Wait()
			})
		})
	}
}

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
				val:  10,
			},
		}

		for i, tc := range testCases {
			tc := tc

			t.Run(fmt.Sprint(i), func(t *testing.T) {
				node := &nodeNonBlocking{val: tc.val}

				amr := newAtomicMarkableReference(node, tc.mark)
				amr.getNode()

				require.Equal(t, node, amr.getNode())
				require.Equal(t, tc.val, amr.getNode().val)
				require.Equal(t, tc.mark, amr.getMark())

				nodeRead, markRead := amr.getBoth()
				require.Equal(t, node, nodeRead)
				require.Equal(t, tc.val, nodeRead.val)
				require.Equal(t, tc.mark, markRead)
			})
		}
	})

	t.Run("mutation", func(t *testing.T) {
		node1 := &nodeNonBlocking{val: 1}
		node2 := &nodeNonBlocking{val: 2}
		mark1 := true
		mark2 := false

		amr := newAtomicMarkableReference(node1, mark1)

		require.True(t, amr.compareAndSet(node1, node2, mark1, mark2))

		require.Equal(t, node2, amr.getNode())
		require.Equal(t, mark2, amr.getMark())
	})
}
