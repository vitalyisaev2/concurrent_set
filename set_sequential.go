package set

import (
	"math"
)

type node struct {
	next  *node
	value int
}

var _ Set = (*sequentialSet)(nil)

// thread-unsafe implementation of linked-list based set.
type sequentialSet struct {
	head *node
}

func (s *sequentialSet) Insert(value int) bool {
	pred := s.head
	curr := pred.next

	for curr.value < value {
		pred = curr
		curr = pred.next
	}

	if curr.value == value {
		return false
	}

	newNode := &node{value: value, next: curr}
	pred.next = newNode

	return true
}

func (s sequentialSet) Contains(value int) bool {
	var (
		pred *node
		curr = s.head.next
	)

	for curr.value < value {
		pred = curr
		curr = pred.next
	}

	return curr.value == value
}

func (s *sequentialSet) Remove(value int) bool {
	pred := s.head
	curr := s.head.next

	for curr.value < value {
		pred = curr
		curr = pred.next
	}

	if curr.value == value {
		pred.next = curr.next

		return true
	}

	return false
}

// NewSequentialSet provides simple thread-unsafe implementation of linked list based set.
func NewSequentialSet() Set {
	// set must contain sentinel nodes with minimal and maximal values
	s := &sequentialSet{}
	s.head = &node{value: -math.MaxInt64}
	s.head.next = &node{value: math.MaxInt64}

	return s
}
