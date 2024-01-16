package reflectx

import (
	"reflect"
)

func IsNilable[T comparable]() bool {
	tt := TypeOfGeneric[T]()

	switch tt.Kind() {
	case reflect.Chan, reflect.Func, reflect.Map, reflect.Pointer, reflect.UnsafePointer, reflect.Interface, reflect.Slice:
		return true
	default:
		return false
	}
}
