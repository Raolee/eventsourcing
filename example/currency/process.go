package currency

import (
	es "eventsourcing"
)

/**
State 수정/변경 하는 Process 를 구현한다.
*/

var (
	Processor *es.Processor[State, Request]
	_         es.Process[State, Request] = CreateCurrencyState
	_         es.Process[State, Request] = AddAmount
	_         es.Process[State, Request] = MinusAmount
	_         es.Process[State, Request] = ChangeStatus
	_         es.Process[State, Request] = ChangeValue
	_         es.Process[State, Request] = ChangeValueV2
	_         es.Process[State, Request] = Burn
)

func init() {
	// Processor 구성
	// 정의한 모든 Event 는 Process 와 매핑되어야 함
	Processor = es.NewProcessor[State, Request]()
	Processor.SetProcess(CreateAmountStateEvent, CreateCurrencyState)
	Processor.SetProcess(AddAmountEvent, AddAmount)
	Processor.SetProcess(MinusAmountEvent, MinusAmount)
	Processor.SetProcess(ChangeStatusEvent, ChangeStatus)
	Processor.SetProcess(ChangeValueEvent, ChangeValue)
	Processor.SetProcess(ChangeValueV2Event, ChangeValueV2) // V2 의 Cmd 를 따로 매핑
	Processor.SetProcess(BurnEvent, Burn)
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
