package eventsourcing

import (
	"fmt"
	"github.com/rs/xid"
	"time"
)

// Event | Event 의 구조체, R 은 Event 의 담길 요청내용(Request) 의 Type
type Event[R any] struct {
	EventId      EventId      `json:"eventId"`
	PartitionKey PartitionKey `json:"partitionKey"`
	*EventType
	EventNo int       `json:"eventNo"`
	EventAt time.Time `json:"eventAt"`
	Request *R        `json:"request"` // Domain 마다 fit 하게 만들어진 구조체를 넣는다
}

func NewEvent[R any](pk PartitionKey, eventType *EventType, no int, request *R) *Event[R] {
	return &Event[R]{
		EventId:      EventId(xid.New().String()),
		PartitionKey: pk,
		EventType:    eventType,
		EventNo:      no,
		EventAt:      time.Now().UTC(),
		Request:      request,
	}
}

// EventType | Event 의 종류를 판별하는데 쓰이는 정보
type EventType struct {
	Domain   Domain       `json:"domain"`
	Name     EventName    `json:"name"`
	Version  EventVersion `json:"version"`
	NeedLock bool         `json:"needLock"` // Event 검사에 lock 이 필요한지 여부
}

func (e *EventType) String() string {
	return fmt.Sprintf("%s_%s_%s", e.Domain, e.Name, e.Version)
}

// EventId | Event 의 고유 아이디, sorted 해야 한다. Event 는 이 type 을 key 로 삼아야 함
type EventId string

// Domain | Event 가 속한 Domain
type Domain string

// PartitionKey | Event 를 저장할 때 partitioning 하는데 사용하는 key, Domain 마다 정의하게되는 각 State 의 key 를 활용하면 된다.
type PartitionKey string

// EventName | Event 의 이름
type EventName string

// EventVersion | Event 의 버전, Domain, EventName 이 같은데 실제 Command 가 달라야 한다면 Version 을 높인 Event를 새로 정의한다.
type EventVersion string
