package reflectx

import (
	"reflect"
)

func TypeOfGeneric[T any]() reflect.Type {
	var zero [0]T
	tt := reflect.TypeOf(zero).Elem()
	return tt
}
