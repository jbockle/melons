package eventsourcing

import "fmt"

///
type ConcurrencyViolationError struct {
	StreamId        string
	ExpectedVersion int
	ActualVersion   int
}

func (err ConcurrencyViolationError) Error() string {
	return fmt.Sprintf("Concurrency violation: expected version %d, actual version %d", err.ExpectedVersion, err.ActualVersion)
}

///
type EventStreamNotFoundError struct {
	StreamId string
}

func (err EventStreamNotFoundError) Error() string {
	return fmt.Sprintf("Event stream '%s' not found", err.StreamId)
}

///
type EventStreamVersionNotFoundError struct {
	StreamId         string
	RequestedVersion int
	ActualVersion    int
}

func (err EventStreamVersionNotFoundError) Error() string {
	return fmt.Sprintf("Event stream '%s' version %d not found, actual version %d", err.StreamId, err.RequestedVersion, err.ActualVersion)
}
