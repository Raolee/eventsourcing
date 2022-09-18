package eventsourcing

import (
	"fmt"
	"github.com/pkg/errors"
	"strings"
)

// Event Sourcing 에서 다루는 에러를 정의한다.
//
// [AlreadyLockedEvent]
// - 이벤트를 잠금 처리하는데 이미 잠겨저 있을 때 사용하는 에러
// 잠금이 필요한 이벤트는 순서나 동시처리가 불가능한 매우 중요한 이벤트일 것임.
// 따라서 이벤트의 잠금이 이미 되어있는 경우는 에러로 취급하고 크리티컬하게 다루어져야 함.
//
// [AlreadyUnlockedEvent]
// - 이벤트는 잠금 해제하려는데 이미 잠금이 해제되어 있을 때 사용하는 에러
// 이벤트가 잠금되었다가, 잠금이 해제되는 것은 '잠금'->'잠금 해제' 의 짝(pair)가 한 쌍만 있어야 함.
// 이미 해제가 되어 있다면 로직상 '잠금 해제'의 한 쌍을 벗어나는 처리가 이루어졌다는 것이므로 에러로 다뤄야함.
//
// [NoHasCommand]
// - 이벤트에 매핑된 커맨드가 없는 경우 사용하는 에러
// 이벤트는 반드시 Command 와 짝이 이루어져야함.
//
// [CommandError]
// - Command 처리 중 에러가 발생하는 경우 사용하는 에러
// Command 는 수행이 State 변경이 go code 상에서만 이뤄지므로 에러가 발생하지 않게 만들어야함.
// 하지만 코드 상 bug 로 에러가 발생한다면, CommandError 를 사용하고 이는 매우 크리티컬하게 다뤄야 함
//
// [ValidateError]
// - Validate 중 에러가 발생하는 경우 사용하는 에러
// Validate 는 storage 에서 가져온 state 와 event 를 go code 상에서만 다루므로 에러가 발생하지 않게 만들어야함.
// CommandError 와 마찬가지로 코드 상 bug 가 있어 에러가 발생한다면, ValidateError 를 사용하고 매우 크리티컬하게 다뤄야 함.
//
// [DispenseEventNoError]
// - Event No 를 발급받다가 에러가 발생할 때 사용하는 에러
// Event No 발급은 Storage 의 특성을 따라가므로, 실패하거나 에러가 발생할 때는 Storage 상의 이슈를 만나게 됨.
// 이 경우, 다른 Event No 발급도 실패할 가능성이 있으므로 크리티컬하게 모니터링하고 다뤄야 함.
//
// [EventStorageError]
// - Event Storage 의 에러가 발생하는 경우 사용하는 에러
//
// [SnapshotStorageError]
// - Snapshot Storage 의 에러가 발생하는 경우 사용하는 에러

// Code | 이벤트 소싱에서 다루는 에러 케이스의 코드 정의
type Code int

const (
	Nothing Code = iota
	AlreadyLockedEvent
	AlreadyUnlockedEvent
	NoHasLockError
	NoHasCommand
	CommandError
	ValidateError
	DispenseEventNoError
	EventStorageError
	SnapshotStorageError
)

// EventSourceError | 이벤트 소싱에서 다루는 에러를 wrapping 한 구조체
type EventSourceError struct {
	Code   Code   // 에러 코드
	err    error  // nullable, 에러 정보
	format string // string 변환 시 사용할 포맷
	args   []any  // string 변환 시 사용할 args
}

func newEventSourceError(code Code, err error, format string, args ...any) *EventSourceError {
	return &EventSourceError{
		Code:   code,
		err:    err,
		format: format,
		args:   args,
	}
}

// error interface 를 구현
func (e *EventSourceError) Error() string {
	var returnError error
	switch {
	case len(strings.TrimSpace(e.format)) == 0: // 포맷이 없는 경우
		if e.err != nil {
			returnError = e.err // 에러가 있으면 에러 그 자체를 찍음
		} else {
			returnError = errors.New("error info is nil") // 에러도 없는 경우 처리
		}
	default:
		msg := fmt.Sprintf(e.format, e.args...)
		if e.err != nil {
			returnError = errors.Wrap(e.err, msg) // 에러가 struct 에 포함되어 있으면 wrapping
		} else {
			returnError = errors.New(msg) // 에러가 없으면 msg 만 에러로 변환
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
