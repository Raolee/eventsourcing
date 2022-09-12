package storage

import (
	"errors"
	"eventsourcing"
	"eventsourcing/item"
	"eventsourcing/item/command"
)

type MockStateSnapshotStorage struct {
	stateMap          map[eventsourcing.PartitionKey]command.State
	snapShotLastEvent map[eventsourcing.PartitionKey]eventsourcing.Event
}

func NewMockStateSnapshotStorage() StateSnapshotStorage {
	return &MockStateSnapshotStorage{
		stateMap:          make(map[eventsourcing.PartitionKey]*command.State),
		snapShotLastEvent: make(map[eventsourcing.PartitionKey]eventsourcing.Event),
	}
}

func (m *MockStateSnapshotStorage) CreateSnapshot(events ...item.Event) error {
	state, err := eventsourcing.ReplayEventsWithoutState(events...)
	if err != nil {
		return err
	}
	m.stateMap[state.PartitionKey()] = state                          // replay 된 state 를 저장
	m.snapShotLastEvent[state.PartitionKey()] = events[len(events)-1] // snapshot 을 구성하는 마지막 event 를 저장
	return nil
}

func (m *MockStateSnapshotStorage) UpdateSnapshot(event item.Event) error {
	_, ok := m.stateMap[event.PartitionKey]
	if !ok {
		return errors.New("state is nil")
	}
	state, err := eventsourcing.ReplayEventsWithState(m.stateMap[event.PartitionKey], event)
	if err != nil {
		return err
	}
	m.stateMap[state.PartitionKey()] = state
	m.snapShotLastEvent[state.PartitionKey()] = event
	return nil
}

func (m *MockStateSnapshotStorage) GetSnapshot(key eventsourcing.PartitionKey) (bool, *command.State, eventsourcing.Event) {
	state, ok := m.stateMap[key]
	lastEvent := m.snapShotLastEvent[key]
	return ok, state, lastEvent
}
