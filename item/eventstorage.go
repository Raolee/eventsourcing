package item

import (
	"errors"
	"sync"
)

// EventStorage | Event 저장소의 인터페이스
type EventStorage interface {
	SetEvent(event *Event) error                  // event 를 저장
	GetEvent(id EventId) (*Event, error)          // event 를 조회
	GetEvents(key PartitionKey) ([]*Event, error) // partition key 의 event list 를 조회
}

// MockEventStorage | Mock-up Event 저장소
type MockEventStorage struct {
	eventStorage     map[EventId]*Event        // 정렬은 되지 않은 event storage
	partitionStorage map[PartitionKey][]*Event // partition key 별로 모아둔 event storage
	sync.Mutex
}

func NewMockEventStorage() EventStorage {
	return &MockEventStorage{
		eventStorage:     make(map[EventId]*Event),
		partitionStorage: make(map[PartitionKey][]*Event),
	}
}

func (m *MockEventStorage) initPartition(key PartitionKey) {
	if _, ok := m.partitionStorage[key]; !ok {
		m.Lock()
		defer m.Unlock()
		if _, ok := m.partitionStorage[key]; !ok {
			m.partitionStorage[key] = make([]*Event, 0)
		}
	}
}

func (m *MockEventStorage) SetEvent(event *Event) error {
	m.initPartition(event.PartitionKey)
	m.eventStorage[event.Id] = event
	m.partitionStorage[event.PartitionKey] = append(m.partitionStorage[event.PartitionKey], event)
	return nil
}

func (m *MockEventStorage) GetEvent(id EventId) (*Event, error) {
	event, ok := m.eventStorage[id]
	if ok {
		return event, nil
	}
	return nil, errors.New("not exists event")
}

func (m *MockEventStorage) GetEvents(key PartitionKey) ([]*Event, error) {
	m.initPartition(key)
	events := m.partitionStorage[key]
	if events == nil || len(events) == 0 {
		return nil, errors.New("not exists events")
	}
	return events, nil
}
