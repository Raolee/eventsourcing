package item

import "errors"

// ReplayEventsWithoutState | 첫 event 부터 replay
func ReplayEventsWithoutState(events ...*Event) (*State, error) {
	state := &State{} // must be not nil
	return replayEvents(state, events)
}

// ReplayEventsWithState | state 부터 event 를 적용하여 replay
func ReplayEventsWithState(state *State, events ...*Event) (*State, error) {
	return replayEvents(state, events)
}

func replayEvents(state *State, events EventList) (*State, error) {
	for e := range events.Iterate() {
		commandFunc, ok := EventCommandMap[e.Name]
		if !ok {
			return state, errors.New("not defined event")
		}
		state = commandFunc(state, e)
	}
	return state, nil
}
