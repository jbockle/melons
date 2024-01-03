package ioc

import (
	"reflect"
)

type transientResolver struct {
	descriptors []*ServiceDescriptor
	last        *ServiceDescriptor
}

func (r *transientResolver) With(descriptor *ServiceDescriptor) {
	r.descriptors = append(r.descriptors, descriptor)
	r.last = descriptor
}

func (r *transientResolver) Resolve() any {
	switch r.last.InstanceOrFactory.Kind() {
	case reflect.Func:
		return r.last.InstanceOrFactory.Call(nil)[0].Interface()
	default:
		panic("not implemented")
	}
}

func (r *transientResolver) ResolveAll() []any {
	services := make([]any, 0)
	for _, descriptor := range r.descriptors {
		switch descriptor.InstanceOrFactory.Kind() {
		case reflect.Func:
			services = append(services, descriptor.InstanceOrFactory.Call(nil)[0].Interface())
		default:
			panic("not implemented")
		}
	}

	return services
}
