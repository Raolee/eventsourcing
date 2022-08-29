package item

import "testing"

func TestEventReplay(t *testing.T) {
	TestMockEventStorage(t) // MockEventStorage 에 쌓아둔 이벤트를 이용

	events, err := mockStorage.GetEvents(PartitionKey(partitionKey))
	if err != nil {
		t.Error(err)
	}

	replayedState, err := ReplayEventsWithoutState(events...)
	if err != nil {
		t.Error(err)
	}
	t.Log(replayedState)
}
