package ioc

import (
	"fmt"
	"melons/pkg/reflectx"
	"reflect"
)

type Resolver interface {
	Resolve() any
	ResolveAll() []any
	With(descriptor *ServiceDescriptor)
}

var built bool = false
var serviceDescriptors []*ServiceDescriptor = make([]*ServiceDescriptor, 0)
var serviceResolvers = make(map[reflect.Type]Resolver)

type ResolvePanicked struct {
	Reason any
}

func (rp ResolvePanicked) Error() string {
	return fmt.Sprintf("Resolve() panicked: %v", rp.Reason)
}

func addDescriptor(descriptor *ServiceDescriptor) {
	if built {
		panic("Cannot add service descriptor after build")
	}

	serviceDescriptors = append(serviceDescriptors, descriptor)
}

func RegisterSingletonInstance[TService any](instance TService) error {
	serviceType := reflectx.TypeOfGeneric[TService]()

	switch serviceType.Kind() {
	case reflect.Interface, reflect.Struct:
		break
	default:
		return fmt.Errorf("service type %s must be interface or struct", serviceType.Name())
	}

	serviceValue := reflect.ValueOf(instance)

	addDescriptor(&ServiceDescriptor{
		Type:              serviceType,
		InstanceOrFactory: serviceValue,
		Lifetime:          Singleton,
		isInstance:        true,
	})

	return nil
}

func RegisterFactory[TService any](factory func() TService, lifetime ServiceLifetime) error {
	serviceType := reflectx.TypeOfGeneric[TService]()

	switch serviceType.Kind() {
	case reflect.Interface, reflect.Struct:
		break
	default:
		return fmt.Errorf("service type %s must be interface or struct", serviceType.Name())
	}

	serviceValue := reflect.ValueOf(factory)

	addDescriptor(&ServiceDescriptor{
		Type:              serviceType,
		InstanceOrFactory: serviceValue,
		Lifetime:          lifetime,
		isFactory:         true,
	})

	return nil
}

func Build() {
	for _, descriptor := range serviceDescriptors {
		resolver, ok := serviceResolvers[descriptor.Type]
		if !ok {
			switch descriptor.Lifetime {
			case Singleton:
				resolver = &singletonResolver{
					values: make(map[*ServiceDescriptor]any),
				}
			case Transient:
				resolver = &transientResolver{
					descriptors: make([]*ServiceDescriptor, 0),
				}
			}
			serviceResolvers[descriptor.Type] = resolver
		}

		resolver.With(descriptor)
	}
	built = true
}

func reset() {
	serviceDescriptors = make([]*ServiceDescriptor, 0)
	serviceResolvers = make(map[reflect.Type]Resolver)
	built = false
}

func IsRegistered(serviceType reflect.Type) bool {
	_, ok := serviceResolvers[serviceType]
	return ok
}

func Resolve[TService any]() TService {
	serviceType := reflectx.TypeOfGeneric[TService]()
	resolver, ok := serviceResolvers[serviceType]
	if !ok {
		panic(fmt.Errorf("service %s is not registered", serviceType.Name()))
	}

	defer func() {
		if r := recover(); r != nil {
			panic(ResolvePanicked{Reason: r})
		}
	}()

	return resolver.Resolve().(TService)
}

func ResolveAll[TService any]() []TService {
	serviceType := reflectx.TypeOfGeneric[TService]()
	resolver, ok := serviceResolvers[serviceType]
	if !ok {
		panic(fmt.Errorf("service %s is not registered", serviceType.Name()))
	}

	defer func() {
		if r := recover(); r != nil {
			panic(ResolvePanicked{Reason: r})
		}
	}()

	services := make([]TService, 0)

	for _, value := range resolver.ResolveAll() {
		services = append(services, value.(TService))
	}

	return services
}
