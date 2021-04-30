package set

import (
	"sync/atomic"
	"unsafe"
)

const mask uintptr = 1

type nodeNonBlocking struct {
	val  int
	next *nodeNonBlocking
}

type atomicMarkableReference struct {
	value uintptr
}

func (amr *atomicMarkableReference) getNode() *nodeNonBlocking {
	return (*nodeNonBlocking)(unsafe.Pointer((atomic.LoadUintptr(&amr.value)) & ^mask))
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

func (amr *atomicMarkableReference) getBoth() (*nodeNonBlocking, bool) {
	current := atomic.LoadUintptr(&amr.value)
	mark := current & mask
	return (*nodeNonBlocking)(unsafe.Pointer(current & ^mask)), uintptrToBool(mark)
}

func (amr *atomicMarkableReference) compareAndSet(expectedNode, desiredNode *nodeNonBlocking, expectedMark, desiredMark bool) bool {
	expected := amr.combine(expectedNode, expectedMark)
	desired := amr.combine(desiredNode, desiredMark)
	return atomic.CompareAndSwapUintptr(&amr.value, expected, desired)
}

func (amr *atomicMarkableReference) combine(node *nodeNonBlocking, mark bool) uintptr {
	return (uintptr(unsafe.Pointer(node)) & ^mask) | boolToUintptr(mark)
}

func newAtomicMarkableReference(node *nodeNonBlocking, mark bool) *atomicMarkableReference {
	amr := &atomicMarkableReference{}
	amr.value = amr.combine(node, mark)

	return amr
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
