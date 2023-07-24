package util

type Set[T any] struct {
	items []T
	cmp   func(T, T) bool
}

func NewSet[T any](cmp func(T, T) bool) Set[T] {
	return Set[T]{
		items: make([]T, 0, 32),
		cmp:   cmp,
	}
}

func (s Set[T]) Has(item T) bool {
	for i := range s.items {
		if s.cmp(s.items[i], item) {
			return true
		}
	}
	return false
}

func (s *Set[T]) Add(item T) bool {
	if s.Has(item) {
		return true
	}
	s.items = append(s.items, item)
	return false
}

func (s *Set[T]) Remove(item T) bool {
	for i := range s.items {
		if s.cmp(s.items[i], item) {
			s.items[i] = s.items[len(s.items)-1]
			s.items = s.items[:len(s.items)-1]
			return true
		}
	}
	return false
}

func (s *Set[T]) Intersect(rhs Set[T]) Set[T] {
	intersection := NewSet(s.cmp)
	for i := range s.items {
		if rhs.Has(s.items[i]) {
			intersection.Add(s.items[i])
		}
	}
	return intersection
}

func (s *Set[T]) Unity(rhs Set[T]) Set[T] {
	unity := NewSet(s.cmp)
	for i := range s.items {
		unity.Add(s.items[i])
	}
	for i := range rhs.items {
		unity.Add(rhs.items[i])
	}
	return unity
}

func (s *Set[T]) Difference(rhs Set[T]) Set[T] {
	difference := NewSet(s.cmp)
	for i := range s.items {
		if !rhs.Has(s.items[i]) {
			difference.Add(s.items[i])
		}
	}
	return difference
}

func (s Set[T]) IsEmpty() bool {
	return len(s.items) == 0
}

func (s Set[T]) Contains(rhs Set[T]) bool {
	for i := range rhs.items {
		if !s.Has(rhs.items[i]) {
			return false
		}
	}
	return true
}

func (s Set[T]) Items() []T {
	return s.items
}
