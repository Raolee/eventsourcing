package eventsourcing

import (
	"fmt"
	"github.com/pkg/errors"
	"strings"
)

type Code int

const (
	Nothing Code = iota
	AlreadyLockedEvent
	AlreadyUnlockedEvent
	NoHasCommand
	CommandError
	ValidateError
	DispenseEventNoError
	EventStorageError
	SnapshotStorageError
)

type EventSourceError struct {
	Code   Code
	err    error
	format string
	args   []any
}

func newEventSourceError(code Code, err error, format string, args ...any) *EventSourceError {
	return &EventSourceError{
		Code:   code,
		err:    err,
		format: format,
		args:   args,
	}
}

func (e *EventSourceError) Error() string {
	var returnError error
	switch {
	case len(strings.TrimSpace(e.format)) == 0:
		if e.err != nil {
			returnError = e.err
		} else {
			returnError = errors.New("error info is nil")
		}
	default:
		msg := fmt.Sprintf(e.format, e.args...)
		if e.err != nil {
			returnError = errors.Wrap(e.err, msg) // 에러가 struct 에 포함되어 있으면 wrapping
		} else {
			returnError = errors.New(msg)
		}
	}
	return returnError.Error()
}

func NewLockedEventError(err error, pk PartitionKey, et *EventType) error {
	return newEventSourceError(AlreadyLockedEvent, err, "already locked. pk(%s), eventType(%s)", pk, et.String())
}

func NewUnlockedEventError(err error, pk PartitionKey, et *EventType) error {
	return newEventSourceError(AlreadyUnlockedEvent, err, "already unlocked. pk(%s), eventType(%s)", pk, et.String())
}

func NewNoHasCommandError(pk PartitionKey, et *EventType) error {
	return newEventSourceError(NoHasCommand, nil, "no has command. pk(%s), eventType(%s)", pk, et.String())
}

func NewCommandError[R any](err error, pk PartitionKey, e *Event[R]) error {
	return newEventSourceError(CommandError, err, "occur error command. pk(%s), event(%s)", pk, e.String())
}

func NewValidateError[S CommonState[R], R any](err error, pk PartitionKey, s *State[S, R]) error {
	return newEventSourceError(ValidateError, err, "occur error command. pk(%s), state(%s)", pk, s.String())
}

func NewDispenseEventNoError(err error, pk PartitionKey) error {
	return newEventSourceError(DispenseEventNoError, err, "occur error dispense eventNo. pk(%s)", pk)
}

func NewEventStorageError(err error) error {
	return newEventSourceError(EventStorageError, err, "")
}

func NewSnapshotStorageError(err error) error {
	return newEventSourceError(SnapshotStorageError, err, "")
}
