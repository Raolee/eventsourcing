package example

import (
	es "eventsourcing"
)

/**
Event Sourcing 으로 다루는 도메인의 State 를 여기에 구현한다.
*/

var (
	_ es.CommonState[ExampleRequest] = AmountState{}
)

type Status int

const (
	NOTHING = iota
	IDLE
	CLAIM
	BURNED
)

type AmountState struct {
	PartitionKey es.PartitionKey           `json:"partitionKey"`
	Amount       int                       `json:"amount"`
	Status       Status                    `json:"status"`
	Value        *string                   `json:"value"`
	Lock         bool                      `json:"lock"`
	LastEvent    *es.Event[ExampleRequest] `json:"lastEvent"`
}

func NewExampleState(pk es.PartitionKey) *es.State[AmountState, ExampleRequest] {
	return es.NewState[AmountState, ExampleRequest](&AmountState{
		PartitionKey: pk,
		Amount:       0,
		Status:       NOTHING,
		Value:        nil,
		LastEvent:    nil,
	})
}

func (e AmountState) GetPartitionKey() es.PartitionKey {
	return e.PartitionKey
}

func (e AmountState) GetLastEvent() *es.Event[ExampleRequest] {
	return e.LastEvent
}

func (e AmountState) String() string {
	return es.JsonString(e)
}
