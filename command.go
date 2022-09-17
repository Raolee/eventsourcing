package eventsourcing

import "sync"

// 이벤트 소싱에서 사용할 커맨드 인터페이스와 구조체를 정의한다.
//
// type Command[S CommonState[R], R any
// - 이벤트 소싱에서 구현할 Command 의 기본 타입
// - 제네릭 설명
//   1) S : 직접 구현한 State Type 을 넣음
//   2) R : Event[R] 에 넣는 R을 넣음
// - 인자 설명
//   1) *State[S, R] : 현재 State
//   2) *Event[R] : State 에 반영할 이벤트
// - 리턴 설명
//   1) *State[S, R] : State 에 Event 를 반영한 후 변한 State
//
// struct Commander
// - 미리 정의한 Event Type 과 이에 맞춰 구현한 Command 를 Commander 에 Set 하고, Get 할 수 있게 만든 구조체
// - Event Sourcing 의 각 기능은 이 Commander 를 주입받아서 유용하게 사용한다.
// - SetCommand 설명
//   1) Command 를 공통 로직을 태우게 wrapping 하여 저장
// - GetCommand 설명
//   2) 공통 로직을 포함시킨 Command 를 가져옴

// Command | Event 를 실제로 수행하는 Func Type
type Command[S CommonState[R], R any] func(state *State[S, R], event *Event[R]) *State[S, R]

// Commander | EventType 과 매핑되어 있는 Command 를 관리
type Commander[S CommonState[R], R any] struct {
	mapper   map[string]Command[S, R]
	rwLocker sync.RWMutex
}

func NewCommander[S CommonState[R], R any]() *Commander[S, R] {
	return &Commander[S, R]{
		mapper:   make(map[string]Command[S, R]),
		rwLocker: sync.RWMutex{},
	}
}

// SetCommand | EventType 과 Command 를 설정하기
func (c Commander[S, R]) SetCommand(et EventType, cmd Command[S, R]) {
	c.rwLocker.Lock()
	defer c.rwLocker.Unlock()

	// AOP 를 못하니까... Command 를 다시 Command 로 감싸서 공통 로직을 적용시킨다.
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

// GetCommand | EventType 으로 Command 를 가져오기
func (c Commander[S, R]) GetCommand(et EventType) (cmd Command[S, R], ok bool) {
	c.rwLocker.RLock()
	defer c.rwLocker.RUnlock()
	cmd, ok = c.mapper[et.String()]
	return cmd, ok
}
