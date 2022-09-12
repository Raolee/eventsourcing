package item

import (
	"eventsourcing"
	"eventsourcing/item/storage"
	"testing"
)

func TestEventReplay(t *testing.T) {
	storage.TestMockEventStorage(t) // MockEventStorage 에 쌓아둔 이벤트를 이용

	events, err := storage.mockStorage.GetEvents(PartitionKey(storage.partitionKey))
	if err != nil {
		t.Error(err)
	}

	replayedState, err := eventsourcing.ReplayEventsWithoutState(events...)
	if err != nil {
		t.Error(err)
	}
	t.Log(replayedState)
}
