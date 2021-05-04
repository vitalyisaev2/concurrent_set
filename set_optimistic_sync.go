package set

import (
	"math"
)

var _ Set = (*optimisticSyncSet)(nil)

type optimisticSyncSet struct {
	head *syncNode
}

func (s *optimisticSyncSet) Insert(value int) bool {
	for {
		result, repeat := s.insertLoopBody(value)
		if !repeat {
			return result
		}
	}
}

//nolint:dupl // it's better to copy-paste code than messing with inheritance
func (s *optimisticSyncSet) insertLoopBody(value int) (result, repeat bool) {
	pred := s.head
	curr := pred.next

	for curr.value < value {
		pred = curr
		curr = curr.next
	}

	pred.Lock()
	curr.Lock()

	defer func() {
		curr.Unlock()
		pred.Unlock()
	}()

	if s.validate(pred, curr) {
		if curr.value == value {
			return false, false
		}

		newNode := &syncNode{value: value, next: curr}
		pred.next = newNode

		return true, false
	}

	return false, true
}

func (s *optimisticSyncSet) Contains(value int) bool {
	for {
		result, repeat := s.containsLoopBody(value)
		if !repeat {
			return result
		}
	}
}

func (s *optimisticSyncSet) containsLoopBody(value int) (result, repeat bool) {
	pred := s.head
	curr := pred.next

	for curr.value < value {
		pred = curr
		curr = curr.next
	}

	pred.Lock()
	curr.Lock()

	defer func() {
		curr.Unlock()
		pred.Unlock()
	}()

	if s.validate(pred, curr) {
		return curr.value == value, false
	}

	return false, true
}

func (s *optimisticSyncSet) Remove(value int) bool {
	for {
		result, repeat := s.removeLoopBody(value)
		if !repeat {
			return result
		}
	}
}

func (s *optimisticSyncSet) removeLoopBody(value int) (result, repeat bool) {
	pred := s.head
	curr := s.head.next

	for curr.value < value {
		pred = curr
		curr = curr.next
	}

	pred.Lock()
	curr.Lock()

	defer func() {
		curr.Unlock()
		pred.Unlock()
	}()

	if s.validate(pred, curr) {
		if curr.value == value {
			pred.next = curr.next

			return true, false
		}

		return false, false
	}

	return false, true
}

func (s *optimisticSyncSet) validate(pred, curr *syncNode) bool {
	for n := s.head; n.value <= pred.value; n = n.next {
		if n == pred {
			return pred.next == curr
		}
	}

	return false
}

// NewOptimisticSyncSet provides optimistic thread-safe set implementation with a mutex in every list node.
func NewOptimisticSyncSet() Set {
	// set must contain sentinel nodes with minimal and maximal values
	s := &optimisticSyncSet{}
	s.head = &syncNode{value: -math.MaxInt64}
	s.head.next = &syncNode{value: math.MaxInt64}

	return s
}
