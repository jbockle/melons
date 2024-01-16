package operators

// returns the first non-nil/zero value in a list of values.
//
// If the value is a zero value, it will be considered nil.
//
// If all values are zero values, the zero value is returned
func Coalesce[T comparable](values ...T) T {
	var zeroValue T
	for _, value := range values {
		if value != zeroValue {
			return value
		}
	}

	return zeroValue
}
