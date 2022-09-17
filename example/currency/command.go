package currency

import (
	es "eventsourcing"
)

/**
State 수정/변경 하는 Command 를 구현한다.
*/

var (
	Commander *es.Commander[State, Request]
	_         es.Command[State, Request] = CreateCurrencyState
	_         es.Command[State, Request] = AddAmount
	_         es.Command[State, Request] = MinusAmount
	_         es.Command[State, Request] = ChangeStatus
	_         es.Command[State, Request] = ChangeValue
	_         es.Command[State, Request] = ChangeValueV2
	_         es.Command[State, Request] = Burn
)

func init() {
	// Commander 구성
	// 정의한 모든 Event 는 Command 와 매핑되어야 함
	Commander = es.NewCommander[State, Request]()
	Commander.SetCommand(CreateAmountStateEvent, CreateCurrencyState)
	Commander.SetCommand(AddAmountEvent, AddAmount)
	Commander.SetCommand(MinusAmountEvent, MinusAmount)
	Commander.SetCommand(ChangeStatusEvent, ChangeStatus)
	Commander.SetCommand(ChangeValueEvent, ChangeValue)
	Commander.SetCommand(ChangeValueV2Event, ChangeValueV2) // V2 의 Cmd 를 따로 매핑
	Commander.SetCommand(BurnEvent, Burn)
}

func CreateCurrencyState(s *es.State[State, Request], e *es.Event[Request]) *es.State[State, Request] {
	s = es.NewState[State, Request](&State{
		PartitionKey: e.PartitionKey,
		Amount:       0,
		Status:       NOTHING,
		Value:        nil,
		LastEvent:    e,
	})
	return s
}

func AddAmount(s *es.State[State, Request], e *es.Event[Request]) *es.State[State, Request] {
	s.State().Amount = s.State().Amount + e.Request.Amount
	s.State().LastEvent = e
	return s
}

func MinusAmount(s *es.State[State, Request], e *es.Event[Request]) *es.State[State, Request] {
	s.State().Amount = s.State().Amount - e.Request.Amount
	s.State().LastEvent = e
	return s
}

func ChangeStatus(s *es.State[State, Request], e *es.Event[Request]) *es.State[State, Request] {
	s.State().Status = *e.Request.Status
	s.State().LastEvent = e
	return s
}

func ChangeValue(s *es.State[State, Request], e *es.Event[Request]) *es.State[State, Request] {
	s.State().Value = e.Request.Value
	s.State().LastEvent = e
	return s
}

func ChangeValueV2(s *es.State[State, Request], e *es.Event[Request]) *es.State[State, Request] {
	s.State().Amount = s.State().Amount + e.Request.Amount
	s.State().Value = e.Request.Value
	s.State().LastEvent = e
	return s
}

func Burn(s *es.State[State, Request], e *es.Event[Request]) *es.State[State, Request] {
	s.State().Status = BURNED
	s.State().LastEvent = e
	return s
}
