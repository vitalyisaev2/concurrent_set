package set

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSet(t *testing.T) {
	t.Run("basic operations", func(t *testing.T) {
		set := NewSequentialSet()

		t.Run("ascending insertion", func(t *testing.T) {
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
		set := NewSequentialSet()

		require.True(t, set.Insert(1))
		require.False(t, set.Insert(1))

		require.True(t, set.Insert(2))
		require.False(t, set.Insert(2))
	})

	t.Run("cannot remove the same value twice", func(t *testing.T) {
		set := NewSequentialSet()

		require.True(t, set.Insert(1))
		require.True(t, set.Insert(2))

		require.True(t, set.Remove(2))
		require.False(t, set.Remove(2))

		require.True(t, set.Remove(1))
		require.False(t, set.Remove(1))
	})
}
