package main

import "sort"

type sorterBase[T interface{}] struct {
	array   []T
	compare func(a, b T) bool
}

/// Used to sort any named object by its human-readable name.
type named interface {
	Name() string
}

/// Used to sort any prioritized object by its priority.
type prioritized interface {
	Priority() int32
}

/// Sorter interface method - returns length of the sorted array.
func (n *sorterBase[_]) Len() int {
	return len(n.array)
}

/// Sorter interface method - exchanges two items.
func (n *sorterBase[_]) Swap(i, j int) {
	tmp := n.array[i]
	n.array[i] = n.array[j]
	n.array[j] = tmp
}

/// Sorter interface method - compares two items.
func (n *sorterBase[_]) Less(i, j int) bool {
	return n.compare(n.array[i], n.array[j])
}

/// Given an array of items, order the elements by their Name().
func SortInPlaceByName[T named](data []T) {
	sort.Sort(&sorterBase[T]{
		array: data,
		compare: func(a, b T) bool {
			return a.Name() < b.Name()
		}})
}

/// Given an array of items, order the elements by their Priority().
func SortInPlaceByPriority[T prioritized](data []T) {
	sort.Sort(&sorterBase[T]{
		array: data,
		compare: func(a, b T) bool {
			return a.Priority() < b.Priority()
		}})
}
