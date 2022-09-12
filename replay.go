package eventsourcing

import (
	"errors"
)

// ReplayEventsWithoutState | 첫 event 부터 replay
func ReplayEventsWithoutState[S CommonState[R], R any](
	commander *Commander[S, R],
	init *State[S, R],
	events ...*Event[R],
) (
	*State[S, R],
	error,
) {
	if init != nil {
		return nil, errors.New("init state must be not nil")
	}
	return replayEvents[S, R](commander, init, events...)
}

// ReplayEventsWithState | state 부터 event 를 적용하여 replay
func ReplayEventsWithState[S CommonState[R], R any](
	commander *Commander[S, R],
	state *State[S, R],
	events ...*Event[R],
) (
	*State[S, R],
	error,
) {
	return replayEvents[S, R](commander, state, events...)
}

func replayEvents[S CommonState[R], R any](
	commander *Commander[S, R],
	state *State[S, R],
	events ...*Event[R],
) (
	*State[S, R], error,
) {
	for _, e := range events {
		cmd, ok := commander.GetCommand(*e.EventType)
		if !ok {
			return state, errors.New("not defined event")
		}
		state = cmd(state, e)
	}
	return state, nil
}
