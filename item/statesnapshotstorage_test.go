package item

import "testing"

func TestNewMockStateSnapshotStorage(t *testing.T) {
	TestMockEventStorage(t) // MockEventStorage 를 이용함

	events, err := mockStorage.GetEvents(partitionKey)
	if err != nil {
		t.Error(err)
	}

	storage := NewMockStateSnapshotStorage()

	err = storage.CreateSnapshot(events...)
	if err != nil {
		t.Error(err)
	}
	exists, snapshot, lastEvent := storage.GetSnapshot(partitionKey)
	t.Log(exists)
	t.Log(snapshot)
	t.Log(lastEvent)

	event := NewEvent(ChangeItemOwnerEvent, "v1", partitionKey, NewRequests(&Owner{
		AccountKey: "raol2",
	}))
	err = storage.UpdateSnapshot(event)
	if err != nil {
		t.Error(err)
	}

	exists, snapshot, lastEvent = storage.GetSnapshot(partitionKey)
	t.Log(exists)
	t.Log(snapshot)
	t.Log(lastEvent)
}
