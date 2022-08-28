package item

import "testing"

func TestEventReplay(t *testing.T) {
	TestMockEventStorage(t)

	events, err := mockStorage.GetEvents(PartitionKey(assetKey))
	if err != nil {
		t.Error(err)
	}

	replayedState, err := EventReplay(events)
	if err != nil {
		t.Error(err)
	}
	t.Log(replayedState)
}
