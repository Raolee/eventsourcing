package eventsourcing

import (
	"errors"
	"sync"
)

type ItemEventStorage interface {
	SetItemEvent(event ItemEvent) error
	GetItemEvent(eventId string) (*ItemEvent, error)
	GetItemEvents(key PartitionKey) (*[]ItemEvent, error)
}

type MockItemEventStorage struct {
	eventStorage     map[string]ItemEvent
	partitionStorage map[PartitionKey][]ItemEvent
	sync.Mutex
}

func NewMockItemEventStorage() ItemEventStorage {
	return &MockItemEventStorage{
		eventStorage:     make(map[string]ItemEvent),
		partitionStorage: make(map[PartitionKey][]ItemEvent),
	}
}

func (m *MockItemEventStorage) initPartition(key PartitionKey) {
	if _, ok := m.partitionStorage[key]; !ok {
		m.Lock()
		defer m.Unlock()
		if _, ok := m.partitionStorage[key]; !ok {
			m.partitionStorage[key] = make([]ItemEvent, 0)
		}
	}
}

func (m *MockItemEventStorage) SetItemEvent(event ItemEvent) error {
	m.initPartition(event.PartitionKey)
	m.eventStorage[event.Id] = event
	m.partitionStorage[event.PartitionKey] = append(m.partitionStorage[event.PartitionKey], event)
	return nil
}

func (m *MockItemEventStorage) GetItemEvent(eventId string) (*ItemEvent, error) {
	event, ok := m.eventStorage[eventId]
	if ok {
		return &event, nil
	}
	return nil, errors.New("not exists event")
}

func (m *MockItemEventStorage) GetItemEvents(key PartitionKey) (*[]ItemEvent, error) {
	m.initPartition(key)
	events := m.partitionStorage[key]
	if events == nil || len(events) == 0 {
		return nil, errors.New("not exists events")
	}
	return &events, nil
}
