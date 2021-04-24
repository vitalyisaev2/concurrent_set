package set

import (
	"math"
	"sync"
)

type lazySyncNode struct {
	next   *lazySyncNode
	value  int
	mutex  sync.Mutex
	marked bool
}

type lazySyncSet struct {
	head *lazySyncNode
}

func (s *lazySyncSet) Insert(value int) bool {
	pred := s.head
	curr := pred.next

	for curr.value < value {
		pred = curr
		curr = curr.next
	}

	pred.mutex.Lock()
	curr.mutex.Lock()

	defer func() {
		curr.mutex.Unlock()
		pred.mutex.Unlock()
	}()

	if s.validate(pred, curr) {
		if curr.value == value {
			return false
		} else {
			newNode := &lazySyncNode{value: value, next: curr}
			pred.next = newNode
		}
	}

	return true
}

func (s *lazySyncSet) Contains(value int) bool {
	curr := s.head

	for curr.value < value {
		curr = curr.next
	}

	return curr.value == value && !curr.marked
}

func (s *lazySyncSet) Remove(value int) bool {
	pred := s.head
	curr := s.head.next

	for curr.value < value {
		pred = curr
		curr = curr.next
	}

	pred.mutex.Lock()
	curr.mutex.Lock()

	defer func() {
		curr.mutex.Unlock()
		pred.mutex.Unlock()
	}()

	if s.validate(pred, curr) {
		if curr.value == value {
			curr.marked = true
			pred.next = curr.next

			return true
		}
	}

	return false
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
