package set

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

type kind int8

const (
	sequential kind = iota + 1
	coarseGrained
	fineGrained
)

func (k kind) String() string {
	switch k {
	case sequential:
		return "sequential"
	case coarseGrained:
		return "coarse grained"
	case fineGrained:
		return "fine grained"
	default:
		panic("unknown set kind")
	}
}

type factory struct{}

func (factory) new(k kind) Set {
	switch k {
	case sequential:
		return NewSequentialSet()
	case coarseGrained:
		return NewCoarseGrainedSyncSet()
	case fineGrained:
		return NewFineGrainedSyncSet()
	default:
		panic("unknown set kind")
	}
}

// TestSequential verifies sequential CRUD operations on various set implementations
func TestSequential(t *testing.T) {
	f := factory{}

	kinds := []kind{
		sequential,
		fineGrained,
		coarseGrained,
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

// TestCRUDConcurrent verifies concurrent CRUD operations on various set implementations
func TestCRUDConcurrent(t *testing.T) {
	f := factory{}

	kinds := []kind{
		fineGrained,
		coarseGrained,
	}

	for _, k := range kinds {
		k := k

		const (
			threads = 10
			items   = 100
		)

		t.Run(k.String(), func(t *testing.T) {
			t.Run("concurrent insertion", func(t *testing.T) {
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
			})
		})
	}
}
