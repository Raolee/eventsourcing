package item

import (
	"encoding/json"
	"eventsourcing"
	"eventsourcing/item/command"
	"eventsourcing/item/model"
	"time"
)

var (
	_ eventsourcing.Event = &Event{}
)

const (
	CreateEvent   = eventsourcing.EventName("create")
	SaveDataEvent = eventsourcing.EventName("save_data")
	RemoveEvent   = eventsourcing.EventName("remove")

	RequestMintingEvent = eventsourcing.EventName("request_minting")
	FailedMintingEvent  = eventsourcing.EventName("failed_minting")
	SuccessMintingEvent = eventsourcing.EventName("success_minting")

	RegisterMarketEvent = eventsourcing.EventName("register_market")
	CancelTradingEvent  = eventsourcing.EventName("cancel_trade")

	ChangeOwnerEvent          = eventsourcing.EventName("change_owner")
	ChangeCharacterOwnerEvent = eventsourcing.EventName("change_char_owner")

	RequestEnhanceEvent = eventsourcing.EventName("request_enhance")
	FailedEnhanceEvent  = eventsourcing.EventName("failed_enhance")
	SuccessEnhanceEvent = eventsourcing.EventName("success_enhance")

	BindToCharEvent          = eventsourcing.EventName("bind_to_char")
	FailedBindToCharEvent    = eventsourcing.EventName("failed_bind_to_char")
	SuccessBindToCharEvent   = eventsourcing.EventName("success_bind_to_char")
	UnbindToCharEvent        = eventsourcing.EventName("unbind_to_char")
	FailedUnbindToCharEvent  = eventsourcing.EventName("failed_unbind_to_char")
	SuccessUnbindToCharEvent = eventsourcing.EventName("success_unbind_to_char")

	RequestCollectEvent = eventsourcing.EventName("request_collect")
	FailedCollectEvent  = eventsourcing.EventName("failed_collect")
	SuccessCollectEvent = eventsourcing.EventName("success_collect")

	RequestBurnEvent = eventsourcing.EventName("request_burn")
	FailedBurnEvent  = eventsourcing.EventName("failed_burn")
	SuccessBurnEvent = eventsourcing.EventName("success_burn")
)

var (
	// EventCommandMap | EventName 과 Command 를 연결
	EventCommandMap = map[eventsourcing.EventName]eventsourcing.Command{
		CreateEvent:              command.Command.Create,
		SaveDataEvent:            command.Command.SaveData,
		RemoveEvent:              command.Command.Remove,
		RequestMintingEvent:      command.Command.RequestMinting,
		FailedMintingEvent:       command.Command.FailedMinting,
		SuccessMintingEvent:      command.Command.SuccessMinting,
		RegisterMarketEvent:      command.Command.RegisterMarket,
		CancelTradingEvent:       command.Command.CancelTrace,
		ChangeOwnerEvent:         command.Command.ChangeOwner,
		RequestEnhanceEvent:      command.Command.RequestEnhance,
		FailedEnhanceEvent:       command.Command.FailedEnhance,
		SuccessEnhanceEvent:      command.Command.SuccessEnhance,
		BindToCharEvent:          command.Command.RequestBindToCharacter,
		FailedBindToCharEvent:    command.Command.FailedBindToCharacter,
		SuccessBindToCharEvent:   command.Command.SuccessBindToCharacter,
		UnbindToCharEvent:        command.Command.RequestUnbindToCharacter,
		FailedUnbindToCharEvent:  command.Command.FailedUnbindToCharacter,
		SuccessUnbindToCharEvent: command.Command.SuccessUnbindToCharacter,
		RequestCollectEvent:      command.Command.RequestCollect,
		FailedCollectEvent:       command.Command.FailedCollect,
		SuccessCollectEvent:      command.Command.SuccessCollect,
		RequestBurnEvent:         command.Command.RequestBurn,
		FailedBurnEvent:          command.Command.FailedBurn,
		SuccessBurnEvent:         command.Command.SuccessBurn,
	}
	EventValidatorMap = map[EventName]func(validate Validate, event *Event) (bool, error){
		CreateEvent:         Validate.ValidateCreateItem,
		SaveDataEvent:       Validate.ValidateSaveItemData,
		RemoveEvent:         Validate.ValidateRemoveItem,
		RequestMintingEvent: Validate.ValidateMintingItemRequest,
		FailedMintingEvent:  Validate.ValidateMintingItemFailure,
		SuccessMintingEvent: Validate.ValidateMintingItemSuccess,
		RegisterMarketEvent: Validate.ValidateRegisterMarketItem,
		CancelTradingEvent:  Validate.ValidateCancelMarketItem,
		ChangeOwnerEvent:    Validate.ValidateChangeItemOwner,
		RequestEnhanceEvent: Validate.ValidateEnhanceItemRequest,
		FailedEnhanceEvent:  Validate.ValidateEnhanceItemFailure,
		SuccessEnhanceEvent: Validate.ValidateEnhanceItemSuccess,
		RequestBurnEvent:    Validate.ValidateBurnItemRequest,
		FailedBurnEvent:     Validate.ValidateBurnItemFailure,
		SuccessBurnEvent:    Validate.ValidateBurnItemSuccess,
	}
)

type Event struct {
	Id           eventsourcing.EventId      `json:"id"` // event 의 고유 아이디, sorted 해야 한다.
	Domain       eventsourcing.Domain       `json:"domain"`
	PartitionKey eventsourcing.PartitionKey `json:"partitionKey"` // event 를 묶어줄 key, storage 에서는 partition key 로 쓰인다.
	Name         eventsourcing.EventName    `json:"name"`         // event 명칭, command 와 mapping 되어 있다.
	Version      eventsourcing.EventVersion `json:"version"`      // event 의 버전, Name 별로 versioning 한다.
	EventAt      time.Time                  `json:"eventAt"`      // event 발생 시간
	Requests     *model.RequestBody         `json:"requests"`     // event 의 body, 실제 command 를 이 body 를 활용함
}

func NewEvent(
	id eventsourcing.EventId,
	name eventsourcing.EventName,
	version eventsourcing.EventVersion,
	partitionKey eventsourcing.PartitionKey,
	reqs *model.RequestBody,
) *Event {
	return &Event{
		Id:           id,
		Name:         name,
		Version:      version,
		EventAt:      time.Now(),
		PartitionKey: partitionKey,
		Requests:     reqs,
	}
}

func (e *Event) GetId() eventsourcing.EventId {
	return e.Id
}

func (e *Event) GetDomain() eventsourcing.Domain {
	return e.Domain
}

func (e *Event) GetPartitionKey() eventsourcing.PartitionKey {
	return e.PartitionKey
}

func (e *Event) GetName() eventsourcing.EventName {
	return e.Name
}

func (e *Event) GetVersion() eventsourcing.EventVersion {
	return e.Version
}

func (e *Event) GetEventNo() int {
	//TODO implement me
	panic("implement me")
}

func (e *Event) GetEventAt() time.Time {
	return e.EventAt
}

// EventList | event list 정의
type EventList []*Event

func (e EventList) String() string {
	if e == nil {
		return "[]"
	}
	bytes, err := json.MarshalIndent(e, "", "    ")
	if err != nil {
		return "[]"
	}
	return string(bytes)
}

// Iterate | 채널 방식으로 for range 를 사용할 수 있게 제공, 다른 방식으로 next() 를 구현하는 방법도 있음
func (e EventList) Iterate() <-chan *Event {
	ch := make(chan *Event)
	go func() {
		for _, event := range ([]*Event)(e) {
			ch <- event
		}
		close(ch)
	}()
	return ch
}
