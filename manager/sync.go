package manager

import "eventsourcing"

// TODO 매니저를 역할별로 더 나누어야 할 듯
// 예상
// validator
// producer
// consumer
// replayer
// querier

type Manager[S eventsourcing.CommonState[R], R any] interface {
	Validate(pk eventsourcing.PartitionKey, et *eventsourcing.EventType) error               // 이벤트를 실행해도 되는지 유효성 검사를 한다.
	Put(pk eventsourcing.PartitionKey, et *eventsourcing.EventType, req *R) error            // 이벤트를 저장한다.
	ApplyEvents(pk eventsourcing.PartitionKey) error                                         // 이벤트를 적용한다.
	GetEvents(pk eventsourcing.PartitionKey, eventNo int) ([]*eventsourcing.Event[R], error) // eventNo 보다 큰 이벤트 리스트를 가져온다.
	GetLatestState(pk eventsourcing.PartitionKey) (*eventsourcing.State[S, R], error)        // 이벤트로 리플레이한 최신 스테이트를 가져온다.
	GetStateSnapshot(pk eventsourcing.PartitionKey) (*eventsourcing.State[S, R], error)      // 스냅샷의 스테이트를 가져온다.
}

// baseManager | 가장 기본적인 이벤트 소싱 매니저, 메세지 스트림을 사용하지 않는다.
type baseManager[S eventsourcing.CommonState[R], R any] struct {
	processor *eventsourcing.Processor[S, R]
	validator *eventsourcing.Validator[S, R]
	es        eventsourcing.EventStorage[R]
	ss        eventsourcing.StateSnapshotStorage[S, R]
	rule      *eventsourcing.Rule
}

// NewBaseManager | 기본적인 매니저를 생성한다. 아래의 규칙을 따름
//
// 1. Validate : Put 전에 호출하여 저장가능한지 확인
//
// 2. Put : Validate 에서 문제가 없으면 Event 를 EventStorage 에 저장
//
// 3. ApplyEvents : 아직 반영하지 않은 이벤트를 적용 (= StateSnapshotStorage 저장)
//
// 4. GetEvents : EventStorage 에서 pk 로 이벤트를 조회
//
// 5. GetStateSnapshot : StateSnapshotStorage + EventStorage 를 합쳐서 최신 State 를 조회
func NewBaseManager[S eventsourcing.CommonState[R], R any](
	rule *eventsourcing.Rule,
	c *eventsourcing.Processor[S, R],
	v *eventsourcing.Validator[S, R],
	es eventsourcing.EventStorage[R],
	ss eventsourcing.StateSnapshotStorage[S, R],
) Manager[S, R] {
	r := eventsourcing.NewDefaultRule()
	r.Merge(rule)
	return &baseManager[S, R]{
		processor: c,
		validator: v,
		es:        es,
		ss:        ss,
		rule:      r,
	}
}

func (b *baseManager[S, R]) replay(pk eventsourcing.PartitionKey, current *eventsourcing.State[S, R], events []*eventsourcing.Event[R]) (state *eventsourcing.State[S, R], err error) {
	for _, e := range events {
		cmd, ok := b.processor.GetProcess(*e.EventType)
		if !ok {
			return nil, eventsourcing.NewNoHasCommandError(pk, e.EventType)
		}
		current = cmd(current, e)
	}
	return current, nil
}

// Validate | 이벤트를 적용할 수 있는지 Validating
func (b *baseManager[S, R]) Validate(pk eventsourcing.PartitionKey, et *eventsourcing.EventType) (err error) {
	defer eventsourcing.HandleError(&err)

	// get validates
	validates, ok := b.validator.GetValidates(*et)
	if !ok || len(validates) == 0 {
		return nil // validate 가 지정되지 않았으므로 검사 없이 끝
	}

	// get snapshot
	snapshot, err := b.GetStateSnapshot(pk)
	if err != nil {
		return err
	}
	if snapshot == nil {
		err = b.ApplyEvents(pk) // 스냅샷이 없으면 최신으로 업데이트 한다
		if err != nil {
			return err
		}
		snapshot, err = b.GetStateSnapshot(pk) // 다시 가져옴
	}

	// get the latest event
	latest, err := b.es.GetLastEvent(pk)
	if err != nil {
		return eventsourcing.NewEventStorageError(err)
	}

	// validation
	for _, v := range validates {
		v(latest, snapshot) // validate 가 있는 경우만 체크, 에러가 있으면 panic 발생
	}

	return nil
}

