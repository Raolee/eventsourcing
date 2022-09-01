package item

import (
	"encoding/json"
	"fmt"
	"github.com/rs/xid"
	"time"
)

type EventName string

const (
	CreateItemEvent   = EventName("createItem")
	SaveItemDataEvent = EventName("saveItemData")
	RemoveItemEvent   = EventName("removeItem")

	RequestMintingItemEvent = EventName("requestMintingItem")
	FailedMintingItemEvent  = EventName("failedMintingItem")
	SuccessMintingItemEvent = EventName("successMintingItem")

	RegisterMarketItemEvent = EventName("registerItemToMarket")
	CancelTradingItemEvent  = EventName("cancelTradingItem")

	ChangeItemOwnerEvent      = EventName("changeItemOwner")
	ChangeCharacterOwnerEvent = EventName("changeCharOwner")

	RequestEnhanceItemEvent = EventName("requestEnhanceItem")
	FailedEnhanceItemEvent  = EventName("failedEnhanceItem")
	SuccessEnhanceItemEvent = EventName("successEnhanceItem")

	BindItemToCharEvent        = EventName("bindItemToChar")
	FailedBindItemToCharEvent  = EventName("failedBindItemToChar")
	SuccessBindItemToCharEvent = EventName("successBindItemToChar")

	RequestBindItemToCatalogEvent = EventName("bindItemToCatalog")
	FailedBindItemToCatalogEvent  = EventName("failedBindItemToCatalog")
	SuccessBindItemToCatalogEvent = EventName("successBindItemToCatalog")

	RequestBurnItemEvent = EventName("requestBurnItem")
	FailedBurnItemEvent  = EventName("failedBurnItem")
	SuccessBurnItemEvent = EventName("successBurnItem")
)

var (
	// EventCommandMap | EventName 과 Command 를 연결
	EventCommandMap = map[EventName]func(command Command, event *Event) *State{
		CreateItemEvent:               Command.CreateItem,
		SaveItemDataEvent:             Command.SaveItemData,
		RemoveItemEvent:               Command.RemoveItem,
		RequestMintingItemEvent:       Command.ApplyRequestMintingItem,
		FailedMintingItemEvent:        Command.ApplyFailedMintingItem,
		SuccessMintingItemEvent:       Command.ApplySuccessMintingItem,
		RegisterMarketItemEvent:       Command.RegisterMarketItem,
		CancelTradingItemEvent:        Command.CancelMarketItem,
		ChangeItemOwnerEvent:          Command.ChangeItemOwner,
		RequestEnhanceItemEvent:       Command.ApplyRequestEnhanceItem,
		FailedEnhanceItemEvent:        Command.ApplyFailedEnhanceItem,
		SuccessEnhanceItemEvent:       Command.ApplySuccessEnhanceItem,
		BindItemToCharEvent:           Command.ApplyRequestBindItemToCharacter,
		FailedBindItemToCharEvent:     Command.ApplyFailedBindItemToCharacter,
		SuccessBindItemToCharEvent:    Command.ApplySuccessBindItemToCharacter,
		RequestBindItemToCatalogEvent: Command.ApplyRequestBindItemToCatalog,
		FailedBindItemToCatalogEvent:  Command.ApplyFailedBindItemToCatalog,
		SuccessBindItemToCatalogEvent: Command.ApplySuccessBindItemToCatalog,
		RequestBurnItemEvent:          Command.ApplyRequestBurnItem,
		FailedBurnItemEvent:           Command.ApplyFailedBurnItem,
		SuccessBurnItemEvent:          Command.ApplySuccessBurnItem,
	}
	EventValidatorMap = map[EventName]func(validate Validate, event *Event) (bool, error){
		CreateItemEvent:         Validate.ValidateCreateItem,
		SaveItemDataEvent:       Validate.ValidateSaveItemData,
		RemoveItemEvent:         Validate.ValidateRemoveItem,
		RequestMintingItemEvent: Validate.ValidateMintingItemRequest,
		FailedMintingItemEvent:  Validate.ValidateMintingItemFailure,
		SuccessMintingItemEvent: Validate.ValidateMintingItemSuccess,
		RegisterMarketItemEvent: Validate.ValidateRegisterMarketItem,
		CancelTradingItemEvent:  Validate.ValidateCancelMarketItem,
		ChangeItemOwnerEvent:    Validate.ValidateChangeItemOwner,
		RequestEnhanceItemEvent: Validate.ValidateEnhanceItemRequest,
		FailedEnhanceItemEvent:  Validate.ValidateEnhanceItemFailure,
		SuccessEnhanceItemEvent: Validate.ValidateEnhanceItemSuccess,
		RequestBurnItemEvent:    Validate.ValidateBurnItemRequest,
		FailedBurnItemEvent:     Validate.ValidateBurnItemFailure,
		SuccessBurnItemEvent:    Validate.ValidateBurnItemSuccess,
	}
)

// EventId | Event 의 고유 아이디, sorted 해야 한다. Event 는 이 type 을 key 로 삼아야 함
type EventId string

type Event struct {
	Id           EventId      `json:"id"`           // event 의 고유 아이디, sorted 해야 한다.
	Name         EventName    `json:"name"`         // event 명칭, command 와 mapping 되어 있다.
	Version      string       `json:"version"`      // event 의 버전, Name 별로 versioning 한다.
	EventAt      time.Time    `json:"eventAt"`      // event 발생 시간
	PartitionKey PartitionKey `json:"partitionKey"` // event 를 묶어줄 key, storage 에서는 partition key 로 쓰인다.
	Requests     *RequestBody `json:"requests"`     // event 의 body, 실제 command 를 이 body 를 활용함
}

func NewEvent(name EventName, version string, partitionKey PartitionKey, reqs *RequestBody) *Event {
	return &Event{
		Id:           EventId(xid.New().String()),
		Name:         name,
		Version:      version,
		EventAt:      time.Now(),
		PartitionKey: partitionKey,
		Requests:     reqs,
	}
}

func (e *Event) String() string {
	if e == nil {
		return "{}"
	}
	bytes, err := json.MarshalIndent(e, "", "    ")
	if err != nil {
		return "{}"
	}
	return string(bytes)
}

func (e *Event) EventName() string {
	return fmt.Sprintf("%s_%s", e.Name, e.Version)
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
