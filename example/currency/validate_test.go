package currency

import (
	es "eventsourcing"
	"testing"
)

func TestNoBurned(t *testing.T) {
	type args struct {
		latest   *es.Event[Request]
		snapshot *es.State[State, Request]
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "first_state",
			args: args{
				latest:   NewCreateAmountStateEvent("test", 1, nil),
				snapshot: nil,
			},
		},
		{
			name: "first_state_and_no_burn_event",
			args: args{
				latest:   NewAddAmountEvent("test", 2, &Request{Amount: 100}),
				snapshot: NewState("test"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			NoBurned(tt.args.latest, tt.args.snapshot)
		})
	}
}

func TestNoLock(t *testing.T) {
	type args struct {
		latest   *es.Event[Request]
		snapshot *es.State[State, Request]
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			NoLock(tt.args.latest, tt.args.snapshot)
		})
	}
}
