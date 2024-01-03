package operators

// If returns the first value if the condition is true, otherwise it returns the second value.
func If[T any](condition bool, then T, otherwise T) T {
	if condition {
		return then
	}

	return otherwise
}
