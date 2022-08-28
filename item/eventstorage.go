package item

import (
	"errors"
	"sync"
)

type EventStorage interface {
	SetEvent(event Event) error
	GetEvent(eventId string) (*Event, error)
	GetEvents(key PartitionKey) (*EventList, error)
}

type MockEventStorage struct {
	eventStorage     map[string]Event
	partitionStorage map[PartitionKey][]Event
	sync.Mutex
}

func NewMockEventStorage() EventStorage {
	return &MockEventStorage{
		eventStorage:     make(map[string]Event),
		partitionStorage: make(map[PartitionKey][]Event),
	}
}

func (m *MockEventStorage) initPartition(key PartitionKey) {
	if _, ok := m.partitionStorage[key]; !ok {
		m.Lock()
		defer m.Unlock()
		if _, ok := m.partitionStorage[key]; !ok {
			m.partitionStorage[key] = make([]Event, 0)
		}
	}
}

func (m *MockEventStorage) SetEvent(event Event) error {
	m.initPartition(event.PartitionKey)
	m.eventStorage[event.Id] = event
	m.partitionStorage[event.PartitionKey] = append(m.partitionStorage[event.PartitionKey], event)
	return nil
}

func (m *MockEventStorage) GetEvent(eventId string) (*Event, error) {
	event, ok := m.eventStorage[eventId]
	if ok {
		return &event, nil
	}
	return nil, errors.New("not exists event")
}

func (m *MockEventStorage) GetEvents(key PartitionKey) (*EventList, error) {
	m.initPartition(key)
	events := m.partitionStorage[key]
	if events == nil || len(events) == 0 {
		return nil, errors.New("not exists events")
	}
	return (*EventList)(&events), nil
}
