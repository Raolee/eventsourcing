package storage

import (
	"errors"
	"eventsourcing"
	"eventsourcing/item"
	"sync"
)

// MockEventStorage | Mock-up Event 저장소
type MockEventStorage struct {
	eventStorage     map[eventsourcing.EventId]*eventsourcing.Event        // 정렬은 되지 않은 event storage
	partitionStorage map[eventsourcing.PartitionKey][]*eventsourcing.Event // partition key 별로 모아둔 event storage
	sync.Mutex
}

func NewMockEventStorage() eventsourcing.EventStorage {
	return &MockEventStorage{
		eventStorage:     make(map[eventsourcing.EventId]*eventsourcing.Event),
		partitionStorage: make(map[eventsourcing.PartitionKey][]*eventsourcing.Event),
	}
}

func (m *MockEventStorage) initPartition(key item.PartitionKey) {
	if _, ok := m.partitionStorage[key]; !ok {
		m.Lock()
		defer m.Unlock()
		if _, ok := m.partitionStorage[key]; !ok {
			m.partitionStorage[key] = make([]*item.Event, 0)
		}
	}
}

func (m *MockEventStorage) AddEvent(event eventsourcing.Event) error {
	m.initPartition(event.PartitionKey)
	m.eventStorage[event.Id] = event
	m.partitionStorage[event.PartitionKey] = append(m.partitionStorage[eventsourcing.PartitionKey], event)
	return nil
}

func (m *MockEventStorage) GetEvent(id eventsourcing.EventId) (eventsourcing.Event, error) {
	event, ok := m.eventStorage[id]
	if ok {
		return event, nil
	}
	return nil, errors.New("not exists event")
}

func (m *MockEventStorage) GetEvents(key eventsourcing.PartitionKey) ([]eventsourcing.Event, error) {
	m.initPartition(key)
	events := m.partitionStorage[key]
	if events == nil || len(events) == 0 {
		return nil, errors.New("not exists events")
	}
	return events, nil
}
