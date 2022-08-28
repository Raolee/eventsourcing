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
	Events = map[EventName]func(command Command, event *Event) *State{
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
)

type Event struct {
	Id           string       `json:"id"`
	Name         EventName    `json:"name"`
	Version      string       `json:"version"`
	EventAt      time.Time    `json:"eventAt"`
	PartitionKey PartitionKey `json:"partitionKey"` // partition key
	Requests     *RequestBody `json:"requests"`
}

func NewEvent(name EventName, version, assetKey string, reqs *RequestBody) *Event {
	return &Event{
		Id:           xid.New().String(),
		Name:         name,
		Version:      version,
		EventAt:      time.Now(),
		PartitionKey: PartitionKey(assetKey),
		Requests:     reqs,
	}
}

func (e *Event) String() string {
	if e == nil {
		return "{}"
	}
	json, err := json.MarshalIndent(e, "", "    ")
	if err != nil {
		return "{}"
	}
	return string(json)
}

func (e *Event) EventName() string {
	return fmt.Sprintf("%s_%s", e.Name, e.Version)
}

// EventList event list 정의
type EventList []Event

func (e EventList) String() string {
	if e == nil {
		return "[]"
	}
	json, err := json.MarshalIndent(e, "", "    ")
	if err != nil {
		return "[]"
	}
	return string(json)
}

func (e EventList) Iterate() <-chan *Event {
	ch := make(chan *Event)
	go func() {
		for _, event := range ([]Event)(e) {
			ch <- &event
		}
		close(ch)
	}()
	return ch
}

type PartitionKey string
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

type BodyModel interface {
	bodyModel()
}
type Owner struct {
	OwnerKey string `json:"ownerKey"`
}

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
