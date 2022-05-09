package types

// ListElement is the type for the individual elements in the Queue/Stack versions
// based on a plain list
type ListElement[E any] struct {
	Value E
	Next  *ListElement[E]
}
