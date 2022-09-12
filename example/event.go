package example

import (
	es "eventsourcing"
)

/**
도메인에서 다루는 Event 를 여기에 정의한다
*/
var (
	CreateAmountStateEvent = es.EventType{Domain: "example", Name: "create_amount_state", Version: "v1"}
	ModifyAmountEvent      = es.EventType{Domain: "example", Name: "modify_amount", Version: "v1"}
	ChangeStatusEvent      = es.EventType{Domain: "example", Name: "change_status", Version: "v1"}
	ChangeValueEvent       = es.EventType{Domain: "example", Name: "change_value", Version: "v1", NeedLock: true}
	ChangeValueV2Event     = es.EventType{Domain: "example", Name: "change_value", Version: "v2", NeedLock: true} // 버전 2 이벤트
	LockEvent              = es.EventType{Domain: "example", Name: "lock", Version: "v1"}
	UnlockEvent            = es.EventType{Domain: "example", Name: "unlock", Version: "v1"}
	BurnEvent              = es.EventType{Domain: "example", Name: "burn", Version: "v1", NeedLock: true}
)

type ExampleRequest struct {
	Amount int
	Status *Status
	Value  *string
}

func NewCreateAmountStateEvent(pk es.PartitionKey, no int, request *ExampleRequest) *es.Event[ExampleRequest] {
	return es.NewEvent[ExampleRequest](pk, &CreateAmountStateEvent, no, request)
}

func NewModifyAmountEvent(pk es.PartitionKey, no int, request *ExampleRequest) *es.Event[ExampleRequest] {
	return es.NewEvent[ExampleRequest](pk, &ModifyAmountEvent, no, request)
}

func NewChangeStatusEvent(pk es.PartitionKey, no int, request *ExampleRequest) *es.Event[ExampleRequest] {
	return es.NewEvent[ExampleRequest](pk, &ChangeStatusEvent, no, request)
}

func NewChangeValueEvent(pk es.PartitionKey, no int, request *ExampleRequest) *es.Event[ExampleRequest] {
	return es.NewEvent[ExampleRequest](pk, &ChangeValueEvent, no, request)
}

func NewChangeValueV2Event(pk es.PartitionKey, no int, request *ExampleRequest) *es.Event[ExampleRequest] {
	return es.NewEvent[ExampleRequest](pk, &ChangeValueV2Event, no, request)
}

func NewLockEvent(pk es.PartitionKey, no int, request *ExampleRequest) *es.Event[ExampleRequest] {
	return es.NewEvent[ExampleRequest](pk, &LockEvent, no, request)
}

func NewUnlockEvent(pk es.PartitionKey, no int, request *ExampleRequest) *es.Event[ExampleRequest] {
	return es.NewEvent[ExampleRequest](pk, &UnlockEvent, no, request)
}

func NewBurnEvent(pk es.PartitionKey, no int, request *ExampleRequest) *es.Event[ExampleRequest] {
	return es.NewEvent[ExampleRequest](pk, &BurnEvent, no, request)
}
