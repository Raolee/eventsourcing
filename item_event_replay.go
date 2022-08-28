package eventsourcing

import "errors"

func ItemReplay(events []*ItemEvent) (*ItemState, error) {
	var state *ItemState
	for _, e := range events {
		f, ok := ItemEvents[e.Name]
		if !ok {
			return state, errors.New("not defined event")
		}
		state = f(state, e)
	}
	return state, nil
}
