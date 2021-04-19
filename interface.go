package set

// TODO: generics.
type Set interface {
	Insert(value int) bool
	Contains(value int) bool
	Remove(value int) bool
}
