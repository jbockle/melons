package ioc

import (
	"fmt"
	"reflect"
)

type singletonResolver struct {
	values map[*ServiceDescriptor]any
	last   *ServiceDescriptor
}

func (r *singletonResolver) With(descriptor *ServiceDescriptor) {
	r.values[descriptor] = nil
	r.last = descriptor
}

func (r *singletonResolver) Resolve() any {
	r.resolve(r.last)

	return r.values[r.last]
}

func (r *singletonResolver) ResolveAll() []any {
	services := make([]any, 0)
	for descriptor, value := range r.values {
		if value == nil {
			r.resolve(descriptor)
			value = r.values[descriptor]
		}

		services = append(services, value)
	}

	return services
}

func (r *singletonResolver) resolve(descriptor *ServiceDescriptor) {
	if r.values[descriptor] == nil {
		kind := descriptor.InstanceOrFactory.Kind()
		fmt.Println(fmt.Sprintf("***kind: %v", kind.String()))
		switch kind {
		case reflect.Func:
			r.values[descriptor] = descriptor.InstanceOrFactory.Call(nil)[0].Interface()
		case reflect.Struct, reflect.Ptr:
			r.values[descriptor] = descriptor.InstanceOrFactory.Interface()
		default:
			panic("not implemented")
		}
	}
}
