package sort

import "sort"

/// Used to sort any named object by its human-readable name.
type named interface {
	Name() string
}

/// Used to sort any prioritized object by its priority.
type prioritized interface {
	Priority() int32
}

/// Given an array of items, order the elements by their Name().
func SortInPlaceByName[T named](data []T) {
	sort.Slice(data, func(a, b int) bool {
		return data[a].Name() < data[b].Name()
	})
}

/// Given an array of items, order the elements by their Priority().
func SortInPlaceByPriority[T prioritized](data []T) {
	sort.Slice(data, func(a, b int) bool {
		return data[a].Priority() < data[b].Priority()
	})
}
