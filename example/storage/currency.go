package storage

import (
	es "eventsourcing"
	"eventsourcing/example/currency"
	"sync"
	"sync/atomic"
)

type Counter struct {
	Count int32
}

func (c *Counter) Increase(delta int32) int32 {
	return atomic.AddInt32(&c.Count, delta)
}

type CurrencyMemoryEventStorage struct {
	eventNoStorage   map[es.PartitionKey]*Counter              // pk 의 event 번호를 저장하는 스토리지
	pkGroupStorage   map[es.PartitionKey][]es.EventId          // pk 의 event id 리스트를 저장하는 스토리지
	eventStorage     map[es.EventId]es.Event[currency.Request] // event id 별로 event 를 저장하는 스토리지
	eventLockStorage map[es.PartitionKey]bool                  // pk 의 event 자체의 lock 값을 저장하는 스토리지
	pkLockers        map[es.PartitionKey]*sync.RWMutex         // pk 안에서 dirty read 를 방지하기 위한 RWMutex
	esLocker         sync.Mutex                                // event storage 자체적으로 사용하는 Mutex
}

func NewCurrencyEventStorage() es.EventStorage[currency.Request] {
	return &CurrencyMemoryEventStorage{
		eventNoStorage:   make(map[es.PartitionKey]*Counter),
		pkGroupStorage:   make(map[es.PartitionKey][]es.EventId),
		eventStorage:     make(map[es.EventId]es.Event[currency.Request]),
		eventLockStorage: make(map[es.PartitionKey]bool),
		pkLockers:        make(map[es.PartitionKey]*sync.RWMutex),
	}
}

func (a *CurrencyMemoryEventStorage) getPkLocker(pk es.PartitionKey) *sync.RWMutex {
	// pk locker 중복 할당 방지
	if _, ok := a.pkLockers[pk]; !ok {
		a.esLocker.Lock()
		defer a.esLocker.Unlock()
		if _, ok = a.pkLockers[pk]; !ok {
			a.pkLockers[pk] = &sync.RWMutex{}
		}
	}
	// locker ptr 을 리턴해야 copy 이슈로 lock 이 걸리지 않는 이슈가 발생하지 않음
	return a.pkLockers[pk]
}

func (a *CurrencyMemoryEventStorage) IncreaseEventNo(pk es.PartitionKey) (eventNo int, err error) {
	// event No 는 pk 별로 atomic 하게 증가시켜야 함
	// 따라서, increase 시 pk 별로 mutex lock 을 건다.

	// counter 중복 할당 방지
	if _, ok := a.eventNoStorage[pk]; !ok {
		pkLocker := a.getPkLocker(pk) // pk 에 대해서만 lock 이 걸리면 되므로 pk pkLockers 를 가져온다
		pkLocker.Lock()
		defer pkLocker.Unlock()
		if _, ok = a.eventLockStorage[pk]; !ok {
			a.eventNoStorage[pk] = &Counter{0}
		}
	}
	counter := a.eventNoStorage[pk]
	return int(counter.Increase(1)), nil
}

func (a *CurrencyMemoryEventStorage) AddEvent(event *es.Event[currency.Request]) error {
	// eventStorage 와 pkGroupStorage 를 둘다 사용하므로
	// 때에 따라, storage 의 초기화가 필요할 수 있다.
	// 중복적인 초기화는 storage 를 날려버릴 수 있으므로 mutex lock 을 건다.

	// pk group storage 를 가져옴, 중복 할당 방지
	if _, ok := a.pkGroupStorage[event.PartitionKey]; !ok {
		pkLocker := a.getPkLocker(event.PartitionKey)
		pkLocker.Lock()
		defer pkLocker.Unlock()
		if _, ok = a.pkGroupStorage[event.PartitionKey]; !ok {
			a.pkGroupStorage[event.PartitionKey] = make([]es.EventId, 0)
		}
	}

	a.eventStorage[event.EventId] = *event
	a.pkGroupStorage[event.PartitionKey] = append(a.pkGroupStorage[event.PartitionKey], event.EventId)
	return nil
}

func (a *CurrencyMemoryEventStorage) GetEvent(id es.EventId) (*es.Event[currency.Request], error) {
	event := a.eventStorage[id]
	if event.EventNo == 0 { // event 가 nil 인 경우는 eventNo 가 초기값인 0이다
		return nil, nil
	}
	return &event, nil
}

func (a *CurrencyMemoryEventStorage) GetEvents(pk es.PartitionKey) ([]*es.Event[currency.Request], error) {
	// 이벤트 리스트를 조회할 때는 dirty read 를 방지하기 위해
	// Read Lock 을 건다

	locker := a.getPkLocker(pk)
	locker.RLock()
	defer locker.RUnlock()

	eventIds := a.pkGroupStorage[pk]
	ptrEvents := make([]*es.Event[currency.Request], len(eventIds))
	for i, eventId := range eventIds {
		event := a.eventStorage[eventId]
		ptrEvents[i] = &event
	}

	return ptrEvents, nil
}

