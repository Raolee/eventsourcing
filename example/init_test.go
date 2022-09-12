package example

import (
	es "eventsourcing"
	"github.com/rs/xid"
	"testing"
)

func TestCommander(t *testing.T) {
	pk := "test_pk"
	e1 := NewModifyAmountEvent(es.PartitionKey(pk), 1, &ExampleRequest{
		Amount: 100,
	})
	idleStatus := IDLE
	e2 := NewChangeStatusEvent(es.PartitionKey(pk), 2, &ExampleRequest{
		Status: (*Status)(&idleStatus),
	})
	value := "value test"
	e3 := NewChangeValueEvent(es.PartitionKey(pk), 3, &ExampleRequest{
		Value: &value,
	})
	e4 := NewLockEvent(es.PartitionKey(pk), 4, nil)
	e5 := NewUnlockEvent(es.PartitionKey(pk), 5, nil)
	e6 := NewModifyAmountEvent(es.PartitionKey(pk), 1, &ExampleRequest{
		Amount: 50,
	})
	e7 := NewModifyAmountEvent(es.PartitionKey(pk), 1, &ExampleRequest{
		Amount: -100,
	})
	e8 := NewModifyAmountEvent(es.PartitionKey(pk), 1, &ExampleRequest{
		Amount: 300,
	})

	events := []*es.Event[ExampleRequest]{
		e1, e2, e3, e4, e5, e6, e7, e8,
	}

	state := NewExampleState(es.PartitionKey(xid.New().String()))

	for _, evt := range events {
		f, _ := ExampleCommander.GetCommand(*evt.EventType)
		state = f(state, evt)
		t.Log(state)
	}
}
