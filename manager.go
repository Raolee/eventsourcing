package eventsourcing

type Manager[S CommonState[R], R any] interface {
	Validate(pk PartitionKey, et *EventType) error          // 이벤트를 실행해도 되는지 유효성 검사를 한다.
	Put(pk PartitionKey, et *EventType, req *R) error       // 이벤트를 저장한다.
	UpdateStateSnapshot(pk PartitionKey) error              // State Snapshot 을 최신 event 로 업데이트 한다.
	GetEvents(pk PartitionKey) ([]*Event[R], error)         // 이벤트 리스트를 가져온다.
	GetLatestState(pk PartitionKey) (*State[S, R], error)   // 이벤트로 리플레이한 최신 스테이트를 가져온다.
	GetStateSnapshot(pk PartitionKey) (*State[S, R], error) // 스냅샷의 스테이트를 가져온다.
}

// baseManager | 가장 기본적인 이벤트 소싱 매니저
type baseManager[S CommonState[R], R any] struct {
	c    *Commander[S, R]
	v    *Validator[S, R]
	es   EventStorage[R]
	ss   StateSnapshotStorage[S, R]
	rule *Rule
}

// NewBaseManager | 기본적인 매니저를 생성한다. 아래의 규칙을 따름
//
// 1. Validate : Put 전에 호출하여 저장가능한지 확인
//
// 2. Put : Validate 에서 문제가 없으면 Event 를 EventStorage 에 저장
//
// 3. UpdateStateSnapshot : StateSnapshot 을 최신 event 로 업데이트
//
// 4. GetEvents : EventStorage 에서 pk 로 이벤트를 조회
//
// 5. GetStateSnapshot : StateSnapshotStorage + EventStorage 를 합쳐서 최신 State 를 조회
func NewBaseManager[S CommonState[R], R any](
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

func (b *baseManager[S, R]) lock(pk PartitionKey, et *EventType) error {
	if et.NeedLock {
		already, err := b.es.Lock(pk)
		if err != nil {
			return NewLockedEventError(err, pk, et)
		}
		if already {
			return NewLockedEventError(nil, pk, et)
		}
	}
	return nil
}

func (b *baseManager[S, R]) unlock(pk PartitionKey, et *EventType) error {
	if et.NeedLock {
		already, err := b.es.Unlock(pk)
		if err != nil {
			return err
		}
		if already {
			return nil
		}
	}
	return nil
}

func (b *baseManager[S, R]) replay(pk PartitionKey, current *State[S, R], events []*Event[R]) (state *State[S, R], err error) {
	for _, e := range events {
		cmd, ok := b.c.GetCommand(*e.EventType)
		if !ok {
			return nil, NewNoHasCommandError(pk, e.EventType)
		}
		current = cmd(current, e)
	}
	return current, nil
}

func (b *baseManager[S, R]) Validate(pk PartitionKey, et *EventType) (err error) {
	defer handleError(&err)

	lock, e := b.es.GetLock(pk) // 이벤트의 잠금 상태를 가져옴
	if e != nil {
		return e // 에러면 리턴
	}
	if lock {
		return NewLockedEventError(nil, pk, et) // 잠겼으면 Validate 실패
	}

	// get validates
	validates, ok := b.v.GetValidates(*et)
	if !ok || len(validates) == 0 {
		return nil // validate 가 지정되지 않았으므로 검사 없이 끝
	}

	// get snapshot
	snapshot, err := b.GetStateSnapshot(pk)
	if err != nil {
		return err
	}
	if snapshot == nil {
		err = b.UpdateStateSnapshot(pk) // 스냅샷이 없으면 최신으로 업데이트 한다
		if err != nil {
			return err
		}
		snapshot, err = b.GetStateSnapshot(pk) // 다시 가져옴
	}

	// get the latest event
	latest, err := b.es.GetLastEvent(pk)
	if err != nil {
		return NewEventStorageError(err)
	}

	// validation
	for _, v := range validates {
		v(latest, snapshot) // validate 가 있는 경우만 체크, 에러가 있으면 panic 발생
	}

	return nil
}

func (b *baseManager[S, R]) Put(pk PartitionKey, et *EventType, req *R) (err error) {
	defer handleError(&err)
	err = b.lock(pk, et)
	if err != nil {
		return err
	}
	defer func(err *error) {
		e := b.unlock(pk, et)
		if e != nil {
			*err = e
		}
	}(&err)

	// event 생성 및 command 실행
	var no int
	no, err = b.es.IncreaseEventNo(pk) // 이벤트 번호를 받아옴
	if err != nil {
		return NewDispenseEventNoError(err, pk)
	}
	event := NewEvent[R](pk, et, no, req) // 이벤트 생성
	err = b.es.AddEvent(event)            // 이벤트를 EventStorage 에 추가
	if err != nil {
		return NewEventStorageError(err)
	}
	return nil
}

func (b *baseManager[S, R]) UpdateStateSnapshot(pk PartitionKey) (err error) {
	defer handleError(&err)

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
	var current *State[S, R] = nil
	current, err = b.ss.GetSnapshot(pk)
	if err != nil {
		return NewSnapshotStorageError(err)
	}

	var last *Event[R]
	last, err = b.es.GetLastEvent(pk)

	if (*current.State()).GetLastEvent().EventNo == last.EventNo {
		return nil // 이미 최신이므로 건너뜀
	}

	// replay 할 event 리스트를 만듦
	var events []*Event[R]
	if current != nil { // 스냅샷이 존재하는 경우, snapshot 이후의 events 만 가져온다
		s := *current.State()
		events, err = b.es.GetEventsAfterEventNo(pk, s.GetLastEvent().EventNo) // 마지막 이벤트 이후의 이벤트 리스트를 조회
	} else { // 스냅샷이 존재하지 않는 경우, 전체를 가져온다
		events, err = b.es.GetEvents(pk) // 처음부터 끝까지의 이벤트 리스트를 조회
	}

	if len(events) == 0 {
		return nil // 이미 스냅샷이 최신이므로 리턴
	}

	// replay events, event 로 현재 state 를 만든다
	current, err = b.replay(pk, current, events)
	if err != nil {
		return err
	}

	// snapshot 에 저장
	err = b.ss.SaveSnapshot(pk, current)
	if err != nil {
		return NewSnapshotStorageError(err)
	}
	return nil
}

func (b *baseManager[S, R]) GetEvents(pk PartitionKey) (events []*Event[R], err error) {
	defer handleError(&err)

	events, err = b.es.GetEvents(pk)
	if err != nil {
		return nil, NewEventStorageError(err)
	}
	return
}

func (b *baseManager[S, R]) GetLatestState(pk PartitionKey) (state *State[S, R], err error) {
	defer handleError(&err)

	var events []*Event[R]
	events, err = b.es.GetEvents(pk)
	if err != nil {
		return nil, NewEventStorageError(err)
	}
	state, err = b.replay(pk, nil, events)
	if err != nil {
		return nil, err
	}
	return
}

func (b *baseManager[S, R]) GetStateSnapshot(pk PartitionKey) (state *State[S, R], err error) {
	defer handleError(&err)

	state, err = b.ss.GetSnapshot(pk)
	if err != nil {
		return nil, NewSnapshotStorageError(err)
	}
	return
}
