package syncol

// Comparator should return true if lhs is less than rhs.
type Comparator[T any] func(lhs, rhs T) bool
