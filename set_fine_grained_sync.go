package set

import (
	"math"
	"sync"
)

type syncNode struct {
	next *syncNode
	sync.Mutex
	value int
}

var _ Set = (*fineGrainedSyncSet)(nil)

type fineGrainedSyncSet struct {
	head *syncNode
}

func (s *fineGrainedSyncSet) Insert(value int) bool {
	// it looks impossible to use defers here
	s.head.Lock()

	pred := s.head
	curr := pred.next

	curr.Lock()

	for curr.value < value {
		pred.Unlock()

		pred = curr
		curr = curr.next

		curr.Lock()
	}

	defer func() {
		curr.Unlock()
		pred.Unlock()
	}()

	if curr.value == value {
		return false
	}

	newNode := &syncNode{value: value, next: curr}
	pred.next = newNode

	return true
}

func (s fineGrainedSyncSet) Contains(value int) bool {
	s.head.Lock()

	pred := s.head
	curr := pred.next

	curr.Lock()

	for curr.value < value {
		pred.Unlock()

		pred = curr
		curr = curr.next

		curr.Lock()
	}

	defer func() {
		curr.Unlock()
		pred.Unlock()
	}()

	return curr.value == value
}

func (s *fineGrainedSyncSet) Remove(value int) bool {
	s.head.Lock()

	pred := s.head
	curr := pred.next

	curr.Lock()

	for curr.value < value {
		pred.Unlock()
		pred = curr
		curr = pred.next
		curr.Lock()
	}

	defer func() {
		curr.Unlock()
		pred.Unlock()
	}()

	if curr.value == value {
		pred.next = curr.next

		return true
	}

	return false
}

// NewFineGrainedSyncSet provides more optimal thread-safe set implementation with a mutex in every list node.
func NewFineGrainedSyncSet() Set {
	// set must contain sentinel nodes with minimal and maximal values
	s := &fineGrainedSyncSet{}
	s.head = &syncNode{value: -math.MaxInt64}
	s.head.next = &syncNode{value: math.MaxInt64}

	return s
}
