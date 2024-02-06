package set

import (
	"encoding/json"
	"fmt"
)

type Set[T comparable] struct {
	value map[T]struct{}
}

func New[T comparable]() *Set[T] {
	return &Set[T]{value: make(map[T]struct{})}
}

func NewFrom[T comparable](items ...T) *Set[T] {
	s := New[T]()
	for _, item := range items {
		s.Add(item)
	}
	return s
}

func NewFromSlice[T comparable](items []T) *Set[T] {
	s := New[T]()
	for _, item := range items {
		s.Add(item)
	}
	return s
}

func (s *Set[T]) init() {
	if s.value == nil {
		s.value = make(map[T]struct{})
	}
}

func (s *Set[T]) Add(item T) bool {
	s.init()
	_, found := s.value[item]
	if !found {
		s.value[item] = struct{}{}
		return true
	}
	return false
}

func (s *Set[T]) AddAll(items ...T) {
	s.init()
	for _, item := range items {
		s.Add(item)
	}
}

func (s *Set[T]) Remove(item T) bool {
	s.init()
	_, found := s.value[item]
	if found {
		delete(s.value, item)
	}
	return found
}

func (s *Set[T]) Contains(item T) bool {
	s.init()
	_, found := s.value[item]
	return found
}

func (s *Set[T]) Size() int {
	s.init()
	return len(s.value)
}

func (s *Set[T]) ToSlice() []T {
	s.init()
	var result []T
	for k := range s.value {
		result = append(result, k)
	}

	return result
}

func (s *Set[T]) Clear() {
	s.value = make(map[T]struct{})
}

func (s *Set[T]) Union(other *Set[T]) *Set[T] {
	s.init()
	result := New[T]()
	for k := range s.value {
		result.Add(k)
	}
	for k := range other.value {
		result.Add(k)
	}
	return result
}

func (s *Set[T]) Intersection(other *Set[T]) *Set[T] {
	s.init()
	result := New[T]()
	for k := range s.value {
		if other.Contains(k) {
			result.Add(k)
		}
	}
	return result
}

func (s *Set[T]) Difference(other *Set[T]) *Set[T] {
	s.init()
	result := New[T]()
	for k := range s.value {
		if !other.Contains(k) {
			result.Add(k)
		}
	}
	return result
}

func (s *Set[T]) IsSubset(other *Set[T]) bool {
	s.init()
	for k := range s.value {
		if !other.Contains(k) {
			return false
		}
	}
	return true
}

func (s *Set[T]) IsSuperset(other *Set[T]) bool {
	s.init()
	return other.IsSubset(s)
}

func (s *Set[T]) Equal(other *Set[T]) bool {
	s.init()
	return s.IsSubset(other) && s.IsSuperset(other)
}

func (s *Set[T]) Clone() *Set[T] {
	s.init()
	result := New[T]()
	for k := range s.value {
		result.Add(k)
	}
	return result
}

func (s *Set[T]) String() string {
	s.init()
	return fmt.Sprintf("Set[%d]", s.Size())
}

func (s *Set[T]) IsEmpty() bool {
	s.init()
	return s.Size() == 0
}

func (s *Set[T]) MarshalJSON() ([]byte, error) {
	s.init()
	return json.Marshal(s.ToSlice())
}

func (s *Set[T]) UnmarshalJSON(data []byte) error {
	s.init()
	var slice []T

	if s.value == nil {
		s.value = make(map[T]struct{})
	}

	if err := json.Unmarshal(data, &slice); err != nil {
		return err
	}

	for _, item := range slice {
		s.Add(item)
	}

	return nil
}
