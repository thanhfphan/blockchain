package set

import "golang.org/x/exp/maps"

const minSetSize = 16

type Set[T comparable] map[T]struct{}

func New[T comparable](initSize uint) Set[T] {
	return make(map[T]struct{}, initSize)
}

func (s *Set[T]) resize(size int) {
	if *s == nil {
		if minSetSize > size {
			size = minSetSize
		}
		*s = make(map[T]struct{}, size)
	}
}

func (s *Set[T]) Add(elts ...T) {
	s.resize(2 * len(elts))
	for _, item := range elts {
		(*s)[item] = struct{}{}
	}
}

func (s *Set[T]) Contains(elt T) bool {
	_, exist := (*s)[elt]
	return exist
}

func (s Set[_]) Len() int {
	return len(s)
}

func (s *Set[_]) Clear() {
	maps.Clear(*s)
}

func (s Set[T]) List() []T {
	return maps.Keys(s)
}
