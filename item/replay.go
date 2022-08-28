package item

import "errors"

func EventReplay(events *EventList) (*State, error) {
	state := &State{} // must be not nil
	for e := range events.Iterate() {
		f, ok := Events[e.Name]
		if !ok {
			return state, errors.New("not defined event")
		}
		state = f(state, e)
	}
	return state, nil
}
