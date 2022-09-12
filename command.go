package eventsourcing

import "sync"

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
	c.mapper[et.String()] = cmd
}

// GetCommand | EventType 으로 Command 를 가져오기
func (c Commander[S, R]) GetCommand(et EventType) (cmd Command[S, R], ok bool) {
	c.rwLocker.RLock()
	defer c.rwLocker.RUnlock()
	cmd, ok = c.mapper[et.String()]
	return cmd, ok
}
