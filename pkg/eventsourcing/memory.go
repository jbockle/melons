package eventsourcing

import (
	"context"
	"time"
)

// /
type MemoryEventRecord struct {
	EventRecord
	streamId      string
	streamVersion int
	event         *Event
	metadata      map[string]any
	timestamp     time.Time
}

func (record *MemoryEventRecord) StreamId() string {
	return record.streamId
}

func (record *MemoryEventRecord) StreamVersion() int {
	return record.streamVersion
}

func (record *MemoryEventRecord) Event() *Event {
	return record.event
}

func (record *MemoryEventRecord) Metadata() map[string]any {
	return record.metadata
}

func (record *MemoryEventRecord) Timestamp() time.Time {
	return record.timestamp
}

// /
type MemoryEventStream struct {
	EventStream
	id      string
	version int
	events  []MemoryEventRecord
}

func (stream *MemoryEventStream) Id() string {
	return stream.id
}

func (stream *MemoryEventStream) Version() int {
	return stream.version
}

func (stream *MemoryEventStream) Events() []EventRecord {
	records := make([]EventRecord, len(stream.events))
	for i, event := range stream.events {
		records[i] = &event
	}
	return records
}

func (stream *MemoryEventStream) fromVersion(fromVersion int) *MemoryEventStream {
	records := make([]MemoryEventRecord, 0)
	for _, record := range stream.events {
		if record.StreamVersion() >= fromVersion {
			records = append(records, record)
		}
	}
	return &MemoryEventStream{
		id:      stream.id,
		version: stream.version,
		events:  records,
	}
}

// /
type MemoryEventStore struct {
	streams map[string]*MemoryEventStream
}

func NewMemoryEventStore() *MemoryEventStore {
	return &MemoryEventStore{
		streams: make(map[string]*MemoryEventStream),
	}
}

func (store *MemoryEventStore) Append(
	ctx context.Context,
	streamId string,
	expectedStreamVersion int,
	events ...*Event,
) error {
	stream, ok := store.streams[streamId]
	if !ok {
		stream = &MemoryEventStream{
			id:      streamId,
			version: 0,
			events:  make([]MemoryEventRecord, 0),
		}

		store.streams[streamId] = stream
	}

	if stream.version != expectedStreamVersion {
		return ConcurrencyViolationError{
			StreamId:        streamId,
			ExpectedVersion: expectedStreamVersion,
			ActualVersion:   stream.version,
		}
	}

	for _, event := range events {
		record := MemoryEventRecord{
			streamId:      streamId,
			streamVersion: stream.version + 1,
			event:         event,
			metadata:      GetEventMetadata(ctx),
			timestamp:     time.Now(),
		}
		stream.events = append(stream.events, record)
		stream.version++
	}

	return nil
}

func (store *MemoryEventStore) Load(
	ctx context.Context,
	streamId string,
	fromVersion int,
) (EventStream, error) {
	stream, ok := store.streams[streamId]
	if !ok {
		return nil, EventStreamNotFoundError{
			StreamId: streamId,
		}
	}

	if fromVersion > stream.version {
		return nil, EventStreamVersionNotFoundError{
			StreamId:         streamId,
			RequestedVersion: fromVersion,
			ActualVersion:    stream.version,
		}
	}

	return stream.fromVersion(fromVersion), nil
}
