package example

import (
	es "eventsourcing"
	"eventsourcing/example/currency"
	"github.com/rs/xid"
	"testing"
)

func TestCommander(t *testing.T) {
	pk := "test_pk"
	e1 := currency.NewAddAmountEvent(es.PartitionKey(pk), 1, &currency.Request{
		Amount: 100,
	})
	idleStatus := currency.IDLE
	e2 := currency.NewChangeStatusEvent(es.PartitionKey(pk), 2, &currency.Request{
		Status: (*currency.Status)(&idleStatus),
	})
	value := "value test"
	e3 := currency.NewChangeValueEvent(es.PartitionKey(pk), 3, &currency.Request{
		Value: &value,
	})
	e4 := currency.NewAddAmountEvent(es.PartitionKey(pk), 1, &currency.Request{
		Amount: 50,
	})
	e5 := currency.NewAddAmountEvent(es.PartitionKey(pk), 1, &currency.Request{
		Amount: -100,
	})
	e6 := currency.NewAddAmountEvent(es.PartitionKey(pk), 1, &currency.Request{
		Amount: 300,
	})

	events := []*es.Event[currency.Request]{
		e1, e2, e3, e4, e5, e6,
	}

	state := currency.NewState(es.PartitionKey(xid.New().String()))

	for _, evt := range events {
		f, _ := CurrencyCommander.GetCommand(*evt.EventType)
		state = f(state, evt)
		t.Log(state)
	}
}
