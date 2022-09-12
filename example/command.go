package example

import (
	es "eventsourcing"
)

/**
AmountState 수정/변경 하는 Command 를 구현한다.
*/

var (
	_ es.Command[AmountState, ExampleRequest] = CreateAmountState
	_ es.Command[AmountState, ExampleRequest] = ModifyAmount
	_ es.Command[AmountState, ExampleRequest] = ChangeStatus
	_ es.Command[AmountState, ExampleRequest] = ChangeValue
	_ es.Command[AmountState, ExampleRequest] = ChangeValueV2
	_ es.Command[AmountState, ExampleRequest] = Lock
	_ es.Command[AmountState, ExampleRequest] = Unlock
)

func CreateAmountState(s *es.State[AmountState, ExampleRequest], e *es.Event[ExampleRequest]) *es.State[AmountState, ExampleRequest] {
	s = es.NewState[AmountState, ExampleRequest](&AmountState{
		PartitionKey: e.PartitionKey,
		Amount:       0,
		Status:       NOTHING,
		Value:        nil,
		Lock:         false,
		LastEvent:    e,
	})
	return s
}

func ModifyAmount(s *es.State[AmountState, ExampleRequest], e *es.Event[ExampleRequest]) *es.State[AmountState, ExampleRequest] {
	s.State().Amount = s.State().Amount + e.Request.Amount
	s.State().LastEvent = e
	return s
}

func ChangeStatus(s *es.State[AmountState, ExampleRequest], e *es.Event[ExampleRequest]) *es.State[AmountState, ExampleRequest] {
	s.State().Status = *e.Request.Status
	s.State().LastEvent = e
	return s
}

func ChangeValue(s *es.State[AmountState, ExampleRequest], e *es.Event[ExampleRequest]) *es.State[AmountState, ExampleRequest] {
	s.State().Value = e.Request.Value
	s.State().LastEvent = e
	return s
}

func ChangeValueV2(s *es.State[AmountState, ExampleRequest], e *es.Event[ExampleRequest]) *es.State[AmountState, ExampleRequest] {
	s.State().Value = e.Request.Value
	s.State().Lock = false // V2 는 unlock 상태로 만드는 액션까지 같이 들어간다
	s.State().LastEvent = e
	return s
}

func Lock(s *es.State[AmountState, ExampleRequest], e *es.Event[ExampleRequest]) *es.State[AmountState, ExampleRequest] {
	s.State().Lock = true
	s.State().LastEvent = e
	return s
}

func Unlock(s *es.State[AmountState, ExampleRequest], e *es.Event[ExampleRequest]) *es.State[AmountState, ExampleRequest] {
	s.State().Lock = false
	s.State().LastEvent = e
	return s
}

func Burn(s *es.State[AmountState, ExampleRequest], e *es.Event[ExampleRequest]) *es.State[AmountState, ExampleRequest] {
	s.State().Status = BURNED
	s.State().Lock = true
	s.State().LastEvent = e
	return s
}
