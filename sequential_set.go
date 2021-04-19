package set

type node struct {
	value int
	next  *node
}

var _ Set = (*sequentialSet)(nil)

// thread-unsafe implementation of linked-list based set.
type sequentialSet struct {
	head *node
}

func (s *sequentialSet) Insert(value int) bool {
	// edge case when the set is empty
	if s.head == nil {
		s.head = &node{value: value}
		return true
	}

	// only one item in the set
	if s.head.next == nil {
		switch {
		case s.head.value < value:
			s.insertBetween(s.head, nil, value)
			return true
		case s.head.value > value:
			// swap head and new node
			oldHead := s.head
			s.head = &node{value: value, next: oldHead}
			return true
		default:
			return false
		}
	}

	// multiple items in set, seek the predecessor
	pred := s.head
	curr := s.head.next

	for curr.next != nil && curr.value < value {
		pred = curr
		curr = pred.next
	}

	if curr.value == value {
		return false
	}

	s.insertBetween(curr, nil, value)
	return true
}

func (s *sequentialSet) insertBetween(pred, next *node, value int) {
	inserted := &node{value: value}
	inserted.next = next
	pred.next = inserted
}

func (s sequentialSet) Contains(value int) bool {
	if s.head == nil {
		return false
	}

	if s.head.value == value {
		return true
	}

	pred := s.head
	curr := s.head.next

	for curr.value < value {
		pred = curr
		curr = pred.next
	}

	return curr.value == value
}

func (s *sequentialSet) Remove(value int) bool {
	if s.head == nil {
		return false
	}

	if s.head.value == value {
		s.head = nil
		return true
	}

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

func NewSequentialSet() Set {
	return &sequentialSet{}
}