// Put | 이벤트를 저장합니다.
func (b *baseManager[S, R]) Put(pk eventsourcing.PartitionKey, et *eventsourcing.EventType, req *R) (err error) {
	defer eventsourcing.HandleError(&err)

	// event 생성 및 command 실행
	var no int
	no, err = b.es.IncreaseEventNo(pk) // 이벤트 번호를 받아옴
	if err != nil {
		return eventsourcing.NewDispenseEventNoError(err, pk)
	}
	event := eventsourcing.NewEvent[R](pk, et, no, req) // 이벤트 생성
	err = b.es.AddEvent(event)                          // 이벤트를 EventStorage 에 추가
	if err != nil {
		return eventsourcing.NewEventStorageError(err)
	}
	return nil
}

// ApplyEvents | pk 에 쌓여있는 이벤트 들을 적용합니다. => snapshot 에 반영
func (b *baseManager[S, R]) ApplyEvents(pk eventsourcing.PartitionKey) (err error) {
	defer eventsourcing.HandleError(&err)

	/**
	1. 먼저 스냅샷을 먼저 조회
	case 1. 스냅샷이 없는 경우 : 전체 event 를 가져온다
	case 2. 스냅샷이 있는 경우 : 스냅샷의 eventNo 이후로 event 를 가져온다

	2. events 가 조회된 것이 있는지 확인
	case 1. events 가 없다면 이미 스냅샷이 최신 상태
	case 2. events 와 state 로 replay

	3. replay 된 state 를 snapshot 에 저장
	*/
	// 스냅샷이 있는지 조회, 없다면 만들어 주어야 함

	state, err := b.GetStateSnapshot(pk)

	// replay 할 event 리스트를 만듦
	var events []*eventsourcing.Event[R]
	var eventNo int
	if state != nil { // 스냅샷이 존재하는 경우, snapshot 이후의 events 만 가져온다
		eventNo = (*state.State()).GetLastEvent().EventNo
	}

	events, err = b.GetEvents(pk, eventNo)
	if err != nil {
		return
	}

	if len(events) == 0 {
		return nil // 이미 스냅샷이 최신이므로 리턴
	}

	// replay events, event 로 현재 state 를 만든다
	state, err = b.replay(pk, state, events)
	if err != nil {
		return err
	}

	// snapshot 에 저장
	err = b.ss.SaveSnapshot(pk, state)
	if err != nil {
		return eventsourcing.NewSnapshotStorageError(err)
	}
	return nil
}

// GetEvents | pk 의 event 리스트를 가져옵니다
func (b *baseManager[S, R]) GetEvents(pk eventsourcing.PartitionKey, afterEventNo int) (events []*eventsourcing.Event[R], err error) {
	defer eventsourcing.HandleError(&err)

	if afterEventNo == 0 {
		events, err = b.es.GetEvents(pk)
		if err != nil {
			return nil, eventsourcing.NewEventStorageError(err)
		}
	} else {
		events, err = b.es.GetEventsAfterEventNo(pk, afterEventNo) // eventNo 이후의 리스트를 조회
		if err != nil {
			return nil, eventsourcing.NewEventStorageError(err)
		}
	}
	return
}

// GetLatestState | pk 의 이벤트를 replay 해서 최신 state 를 만듭니다.
func (b *baseManager[S, R]) GetLatestState(pk eventsourcing.PartitionKey) (state *eventsourcing.State[S, R], err error) {
	defer eventsourcing.HandleError(&err)

	var events []*eventsourcing.Event[R]
	events, err = b.es.GetEvents(pk)
	if err != nil {
		return nil, eventsourcing.NewEventStorageError(err)
	}
	state, err = b.replay(pk, nil, events)
	if err != nil {
		return nil, err
	}
	return
}

// GetStateSnapshot | 현재 snapshot 에 저장된 이벤트를 가져옵니다.
func (b *baseManager[S, R]) GetStateSnapshot(pk eventsourcing.PartitionKey) (state *eventsourcing.State[S, R], err error) {
	defer eventsourcing.HandleError(&err)

	state, err = b.ss.GetSnapshot(pk)
	if err != nil {
		return nil, eventsourcing.NewSnapshotStorageError(err)
	}
	return
}
