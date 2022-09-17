package storage

import (
	es "eventsourcing"
	"eventsourcing/example/currency"
	"sync"
)

type CurrencyMemoryEventStorage struct {
	eventNoStorage   map[es.PartitionKey]*int
	pkGroupStorage   map[es.PartitionKey][]es.EventId
	eventStorage     map[es.EventId]*es.Event[currency.Request]
	eventLockStorage map[es.PartitionKey]bool
	rwLocker         sync.RWMutex
}

func NewCurrencyEventStorage() *CurrencyMemoryEventStorage {
	return &CurrencyMemoryEventStorage{
		eventNoStorage:   make(map[es.PartitionKey]*int),
		pkGroupStorage:   make(map[es.PartitionKey][]es.EventId),
		eventStorage:     make(map[es.EventId]*es.Event[currency.Request]),
		eventLockStorage: make(map[es.PartitionKey]bool),
	}
}

func (a CurrencyMemoryEventStorage) IncreaseEventNo(pk es.PartitionKey) (eventNo int, err error) {
	a.rwLocker.Lock()
	defer a.rwLocker.Unlock()
	no, ok := a.eventNoStorage[pk]
	if !ok {
		a.eventNoStorage[pk] = 1 // TODO : 09.17 여기부터 구현
		no = 1
	}
}

func (a CurrencyMemoryEventStorage) getCurrentEventNo(pk es.PartitionKey) (eventNo int) {
	a.rwLocker.RLocker()
	defer a.rwLocker.RUnlock()
	return a.eventNoStorage[pk]
}

func (a CurrencyMemoryEventStorage) AddEvent(event *es.Event[currency.Request]) error {
	//TODO implement me
	panic("implement me")
}

func (a CurrencyMemoryEventStorage) GetEvent(id es.EventId) (*es.Event[currency.Request], error) {
	panic("implement me")
}

func (a CurrencyMemoryEventStorage) GetEvents(pk es.PartitionKey) ([]*es.Event[currency.Request], error) {
	//TODO implement me
	panic("implement me")
}

func (a CurrencyMemoryEventStorage) GetEventsAfterEventNo(pk es.PartitionKey, eventNo int) ([]*es.Event[currency.Request], error) {
	//TODO implement me
	panic("implement me")
}

func (a CurrencyMemoryEventStorage) GetLastEvent(pk es.PartitionKey) (*es.Event[currency.Request], error) {
	//TODO implement me
	panic("implement me")
}

func (a CurrencyMemoryEventStorage) GetLock(pk es.PartitionKey) (bool, error) {
	//TODO implement me
	panic("implement me")
}

func (a CurrencyMemoryEventStorage) Lock(pk es.PartitionKey) (already bool, err error) {
	//TODO implement me
	panic("implement me")
}

func (a CurrencyMemoryEventStorage) Unlock(pk es.PartitionKey) (already bool, err error) {
	//TODO implement me
	panic("implement me")
}

type CurrencySnapshotStorage struct {
}

func NewCurrencySnapshotStorage() *CurrencySnapshotStorage {
	return &CurrencySnapshotStorage{}
}

func (a CurrencySnapshotStorage) SaveSnapshot(pk es.PartitionKey, state *es.State[currency.State, currency.Request]) error {
	//TODO implement me
	panic("implement me")
}

func (a CurrencySnapshotStorage) GetSnapshot(pk es.PartitionKey) (state *es.State[currency.State, currency.Request], err error) {
	//TODO implement me
	panic("implement me")
}
