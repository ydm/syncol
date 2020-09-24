package syncol

// Comparator should return true if lhs is less than rhs.
type Comparator func(lhs, rhs interface{}) bool
