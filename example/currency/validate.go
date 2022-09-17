package currency

import (
	es "eventsourcing"
)

/**
Event 를 받아들일지 말지 결정하는 Validate func 를 여기에 구현한다.
*/

var (
	Validator *es.Validator[State, Request]
	_         es.Validate[State, Request] = NoBurned
)

func init() {
	// Validator 구성
	// 등록되지 않은 Event 는 Validate 가 없는 액션
	Validator = es.NewValidator[State, Request]()
	Validator.SetValidates(AddAmountEvent, NoBurned)
	Validator.SetValidates(MinusAmountEvent, NoBurned)
	Validator.SetValidates(ChangeStatusEvent, NoBurned)
	Validator.SetValidates(ChangeValueEvent, NoBurned)
	Validator.SetValidates(ChangeValueV2Event, NoBurned)
	Validator.SetValidates(BurnEvent, NoBurned)
}

func NoBurned(latest *es.Event[Request], snapshot *es.State[State, Request]) {
	if snapshot != nil {
		if snapshot.State().Status == BURNED {
			panic("status is burned")
		}
	}
	if snapshot == nil {
		if latest.EventType.String() == BurnEvent.String() {
			panic("status is burned")
		}
	}
	if latest.EventNo > snapshot.State().LastEvent.EventNo {
		if latest.EventType.String() == BurnEvent.String() {
			panic("status is burned")
		}
	}
}
