package example

import (
	es "eventsourcing"
)

/**
Event 를 받아들일지 말지 결정하는 Validate func 를 여기에 구현한다.
*/

var (
	_ es.Validate[AmountState, ExampleRequest] = NoBurned
	_ es.Validate[AmountState, ExampleRequest] = NoLock
)

func NoBurned(latest *es.Event[ExampleRequest], snapshot *es.State[AmountState, ExampleRequest]) {
	if snapshot != nil {
		if snapshot.State().Status == BURNED {
			panic("status is burned")
		}
	}
	if latest.EventNo > snapshot.State().LastEvent.EventNo {
		if latest.EventType.String() == BurnEvent.String() {
			panic("status is burned")
		}
	}
}

func NoLock(latest *es.Event[ExampleRequest], snapshot *es.State[AmountState, ExampleRequest]) {
	if latest.EventNo != snapshot.State().LastEvent.EventNo {
		// 스냅샷에 아직 이벤트가 반영되지 않은 상태의 유효성 검사
	}
	if snapshot.State().Lock {
		panic("no change value, because locked state")
	}
}
