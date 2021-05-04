package set

/*
	FIXME: nonblocking implementation seems to be broken. Look at this code:
	https://github.com/gramoli/synchrobench/blob/master/java/src/linkedlists/lockfree/NonBlockingLinkedListSet.java#L43
	This line causes nil pointer dereference when trying to add first item into the empty set
	(with only two sentinel nodes).
*/

import (
	"math"
	"sync/atomic"
	"unsafe"
)

const mask uintptr = 1

type nonBlockingNode struct {
	next  *atomicMarkableReference
	value int
}

type atomicMarkableReference struct {
	value uintptr
}

func (amr *atomicMarkableReference) getNode() *nonBlockingNode {
	return (*nonBlockingNode)(unsafe.Pointer((atomic.LoadUintptr(&amr.value)) & ^mask))
}

func (amr *atomicMarkableReference) getMark() bool {
	current := atomic.LoadUintptr(&amr.value) & mask
	switch current {
	case 1:
		return true
	case 0:
		return false
	default:
		panic(current)
	}
}

func (amr *atomicMarkableReference) getBoth() (*nonBlockingNode, bool) {
	current := atomic.LoadUintptr(&amr.value)
	mark := current & mask

	return (*nonBlockingNode)(unsafe.Pointer(current & ^mask)), uintptrToBool(mark)
}

func (amr *atomicMarkableReference) compareAndSet(expectedNode, desiredNode *nonBlockingNode, expectedMark, desiredMark bool) bool {
	expected := amr.combine(expectedNode, expectedMark)
	desired := amr.combine(desiredNode, desiredMark)

	return atomic.CompareAndSwapUintptr(&amr.value, expected, desired)
}

func (amr *atomicMarkableReference) combine(node *nonBlockingNode, mark bool) uintptr {
	return (uintptr(unsafe.Pointer(node)) & ^mask) | boolToUintptr(mark)
}

func boolToUintptr(b bool) uintptr {
	if b {
		return 1
	}

	return 0
}

func uintptrToBool(val uintptr) bool {
	switch val {
	case 0:
		return false
	case 1:
		return true
	default:
		panic(val)
	}
}

func newAtomicMarkableReference(node *nonBlockingNode, mark bool) *atomicMarkableReference {
	amr := &atomicMarkableReference{}
	amr.value = amr.combine(node, mark)

	return amr
}

type window struct {
	pred, curr *nonBlockingNode
}

func findWindow(head *nonBlockingNode, val int) *window {
	var (
		pred, curr, succ *nonBlockingNode
		snip             bool
		marked           bool
	)

LOOP:
	for {
		pred = head
		curr = pred.next.getNode()
		for {
			// FIXME: you can't find succ if there are only two sentinel nodes in set
			succ, marked = curr.next.getBoth()
			for marked {
				snip = pred.next.compareAndSet(curr, succ, false, false)
				if !snip {
					continue LOOP
				}

				curr = succ
				succ, marked = curr.next.getBoth()
			}

			if curr.value > val {
				return &window{pred: pred, curr: curr}
			}

			pred = curr
			curr = succ
		}
	}
}

type nonBlockingSet struct {
	head *nonBlockingNode
}

func (s *nonBlockingSet) Insert(value int) bool {
	for {
		w := findWindow(s.head, value)
		pred := w.pred
		curr := w.curr

		if curr.value == value {
			return false
		}

		newNode := &nonBlockingNode{value: value}
		newNode.next = newAtomicMarkableReference(curr, false)

		if pred.next.compareAndSet(curr, newNode, false, false) {
			return true
		}
	}
}

func (s *nonBlockingSet) Contains(value int) bool {
	curr := s.head

	for curr.value < value {
		curr = curr.next.getNode()
	}

	return curr.value == value && !curr.next.getMark()
}

func (s *nonBlockingSet) Remove(value int) bool {
	for {
		w := findWindow(s.head, value)
		pred := w.pred
		curr := w.curr

		if curr.value != value {
			return false
		}

		succ := curr.next.getNode()
		snip := curr.next.compareAndSet(succ, succ, false, true)

		if !snip {
			continue
		}

		pred.next.compareAndSet(curr, succ, false, false)

		return true
	}
}

// NewNonBlockingSyncSet builds wait-free implementation of set.
func NewNonBlockingSyncSet() Set {
	s := &nonBlockingSet{}
	s.head = &nonBlockingNode{value: -math.MaxInt64}
	tail := &nonBlockingNode{value: math.MaxInt64}
	s.head.next = newAtomicMarkableReference(tail, false)

	return s
}
