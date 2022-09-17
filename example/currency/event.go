package currency

import (
	es "eventsourcing"
)

/**
도메인에서 다루는 Event 를 여기에 정의한다
*/
var (
	CreateAmountStateEvent = es.EventType{Domain: "currency", Name: "create_currency_state", Version: "v1"}
	AddAmountEvent         = es.EventType{Domain: "currency", Name: "add_amount", Version: "v1"}
	MinusAmountEvent       = es.EventType{Domain: "currency", Name: "minus_amount", Version: "v1"}
	ChangeStatusEvent      = es.EventType{Domain: "currency", Name: "change_status", Version: "v1"}
	ChangeValueEvent       = es.EventType{Domain: "currency", Name: "change_value", Version: "v1", NeedLock: true}
	ChangeValueV2Event     = es.EventType{Domain: "currency", Name: "change_value", Version: "v2", NeedLock: true} // 버전 2 이벤트
	BurnEvent              = es.EventType{Domain: "currency", Name: "burn", Version: "v1", NeedLock: true}
)

type Request struct {
	Amount int
	Status *Status
	Value  *string
}

func NewCreateAmountStateEvent(pk es.PartitionKey, no int, request *Request) *es.Event[Request] {
	return es.NewEvent[Request](pk, &CreateAmountStateEvent, no, request)
}

func NewAddAmountEvent(pk es.PartitionKey, no int, request *Request) *es.Event[Request] {
	return es.NewEvent[Request](pk, &AddAmountEvent, no, request)
}

func NewMinusAmountEvent(pk es.PartitionKey, no int, request *Request) *es.Event[Request] {
	return es.NewEvent[Request](pk, &MinusAmountEvent, no, request)
}

func NewChangeStatusEvent(pk es.PartitionKey, no int, request *Request) *es.Event[Request] {
	return es.NewEvent[Request](pk, &ChangeStatusEvent, no, request)
}

func NewChangeValueEvent(pk es.PartitionKey, no int, request *Request) *es.Event[Request] {
	return es.NewEvent[Request](pk, &ChangeValueEvent, no, request)
}

func NewChangeValueV2Event(pk es.PartitionKey, no int, request *Request) *es.Event[Request] {
	return es.NewEvent[Request](pk, &ChangeValueV2Event, no, request)
}

func NewBurnEvent(pk es.PartitionKey, no int, request *Request) *es.Event[Request] {
	return es.NewEvent[Request](pk, &BurnEvent, no, request)
}
