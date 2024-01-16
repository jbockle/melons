package operators

import (
	"melons/pkg/reflectx"
	"reflect"
)

func AssignIfNilOrZero[T comparable](target *T, value T) {
	if reflectx.IsNilable[T]() {
		if reflect.ValueOf(*target).IsNil() {
			target = &value
		}

		return
	}

	var zeroValue T

	if *target == zeroValue {
		target = &value
	}
}
