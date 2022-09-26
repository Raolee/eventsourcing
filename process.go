package eventsourcing

import "sync"

// Event Sourcing 에서 사용할 Process 인터페이스와 구조체를 정의한다.
//
// type Process[S CommonState[R], R any
// - 이벤트 소싱에서 구현할 Process 의 기본 타입
// - 제네릭 설명
//   1) S : 직접 구현한 State Type 을 넣음
//   2) R : Event[R] 에 넣는 R을 넣음
// - 인자 설명
//   1) *State[S, R] : 현재 State
//   2) *Event[R] : State 에 반영할 이벤트
// - 리턴 설명
//   1) *State[S, R] : State 에 Event 를 반영한 후 변한 State
//
// struct Processor
// - 미리 정의한 Event Type 과 이에 맞춰 구현한 Process 를 Processor 에 Set 하고, Get 할 수 있게 만든 구조체
// - Event Sourcing 의 각 기능은 이 Processor 를 주입받아서 유용하게 사용한다.
// - SetProcess 설명
//   1) Process 를 공통 로직을 태우게 wrapping 하여 저장
// - GetProcess 설명
//   2) 공통 로직을 포함시킨 Process 를 가져옴

// Process | Event 를 실제로 수행하는 Func Type
type Process[S CommonState[R], R any] func(state *State[S, R], event *Event[R]) *State[S, R]

// Processor | EventType 과 매핑되어 있는 Process 를 관리
type Processor[S CommonState[R], R any] struct {
	mapper   map[string]Process[S, R] // key : event type, value : process
	rwLocker sync.RWMutex
}

func NewProcessor[S CommonState[R], R any]() *Processor[S, R] {
	return &Processor[S, R]{
		mapper:   make(map[string]Process[S, R]),
		rwLocker: sync.RWMutex{},
	}
}

// SetProcess | EventType 과 Process 를 설정하기
func (c Processor[S, R]) SetProcess(et EventType, cmd Process[S, R]) {
	c.rwLocker.Lock()
	defer c.rwLocker.Unlock()

	// AOP 를 못하니까... Process 를 다시 Process 로 감싸서 공통 로직을 적용시킨다.
	wrapped := func(state *State[S, R], event *Event[R]) *State[S, R] {
		if state != nil {
			if event.EventNo <= (*state.State()).GetLastEvent().EventNo {
				return state // state 최신 이벤트 번호가 요청온 event 번호보다 더 크거나 같다면, event 처리 무시
			}
		}
		return cmd(state, event)
	}
	c.mapper[et.String()] = wrapped
}

// GetProcess | EventType 으로 Process 를 가져오기
func (c Processor[S, R]) GetProcess(et EventType) (cmd Process[S, R], ok bool) {
	c.rwLocker.RLock()
	defer c.rwLocker.RUnlock()
	cmd, ok = c.mapper[et.String()]
	return cmd, ok
}
