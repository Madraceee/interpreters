package utils

type Stack[T any] struct {
	Push    func(T)
	Pop     func() T
	Top     func() T
	Length  func() int
	IsEmpty func() bool
	Itr     func() (next func() T, has func() bool)
}

// Make Iterator better
// Current one seems shitty
func NewStack[T any]() Stack[T] {
	slice := make([]T, 0)
	return Stack[T]{
		Push: func(t T) {
			slice = append(slice, t)
		},
		Pop: func() T {
			top := slice[len(slice)-1]
			slice = slice[:len(slice)-1]
			return top
		},
		Top: func() T {
			return slice[len(slice)-1]
		},
		Length: func() int {
			return len(slice)
		},
		IsEmpty: func() bool {
			return len(slice) == 0
		},
		Itr: func() (next func() T, has func() bool) {
			i := 0
			return func() T {
					val := slice[i]
					i += 1
					return val
				},
				func() bool {
					return i < len(slice)
				}
		},
	}
}
