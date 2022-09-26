package manager

import (
	"eventsourcing"
)

// TODO : async 를 처리하는 도메인 추가 필요
// [MEMO] : Async 는 매니저가 불필요한 것 같음
type asyncManager[S eventsourcing.CommonState[R], R any] struct {
	processor              *eventsourcing.Processor[S, R]
	validator              *eventsourcing.Validator[S, R]
	eventStorage           eventsourcing.EventStorage[R]
	snapshotStorage        eventsourcing.StateSnapshotStorage[S, R]
	latestEventTypeStorage eventsourcing.LatestEventTypeStorage
	rule                   *eventsourcing.Rule
}

func NewAsyncManager[S eventsourcing.CommonState[R], R any](
	rule *eventsourcing.Rule,
	p *eventsourcing.Processor[S, R],
	v *eventsourcing.Validator[S, R],
	es eventsourcing.EventStorage[R],
	ss eventsourcing.StateSnapshotStorage[S, R],
	lets eventsourcing.LatestEventTypeStorage,
) Manager[S, R] {
	r := eventsourcing.NewDefaultRule()
	r.Merge(rule)
	return &asyncManager[S, R]{
		processor:              p,
		validator:              v,
		eventStorage:           es,
		snapshotStorage:        ss,
		latestEventTypeStorage: lets,
		rule:                   r,
	}
}

func (e *asyncManager[S, R]) Validate(pk eventsourcing.PartitionKey, et *eventsourcing.EventType) error {
	//TODO implement me
	panic("implement me")
}

func (e *asyncManager[S, R]) Put(pk eventsourcing.PartitionKey, et *eventsourcing.EventType, req *R) error {
	//TODO implement me
	panic("implement me")
}

func (e *asyncManager[S, R]) ApplyEvents(pk eventsourcing.PartitionKey) error {
	//TODO implement me
	panic("implement me")
}

func (e *asyncManager[S, R]) GetEvents(pk eventsourcing.PartitionKey, eventNo int) ([]*eventsourcing.Event[R], error) {
	//TODO implement me
	panic("implement me")
}

func (e *asyncManager[S, R]) GetLatestState(pk eventsourcing.PartitionKey) (*eventsourcing.State[S, R], error) {
	//TODO implement me
	panic("implement me")
}

func (e *asyncManager[S, R]) GetStateSnapshot(pk eventsourcing.PartitionKey) (*eventsourcing.State[S, R], error) {
	//TODO implement me
	panic("implement me")
}
