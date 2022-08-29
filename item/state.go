package item

import (
	"encoding/json"
	"time"
)

type Command interface {
	CreateItem(*Event) *State
	SaveItemData(*Event) *State
	RemoveItem(*Event) *State
	MintingItemRequest(*Event) *State
	MintingItemFailure(*Event) *State
	MintingItemSuccess(*Event) *State
	RegisterMarketItem(*Event) *State
	CancelMarketItem(*Event) *State
	ChangeItemOwner(*Event) *State
	EnhancingItemRequest(*Event) *State
	EnhancingItemFailure(*Event) *State
	EnhancingItemSuccess(*Event) *State
	BurningItemRequest(*Event) *State
	BurningItemFailure(*Event) *State
	BurningItemSuccess(*Event) *State
}

// Status | Item State 가 가지는 상태 값
type Status string

const (
	Created   = Status("created")
	Removed   = Status("removed")
	Minting   = Status("minting")
	Onchain   = Status("onchain")
	Trading   = Status("trading")
	Enhancing = Status("enhancing")
	Burning   = Status("burning")
	Burned    = Status("burned")
)

// Partition | State 는 Partition interface 를 구현한다. 즉, State 는 Partition key 를 꼭 가지고 있어야 한다는 의미
type Partition interface {
	PartitionKey() PartitionKey
}

// State | Item Domain 의 State, Command interface 를 구현하여 replay 할 수 있게 만든다.
type State struct {
	AssetKey      PartitionKey `json:"assetKey"`
	Status        Status       `json:"status"`
	Owner         Owner        `json:"owner"`
	OnchainLink   OnchainLink  `json:"onchainLink"`
	Data          Data         `json:"data"`
	Lock          bool         `json:"lock"`
	CreatedAt     time.Time    `json:"createdAt"`
	LastEventAt   time.Time    `json:"lastEventAt"`
	LastEventId   EventId      `json:"lastEventId"`   // 현 State 가 만들어진 마지막 EventId
	LastEventName EventName    `json:"lastEventName"` // 현 STate 가 만들어진 마지막 EventName
}

func (i *State) PartitionKey() PartitionKey {
	return i.AssetKey
}

func (i *State) String() string {
	if i == nil {
		return "{}"
	}
	marshal, err := json.MarshalIndent(i, "", "    ")
	if err != nil {
		return "{}"
	}
	return string(marshal)
}

func (i *State) CreateItem(event *Event) *State {
	return &State{
		AssetKey:    event.PartitionKey,
		Status:      Created,
		Owner:       *event.Requests.GetItemOwner(),
		OnchainLink: *event.Requests.GetItemOnchainLink(),
		Data:        *event.Requests.GetItemData(),
		Lock:        false,
		CreatedAt:   time.Now(),
		LastEventAt: event.EventAt,
		LastEventId: event.Id,
	}
}

func (i *State) setEventInfo(event *Event) {
	i.LastEventAt = event.EventAt
	i.LastEventId = event.Id
	i.LastEventName = event.Name
}

func (i *State) SaveItemData(event *Event) *State {
	i.Data = *event.Requests.GetItemData()
	i.setEventInfo(event)
	return i
}

func (i *State) RemoveItem(event *Event) *State {
	i.Status = Removed
	i.setEventInfo(event)
	return i
}

func (i *State) MintingItemRequest(event *Event) *State {
	i.Status = Minting
	i.Lock = true
	i.setEventInfo(event)
	return i
}

func (i *State) MintingItemFailure(event *Event) *State {
	i.Status = Created
	i.Lock = false
	i.setEventInfo(event)
	return i
}

func (i *State) MintingItemSuccess(event *Event) *State {
	i.OnchainLink = *event.Requests.GetItemOnchainLink()
	i.Status = Onchain
	i.Lock = false
	i.setEventInfo(event)
	return i
}

func (i *State) RegisterMarketItem(event *Event) *State {
	i.Status = Trading
	i.Lock = true
	i.setEventInfo(event)
	return i
}

func (i *State) CancelMarketItem(event *Event) *State {
	i.Status = Onchain
	i.Lock = false
	i.setEventInfo(event)
	return i
}

func (i *State) ChangeItemOwner(event *Event) *State {
	i.Owner = *event.Requests.GetItemOwner()
	i.Status = Onchain
	i.Lock = false
	i.setEventInfo(event)
	return i
}

func (i *State) EnhancingItemRequest(event *Event) *State {
	i.Status = Enhancing
	i.Lock = true
	i.setEventInfo(event)
	return i
}

func (i *State) EnhancingItemFailure(event *Event) *State {
	i.Status = Onchain
	i.Lock = false
	i.setEventInfo(event)
	return i
}

func (i *State) EnhancingItemSuccess(event *Event) *State {
	i.Data = *event.Requests.GetItemData()
	i.Status = Onchain
	i.Lock = true
	i.setEventInfo(event)
	return i
}

func (i *State) BurningItemRequest(event *Event) *State {
	i.Status = Burning
	i.Lock = true
	i.setEventInfo(event)
	return i
}

func (i *State) BurningItemFailure(event *Event) *State {
	i.Status = Onchain
	i.Lock = false
	i.setEventInfo(event)
	return i
}

func (i *State) BurningItemSuccess(event *Event) *State {
	i.Status = Burned
	i.Lock = true
	i.setEventInfo(event)
	return i
}
