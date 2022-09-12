package example

import es "eventsourcing"

var (
	ExampleCommander *es.Commander[AmountState, ExampleRequest]
	ExampleValidator *es.Validator[AmountState, ExampleRequest]
)

func init() {
	// Commander 구성
	// 정의한 모든 Event 는 Command 와 매핑되어야 함
	ExampleCommander = es.NewCommander[AmountState, ExampleRequest]()
	ExampleCommander.SetCommand(CreateAmountStateEvent, CreateAmountState)
	ExampleCommander.SetCommand(ModifyAmountEvent, ModifyAmount)
	ExampleCommander.SetCommand(ChangeStatusEvent, ChangeStatus)
	ExampleCommander.SetCommand(ChangeValueEvent, ChangeValue)
	ExampleCommander.SetCommand(ChangeValueV2Event, ChangeValueV2) // V2 의 Cmd 를 따로 매핑
	ExampleCommander.SetCommand(LockEvent, Lock)
	ExampleCommander.SetCommand(UnlockEvent, Unlock)
	ExampleCommander.SetCommand(BurnEvent, Burn)

	// Validator 구성
	// 등록되지 않은 Event 는 Validate 가 없는 액션
	ExampleValidator = es.NewValidator[AmountState, ExampleRequest]()
	ExampleValidator.SetValidates(ModifyAmountEvent, NoBurned)
	ExampleValidator.SetValidates(ChangeStatusEvent, NoBurned, NoLock)
	ExampleValidator.SetValidates(ChangeValueEvent, NoBurned, NoLock)
	ExampleValidator.SetValidates(ChangeValueV2Event, NoBurned)
	ExampleValidator.SetValidates(LockEvent, NoBurned)
	ExampleValidator.SetValidates(UnlockEvent, NoBurned)
	ExampleValidator.SetValidates(BurnEvent, NoBurned)
}
