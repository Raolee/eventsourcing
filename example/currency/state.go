package currency

import (
	es "eventsourcing"
)

/**
Event Sourcing 으로 다루는 도메인의 State 를 여기에 구현한다.
*/

var (
	_ es.CommonState[Request] = State{}
)

type Status int

const (
	NOTHING = iota
	IDLE
	CLAIM
	BURNED
)

type State struct {
	PartitionKey es.PartitionKey    `json:"partitionKey"`
	Amount       int                `json:"amount"`
	Status       Status             `json:"status"`
	Value        *string            `json:"value"`
	LastEvent    *es.Event[Request] `json:"lastEvent"`
}

func NewState(pk es.PartitionKey) *es.State[State, Request] {
	return es.NewState[State, Request](&State{
		PartitionKey: pk,
		Amount:       0,
		Status:       NOTHING,
		Value:        nil,
		LastEvent:    nil,
	})
}

func (e State) GetPartitionKey() es.PartitionKey {
	return e.PartitionKey
}

func (e State) GetLastEvent() *es.Event[Request] {
	return e.LastEvent
}

func (e State) String() string {
	return es.JsonString(e)
}
