package container

type Stack[E any] interface {
	Push(E)
	Pop() (e E, ok bool)
}

type Queue[E any] interface {
	Enqueue(E)
	Dequeue() (e E, ok bool)
}

type Countable interface {
	Len() int
}
