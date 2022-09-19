package eventsourcing

type eventStreamManager[S CommonState[R], R any] struct {
	c    *Commander[S, R]
	v    *Validator[S, R]
	es   EventStorage[R]
	ss   StateSnapshotStorage[S, R]
	rule *Rule
}

func NewEventStreamManager[S CommonState[R], R any](
	rule *Rule,
	c *Commander[S, R],
	v *Validator[S, R],
	es EventStorage[R],
	ss StateSnapshotStorage[S, R],
) Manager[S, R] {
	r := newDefaultRule()
	r.Merge(rule)
	return &baseManager[S, R]{
		c:    c,
		v:    v,
		es:   es,
		ss:   ss,
		rule: r,
	}
}

func (e *eventStreamManager[S, R]) Validate(pk PartitionKey, et *EventType) error {
	//TODO implement me
	panic("implement me")
}

func (e *eventStreamManager[S, R]) Put(pk PartitionKey, et *EventType, req *R) error {
	//TODO implement me
	panic("implement me")
}

func (e *eventStreamManager[S, R]) ApplyEvents(pk PartitionKey) error {
	//TODO implement me
	panic("implement me")
}

func (e *eventStreamManager[S, R]) GetEvents(pk PartitionKey, eventNo int) ([]*Event[R], error) {
	//TODO implement me
	panic("implement me")
}

func (e *eventStreamManager[S, R]) GetLatestState(pk PartitionKey) (*State[S, R], error) {
	//TODO implement me
	panic("implement me")
}

func (e *eventStreamManager[S, R]) GetStateSnapshot(pk PartitionKey) (*State[S, R], error) {
	//TODO implement me
	panic("implement me")
}
