package datastruct

const default_init_size = 13

type Stack[T any] struct {
	data []T
	pos  int

	empty T
}

func NewStack[T any]() *Stack[T] {
	return NewStackSize[T](default_init_size)
}

func NewStackSize[T any](initSize int) *Stack[T] {
	if initSize <= 0 {
		initSize = default_init_size
	}
	s := &Stack[T]{data: make([]T, initSize), pos: -1}
	return s
}

func (s *Stack[T]) resize() {
	l := len(s.data)
	newData := make([]T, l*2)
	copy(newData, s.data)
	s.data = newData
}

func (s *Stack[T]) Cap() int {
	return len(s.data) - s.pos - 1
}

func (s *Stack[T]) Push(e T) {
	if s.Cap() <= 0 {
		s.resize()
	}
	s.pos++
	s.data[s.pos] = e
}

func (s *Stack[T]) Pop() (e T) {
	if s.pos >= 0 {
		e = s.data[s.pos]
		s.data[s.pos] = s.empty
		s.pos--
	}
	return
}

func (s *Stack[T]) IsEmpty() bool {
	return s.pos < 0
}

func (s *Stack[T]) Size() int {
	return s.pos + 1
}

func (s *Stack[T]) Peek() T {
	if !s.IsEmpty() {
		return s.data[0]
	}
	return s.empty
}
