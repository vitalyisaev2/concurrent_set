package set

import (
	"sync"
)

var _ Set = (*coarseGrainedSyncSet)(nil)

type coarseGrainedSyncSet struct {
	sequentialSet Set
	mutex         sync.RWMutex
}

func (c *coarseGrainedSyncSet) Insert(value int) bool {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	return c.sequentialSet.Insert(value)
}

func (c *coarseGrainedSyncSet) Contains(value int) bool {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	return c.sequentialSet.Contains(value)
}

func (c *coarseGrainedSyncSet) Remove(value int) bool {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	return c.sequentialSet.Remove(value)
}

// NewCoarseGrainedSyncSet provides thread-safe implementation of set, utilizing pessimistic locks.
func NewCoarseGrainedSyncSet() Set {
	return &coarseGrainedSyncSet{
		sequentialSet: NewSequentialSet(),
	}
}
