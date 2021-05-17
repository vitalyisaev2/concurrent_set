package set

import (
	"math"
	"sync/atomic"
	"unsafe"
)

type nonBlockingNode struct {
	next  *atomicMarkableReference
	value int
}

type markableReference struct {
	node *nonBlockingNode
	mark bool
}

type atomicMarkableReference struct {
	ref unsafe.Pointer // *markableReference
}

func (amr *atomicMarkableReference) getNode() *nonBlockingNode {
	if amr == nil {
		return nil
	}

	existingRef := (*markableReference)(atomic.LoadPointer(&amr.ref))

	return existingRef.node
}

func (amr *atomicMarkableReference) getMark() bool {
	if amr == nil {
		return false
	}

	existingRef := (*markableReference)(atomic.LoadPointer(&amr.ref))

	return existingRef.mark
}

func (amr *atomicMarkableReference) getBoth() (*nonBlockingNode, bool) {
	if amr == nil {
		return nil, false
	}

	existingRef := (*markableReference)(atomic.LoadPointer(&amr.ref))

	return existingRef.node, existingRef.mark
}

func (amr *atomicMarkableReference) compareAndSet(expectedNode, desiredNode *nonBlockingNode, expectedMark, desiredMark bool) bool {
	if amr == nil {
		return false
	}

	existingRefValue := atomic.LoadPointer(&amr.ref)
	existingRef := (*markableReference)(existingRefValue)

	newRef := &markableReference{node: desiredNode, mark: desiredMark}
	newRefValue := unsafe.Pointer(newRef)

	return existingRef.node == expectedNode &&
		existingRef.mark == expectedMark &&
		atomic.CompareAndSwapPointer(&amr.ref, existingRefValue, newRefValue)
}

func newAtomicMarkableReference(node *nonBlockingNode, mark bool) *atomicMarkableReference {
	ref := &markableReference{node: node, mark: mark}
	return &atomicMarkableReference{ref: unsafe.Pointer(ref)}
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
			succ, marked = curr.next.getBoth()
			for marked {
				snip = pred.next.compareAndSet(curr, succ, false, false)
				if !snip {
					continue LOOP
				}

				curr = succ
				succ, marked = curr.next.getBoth()
			}

			if curr.value >= val {
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

	head := &nonBlockingNode{value: math.MinInt64}
	tail := &nonBlockingNode{value: math.MaxInt64}

	head.next = newAtomicMarkableReference(tail, false)
	tail.next = newAtomicMarkableReference(nil, false)

	s.head = head

	return s
}
