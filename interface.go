// Package set provides various implementations of linked-list based sets
package set

// Set contains unique integer values
// TODO: generics.
type Set interface {
	Insert(value int) bool
	Contains(value int) bool
	Remove(value int) bool
}
