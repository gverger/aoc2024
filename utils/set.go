package utils

import (
	"fmt"
	"slices"
	"strings"
)

type Set[T comparable] map[T]struct{}

func NewSet[T comparable]() Set[T] {
	return make(Set[T])
}

func (s Set[T]) Add(value T) {
	s[value] = struct{}{}
}

func (s Set[T]) Exists(value T) bool {
	_, ok := s[value]
	return ok
}

func (s Set[T]) Delete(value T) {
	delete(s, value)
}

func (s Set[T]) Union(other Set[T]) {
	for v := range other {
		s.Add(v)
	}
}

func (s Set[T]) Intersection(other Set[T]) {
	for v := range s {
		if !other.Exists(v) {
			s.Delete(v)
		}
	}
}

func (s Set[T]) String() string {
	keys := make([]string, 0, len(s))
	for k := range s {
		keys = append(keys, fmt.Sprint(k))
	}
	slices.Sort(keys)
	return "{" + strings.Join(keys, ",") + "}"
}