func (a *CurrencyMemoryEventStorage) GetEventsAfterEventNo(pk es.PartitionKey, eventNo int) ([]*es.Event[currency.Request], error) {
	// event No 이상이 되는 이벤트 리스트를 조회할 때는 dirty read 를 방지하기 위해
	// Read Lock 을 건다

	locker := a.getPkLocker(pk)
	locker.RLock()
	defer locker.RUnlock()

	eventIds := a.pkGroupStorage[pk]
	ptrEvents := make([]*es.Event[currency.Request], 0)
	for _, eventId := range eventIds {
		if e := a.eventStorage[eventId]; e.EventNo > eventNo {
			ptrEvents = append(ptrEvents, &e)
		}
	}
	return ptrEvents, nil
}

func (a *CurrencyMemoryEventStorage) GetLastEvent(pk es.PartitionKey) (*es.Event[currency.Request], error) {
	// 쌓인 마지막 이벤트를 가져올 때, dirty read 를 방지하기 위해
	// Read lock 을 건다
	locker := a.getPkLocker(pk)
	locker.RLock()
	defer locker.RUnlock()

	group, ok := a.pkGroupStorage[pk]
	if !ok {
		return nil, nil
	}
	if len(group) == 0 {
		return nil, nil
	}
	eventId := group[len(group)-1]
	event := a.eventStorage[eventId]
	return &event, nil
}

func (a *CurrencyMemoryEventStorage) GetLock(pk es.PartitionKey) (bool, error) {
	// 현재 storage 에 저장된 lock 값을 읽을 때 dirty read 를 방지하기 위해
	// Read lock 을 이용한다.

	pkLocker := a.getPkLocker(pk)
	pkLocker.RLock()
	defer pkLocker.RUnlock()

	return a.eventLockStorage[pk], nil
}

func (a *CurrencyMemoryEventStorage) Lock(pk es.PartitionKey) (already bool, err error) {
	// lock 을 거는 이벤트와 unlock 을 거는 이벤트가 분리된 경우, lock event 가 사용한다.
	// eventLockStorage 에 lock 여부를 저장할 때도, atomic 하게 작동해야 하므로 pkLockers 를 이용한다.

	pkLocker := a.getPkLocker(pk)
	pkLocker.Lock()
	defer pkLocker.Unlock()

	lockValue, ok := a.eventLockStorage[pk]
	if !ok {
		a.eventLockStorage[pk] = true
		return false, nil
	}
	if lockValue {
		return true, nil
	}
	a.eventLockStorage[pk] = true
	return false, nil
}

func (a *CurrencyMemoryEventStorage) Unlock(pk es.PartitionKey) (already bool, err error) {
	// lock 을 거는 이벤트와 unlock 을 거는 이벤트가 분리된 경우, unlock event 가 사용한다.
	// eventLockStorage 에 lock 여부를 저장할 때도, atomic 하게 작동해야 하므로 pkLockers 를 이용한다.

	pkLocker := a.getPkLocker(pk)
	pkLocker.Lock()
	defer pkLocker.Unlock()

	lockValue, ok := a.eventLockStorage[pk]
	if !ok {
		a.eventLockStorage[pk] = false
		return false, nil
	}
	if !lockValue {
		return true, nil
	}
	a.eventLockStorage[pk] = false
	return false, nil
}

type CurrencySnapshotStorage struct {
	pkSnapshotStorage map[es.PartitionKey]es.State[currency.State, currency.Request]
	pkLockers         map[es.PartitionKey]*sync.RWMutex
	ssLocker          sync.Mutex
}

func NewCurrencySnapshotStorage() es.StateSnapshotStorage[currency.State, currency.Request] {
	return &CurrencySnapshotStorage{
		pkSnapshotStorage: make(map[es.PartitionKey]es.State[currency.State, currency.Request]),
		pkLockers:         make(map[es.PartitionKey]*sync.RWMutex),
	}
}

func (a *CurrencySnapshotStorage) getPkLocker(pk es.PartitionKey) *sync.RWMutex {
	// pk locker 중복 할당 방지
	if _, ok := a.pkLockers[pk]; !ok {
		a.ssLocker.Lock()
		defer a.ssLocker.Unlock()
		if _, ok = a.pkLockers[pk]; !ok {
			a.pkLockers[pk] = &sync.RWMutex{}
		}
	}
	// locker ptr 을 리턴해야 copy 이슈로 lock 이 걸리지 않는 이슈가 발생하지 않음
	return a.pkLockers[pk]
}

func (a *CurrencySnapshotStorage) SaveSnapshot(pk es.PartitionKey, state *es.State[currency.State, currency.Request]) error {
	pkLocker := a.getPkLocker(pk)
	pkLocker.Lock()
	defer pkLocker.Unlock()

	a.pkSnapshotStorage[pk] = *state
	return nil
}

func (a *CurrencySnapshotStorage) GetSnapshot(pk es.PartitionKey) (state *es.State[currency.State, currency.Request], err error) {
	pkLocker := a.getPkLocker(pk)
	pkLocker.Lock()
	defer pkLocker.Unlock()

	snapshot, ok := a.pkSnapshotStorage[pk]
	if !ok {
		return nil, nil
	}
	return &snapshot, nil
}
