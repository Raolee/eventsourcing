package item

import (
	"encoding/json"
	"fmt"
	"github.com/rs/xid"
	"time"
)

type EventName string

const (
	CreateItemEvent           = EventName("createItem")
	SaveItemDataEvent         = EventName("saveItemData")
	RemoveItemEvent           = EventName("removeItem")
	MintingItemRequestEvent   = EventName("mintingItemRequest")
	MintingItemFailureEvent   = EventName("mintingItemFailure")
	MintingItemSuccessEvent   = EventName("mintingItemSuccess")
	RegisterMarketItemEvent   = EventName("registerMarketItem")
	CancelTradingItemEvent    = EventName("cancelTradingItem")
	ChangeItemOwnerEvent      = EventName("successTradingItem")
	EnhancingItemRequestEvent = EventName("enhancingItemRequest")
	EnhancingItemFailureEvent = EventName("enhancingItemFailure")
	EnhancingItemSuccessEvent = EventName("enhancingItemSuccess")
	BurningItemRequestEvent   = EventName("burningItemRequest")
	BurningItemFailureEvent   = EventName("burningItemFailure")
	BurningItemSuccessEvent   = EventName("burningItemSuccess")
)

var (
	// EventCommandMap | EventName 과 Command 를 연결
	EventCommandMap = map[EventName]func(command Command, event *Event) *State{
		CreateItemEvent:           Command.CreateItem,
		SaveItemDataEvent:         Command.SaveItemData,
		RemoveItemEvent:           Command.RemoveItem,
		MintingItemRequestEvent:   Command.MintingItemRequest,
		MintingItemFailureEvent:   Command.MintingItemFailure,
		MintingItemSuccessEvent:   Command.MintingItemSuccess,
		RegisterMarketItemEvent:   Command.RegisterMarketItem,
		CancelTradingItemEvent:    Command.CancelMarketItem,
		ChangeItemOwnerEvent:      Command.ChangeItemOwner,
		EnhancingItemRequestEvent: Command.EnhancingItemRequest,
		EnhancingItemFailureEvent: Command.EnhancingItemFailure,
		EnhancingItemSuccessEvent: Command.EnhancingItemSuccess,
		BurningItemRequestEvent:   Command.BurningItemRequest,
		BurningItemFailureEvent:   Command.BurningItemFailure,
		BurningItemSuccessEvent:   Command.BurningItemSuccess,
	}
	EventValidatorMap = map[EventName]func(validate Validate, event *Event) (bool, error){
		CreateItemEvent:           Validate.ValidateCreateItem,
		SaveItemDataEvent:         Validate.ValidateSaveItemData,
		RemoveItemEvent:           Validate.ValidateRemoveItem,
		MintingItemRequestEvent:   Validate.ValidateMintingItemRequest,
		MintingItemFailureEvent:   Validate.ValidateMintingItemFailure,
		MintingItemSuccessEvent:   Validate.ValidateMintingItemSuccess,
		RegisterMarketItemEvent:   Validate.ValidateRegisterMarketItem,
		CancelTradingItemEvent:    Validate.ValidateCancelMarketItem,
		ChangeItemOwnerEvent:      Validate.ValidateChangeItemOwner,
		EnhancingItemRequestEvent: Validate.ValidateEnhancingItemRequest,
		EnhancingItemFailureEvent: Validate.ValidateEnhancingItemFailure,
		EnhancingItemSuccessEvent: Validate.ValidateEnhancingItemSuccess,
		BurningItemRequestEvent:   Validate.ValidateBurningItemRequest,
		BurningItemFailureEvent:   Validate.ValidateBurningItemFailure,
		BurningItemSuccessEvent:   Validate.ValidateBurningItemSuccess,
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

// PartitionKey | string 을 type 으로 재정의 했다. 명시적으로 보이고, 추가 func 을 더 붙일 수 있기 때문임
type PartitionKey string

// RequestBody | Item 도메인의 Request Body, 계속 추가 될 수 있다.
type RequestBody struct {
	*Owner
	*OnchainLink
	*Data
}

func NewRequests(req BodyModel) *RequestBody {
	if req == nil {
		return &RequestBody{}
	}
	rb := &RequestBody{}
	rb.SetReq(req)
	return rb
}
func (i *RequestBody) SetReq(req BodyModel) {
	switch req.(type) {
	case *Owner:
		i.Owner = req.(*Owner)
	case *OnchainLink:
		i.OnchainLink = req.(*OnchainLink)
	case *Data:
		i.Data = req.(*Data)
	}
}
func (i *RequestBody) GetItemOwner() *Owner {
	if i == nil || i.Owner == nil {
		return nil
	}
	return i.Owner
}
func (i *RequestBody) GetItemOnchainLink() *OnchainLink {
	if i == nil || i.OnchainLink == nil {
		return nil
	}
	return i.OnchainLink
}
func (i *RequestBody) GetItemData() *Data {
	if i == nil || i.Data == nil {
		return nil
	}
	return i.Data
}

// BodyModel | Item RequestBody 를 하나의 type 으로 받기 위해 Trick 을 사용, 이 interface 는 아무런 액션도 하지 않음
type BodyModel interface {
	bodyModel()
}

type Owner struct {
	OwnerKey string `json:"ownerKey"`
}

// bodyModel interface 를 구현했기에, RequestBody 에 다같이 들어갈 수 있는 것
func (i *Owner) bodyModel() {
	return
}

type OnchainLink struct {
	TokenId      string `json:"tokenId"`
	MintingNo    string `json:"mintingNo"`
	ContractAddr string `json:"contractAddr"`
}

func (i *OnchainLink) bodyModel() {
	return
}

type Data struct {
	Data string `json:"data"`
}

func (i *Data) bodyModel() {
	return
}
