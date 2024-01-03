package eventsourcing

import (
	"melons/pkg/operators"
	"melons/pkg/reflectx"
	"reflect"
	"slices"
)

type Event interface {
	EventDescriptor() EventDescriptor
}

type EventDescriptor struct {
	Name    string `json:"name"`
	Version int    `json:"version"`
}

var eventDescriptors = make(map[reflect.Type]EventDescriptor)
var eventTypes = make(map[EventDescriptor]reflect.Type)

func RegisterEventDescriptor[TEvent Event](values ...any) EventDescriptor {
	eventType := reflectx.TypeOfGeneric[TEvent]()

	nameIndex := slices.IndexFunc(values, func(value any) bool {
		return reflect.ValueOf(value).Kind() == reflect.String
	})
	versionIndex := slices.IndexFunc(values, func(value any) bool {
		return reflect.ValueOf(value).Kind() == reflect.Int
	})

	descriptor := eventDescriptors[eventType]

	descriptor = EventDescriptor{
		Name:    operators.If(nameIndex >= 0, values[nameIndex].(string), eventType.Name()),
		Version: operators.If(versionIndex >= 0, values[versionIndex].(int), 1),
	}

	eventDescriptors[eventType] = descriptor
	eventTypes[descriptor] = eventType

	return *&descriptor
}
