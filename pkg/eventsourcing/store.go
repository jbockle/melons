package eventsourcing

import (
	"context"
	"time"
)

type EventRecord interface {
	StreamId() string
	StreamVersion() int
	Event() *Event
	Metadata() map[string]any
	Timestamp() time.Time
}

type EventStream interface {
	Id() string
	Version() int
	Events() []EventRecord
}

type EventStore interface {
	Append(
		ctx context.Context,
		streamId string,
		expectedStreamVersion int,
		events ...*Event,
	) error

	Load(
		ctx context.Context,
		streamId string,
		fromVersion int,
	) (EventStream, error)
}
