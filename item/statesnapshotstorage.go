package item

import "errors"

type StateSnapshotStorage interface {
	CreateSnapshot(events ...*Event) error
	UpdateSnapshot(event *Event) error
	GetSnapshot(key PartitionKey) (bool, *State, *Event)
}

type MockStateSnapshotStorage struct {
	stateMap          map[PartitionKey]*State
	snapShotLastEvent map[PartitionKey]*Event
}

func NewMockStateSnapshotStorage() StateSnapshotStorage {
	return &MockStateSnapshotStorage{
		stateMap:          make(map[PartitionKey]*State),
		snapShotLastEvent: make(map[PartitionKey]*Event),
	}
}

func (m *MockStateSnapshotStorage) CreateSnapshot(events ...*Event) error {
	state, err := ReplayEventsWithoutState(events...)
	if err != nil {
		return err
	}
	m.stateMap[state.PartitionKey()] = state                          // replay 된 state 를 저장
	m.snapShotLastEvent[state.PartitionKey()] = events[len(events)-1] // snapshot 을 구성하는 마지막 event 를 저장
	return nil
}

func (m *MockStateSnapshotStorage) UpdateSnapshot(event *Event) error {
	_, ok := m.stateMap[event.PartitionKey]
	if !ok {
		return errors.New("state is nil")
	}
	state, err := ReplayEventsWithState(m.stateMap[event.PartitionKey], event)
	if err != nil {
		return err
	}
	m.stateMap[state.PartitionKey()] = state
	m.snapShotLastEvent[state.PartitionKey()] = event
	return nil
}

func (m *MockStateSnapshotStorage) GetSnapshot(key PartitionKey) (bool, *State, *Event) {
	state, ok := m.stateMap[key]
	lastEvent := m.snapShotLastEvent[key]
	return ok, state, lastEvent
}
