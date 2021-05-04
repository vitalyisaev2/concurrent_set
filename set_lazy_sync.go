package set

import (
	"math"
	"sync"
)

type lazySyncNode struct {
	next  *lazySyncNode
	value int
	sync.Mutex
	marked bool
}

type lazySyncSet struct {
	head *lazySyncNode
}

func (s *lazySyncSet) Insert(value int) bool {
	for {
		result, repeat := s.insertLoopBody(value)
		if !repeat {
			return result
		}
	}
}

//nolint:dupl // it's better to copy-paste code than messing with inheritance
func (s *lazySyncSet) insertLoopBody(value int) (result, repeat bool) {
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

		newNode := &lazySyncNode{value: value, next: curr}
		pred.next = newNode

		return true, false
	}

	return false, true
}

func (s *lazySyncSet) Contains(value int) bool {
	for {
		result, repeat := s.containsLoopBody(value)
		if !repeat {
			return result
		}
	}
}

func (s *lazySyncSet) containsLoopBody(value int) (result, repeat bool) {
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

func (s *lazySyncSet) Remove(value int) bool {
	for {
		result, repeat := s.removeLoopBody(value)
		if !repeat {
			return result
		}
	}
}

func (s *lazySyncSet) removeLoopBody(value int) (result, repeat bool) {
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
			curr.marked = true
			pred.next = curr.next

			return true, false
		}

		return false, false
	}

	return false, true
}

func (s *lazySyncSet) validate(pred, curr *lazySyncNode) bool {
	return !pred.marked && !curr.marked && pred.next == curr
}

// NewLazySyncSet provides lazy thread-safe set implementation with a mutex in every list node.
func NewLazySyncSet() Set {
	// set must contain sentinel nodes with minimal and maximal values
	s := &lazySyncSet{}
	s.head = &lazySyncNode{value: -math.MaxInt64}
	s.head.next = &lazySyncNode{value: math.MaxInt64}

	return s
}
