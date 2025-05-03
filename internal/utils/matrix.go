package utils

type Matrix[T any] struct {
	storage    [][]T
	cols       int
	actualLine int
}

func NewMatrix[T any](cols int) Matrix[T] {
	s := make([][]T, 1)
	s[0] = make([]T, 0, cols)

	return Matrix[T]{
		storage:    s,
		cols:       cols,
		actualLine: 0,
	}
}

func (m *Matrix[T]) Add(v T) {
	l := len(m.storage[m.actualLine])

	if l+1 > m.cols {
		m.storage = append(m.storage, make([]T, 0, m.cols))
		m.actualLine++
	}

	m.storage[m.actualLine] = append(m.storage[m.actualLine], v)
}

func (m *Matrix[T]) GetMatrix() [][]T {
	return m.storage
}
