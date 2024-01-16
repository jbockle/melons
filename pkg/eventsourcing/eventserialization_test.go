package eventsourcing

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
)

type eventInfo struct {
	EventInfo EventDescriptor `json:"eventInfo"`
}

type eventPayload struct {
	Payload Event `json:"payload"`
}

type MyEventRecord struct {
	eventInfo
	eventPayload
}

type CreatedEvent struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

func (e *CreatedEvent) EventDescriptor() EventDescriptor {
	return RegisterEventDescriptor[*CreatedEvent]("foo", 2)
}

func (er *MyEventRecord) UnmarshalJSON(data []byte) error {
	er.eventInfo = eventInfo{}
	if err := json.Unmarshal(data, &er.eventInfo); err != nil {
		return err
	}

	eventType, ok := eventTypes[er.EventInfo]
	if !ok {
		return fmt.Errorf("Event type not found %+v", er.eventInfo)
	}

	er.eventPayload = eventPayload{}
	er.eventPayload.Payload = reflect.New(eventType.Elem()).Interface().(Event)

	return json.Unmarshal(data, &er.eventPayload)
}

func TestSerializeEvent(t *testing.T) {
	e := &CreatedEvent{
		Id:   "123",
		Name: "foo",
	}

	er := &MyEventRecord{
		eventInfo: eventInfo{
			EventInfo: e.EventDescriptor(),
		},
		eventPayload: eventPayload{
			Payload: e,
		},
	}

	data, err := json.Marshal(er)
	if err != nil {
		t.Error(err)
	}

	er2 := &MyEventRecord{}
	if err := json.Unmarshal(data, er2); err != nil {
		t.Error(err)
	}

	t.Logf("%+v", er2)

	if er.eventInfo != er2.eventInfo {
		t.Error("Event descriptors are not equal")
	}

	expectedPayload := er.eventPayload.Payload.(*CreatedEvent)
	actualPayload := er2.eventPayload.Payload.(*CreatedEvent)

	if *expectedPayload != *actualPayload {
		t.Error("Event payloads are not equal")
	}
}
