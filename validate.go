package eventsourcing

import "sync"

// Validate | Event 를 받아들일지 판단하는 Validate Func Type
type Validate[S CommonState[R], R any] func(latest *Event[R], snapshot *State[S, R])

type Validator[S CommonState[R], R any] struct {
	mapper   map[string][]Validate[S, R]
	rwLocker sync.RWMutex
}

func NewValidator[S CommonState[R], R any]() *Validator[S, R] {
	return &Validator[S, R]{
		mapper:   make(map[string][]Validate[S, R]),
		rwLocker: sync.RWMutex{},
	}
}

// SetValidates | EventType 과 Validate 를 설정하기
func (v Validator[S, R]) SetValidates(et EventType, validates ...Validate[S, R]) {
	v.rwLocker.Lock()
	defer v.rwLocker.Unlock()
	v.mapper[et.String()] = validates
}

// GetValidates | EventType 으로 Validate 를 가져오기
func (v Validator[S, R]) GetValidates(et EventType) (validates []Validate[S, R], ok bool) {
	v.rwLocker.RLock()
	defer v.rwLocker.RUnlock()
	validates, ok = v.mapper[et.String()]
	return validates, ok
}
