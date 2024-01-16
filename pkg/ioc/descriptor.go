package ioc

import (
	"reflect"
)

type ServiceDescriptor struct {
	Lifetime          ServiceLifetime
	Type              reflect.Type
	InstanceOrFactory reflect.Value
	isInstance        bool
	isFactory         bool
}

func (sd *ServiceDescriptor) IsInstance() bool {
	return sd.isInstance
}

func (sd *ServiceDescriptor) IsFactory() bool {
	return sd.isFactory
}

func (sd *ServiceDescriptor) IsSingleton() bool {
	return sd.Lifetime == Singleton
}

func (sd *ServiceDescriptor) IsTransient() bool {
	return sd.Lifetime == Transient
}
