package eventsourcing

// CommonState | Domain 마다 정의하는 State 가 구현해야할 인터페이스, R는 Event 의 Request 구조체 타입
type CommonState[R any] interface {
	GetPartitionKey() PartitionKey // 파티션 키를 가져온다
	GetLastEvent() *Event[R]       // State 의 마지막 이벤트를 가져온다
	String() string                // State 의 ToString() func
}

type State[S CommonState[R], R any] struct {
	state *S
}

func NewState[S CommonState[R], R any](state *S) *State[S, R] {
	return &State[S, R]{
		state: state,
	}
}

func (s *State[S, R]) State() *S {
	return s.state
}

func (s *State[S, R]) String() string {
	return (*s.state).String()
}
