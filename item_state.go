package eventsourcing

import (
	"time"
)

type ItemCommand interface {
	CreateItem(*ItemEvent) *ItemState
	SaveItemData(*ItemEvent) *ItemState
	RemoveItem(*ItemEvent) *ItemState
	MintingItemRequest(*ItemEvent) *ItemState
	MintingItemFailure(*ItemEvent) *ItemState
	MintingItemSuccess(*ItemEvent) *ItemState
	RegisterMarketItem(*ItemEvent) *ItemState
	CancelMarketItem(*ItemEvent) *ItemState
	ChangeItemOwner(*ItemEvent) *ItemState
	EnhancingItemRequest(*ItemEvent) *ItemState
	EnhancingItemFailure(*ItemEvent) *ItemState
	EnhancingItemSuccess(*ItemEvent) *ItemState
	BurningItemRequest(*ItemEvent) *ItemState
	BurningItemFailure(*ItemEvent) *ItemState
	BurningItemSuccess(*ItemEvent) *ItemState
}

type ItemStatus string

const (
	Created   = ItemStatus("created")
	Removed   = ItemStatus("removed")
	Minting   = ItemStatus("minting")
	Onchain   = ItemStatus("onchain")
	Trading   = ItemStatus("trading")
	Enhancing = ItemStatus("enhancing")
	Burning   = ItemStatus("burning")
	Burned    = ItemStatus("burned")
)

type ItemState struct {
	AssetKey    string
	Status      ItemStatus
	Owner       ItemOwner
	OnchainLink ItemOnchainLink
	Data        ItemData
	Lock        bool
	CreatedAt   time.Time
	LastEventAt time.Time
	LastEventId string
}

func (i *ItemState) CreateItem(event *ItemEvent) *ItemState {
	return &ItemState{
		AssetKey:    string(event.PartitionKey),
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

func (i *ItemState) SaveItemData(event *ItemEvent) *ItemState {
	i.Data = *event.Requests.GetItemData()
	i.LastEventAt = event.EventAt
	i.LastEventId = event.Id
	return i
}

func (i *ItemState) RemoveItem(event *ItemEvent) *ItemState {
	i.Status = Removed
	i.LastEventAt = event.EventAt
	i.LastEventId = event.Id
	return i
}

func (i *ItemState) MintingItemRequest(event *ItemEvent) *ItemState {
	i.Status = Minting
	i.Lock = true
	i.LastEventAt = event.EventAt
	i.LastEventId = event.Id
	return i
}

func (i *ItemState) MintingItemFailure(event *ItemEvent) *ItemState {
	i.Status = Created
	i.Lock = false
	i.LastEventAt = event.EventAt
	i.LastEventId = event.Id
	return i
}

func (i *ItemState) MintingItemSuccess(event *ItemEvent) *ItemState {
	i.OnchainLink = *event.Requests.GetItemOnchainLink()
	i.Status = Onchain
	i.Lock = false
	i.LastEventAt = event.EventAt
	i.LastEventId = event.Id
	return i
}

func (i *ItemState) RegisterMarketItem(event *ItemEvent) *ItemState {
	i.Status = Trading
	i.Lock = true
	i.LastEventAt = event.EventAt
	i.LastEventId = event.Id
	return i
}

func (i *ItemState) CancelMarketItem(event *ItemEvent) *ItemState {
	i.Status = Onchain
	i.Lock = false
	i.LastEventAt = event.EventAt
	i.LastEventId = event.Id
	return i
}

func (i *ItemState) ChangeItemOwner(event *ItemEvent) *ItemState {
	i.Owner = *event.Requests.GetItemOwner()
	i.Status = Onchain
	i.Lock = false
	i.LastEventAt = event.EventAt
	i.LastEventId = event.Id
	return i
}

func (i *ItemState) EnhancingItemRequest(event *ItemEvent) *ItemState {
	i.Status = Enhancing
	i.Lock = true
	i.LastEventAt = event.EventAt
	i.LastEventId = event.Id
	return i
}

func (i *ItemState) EnhancingItemFailure(event *ItemEvent) *ItemState {
	i.Status = Onchain
	i.Lock = false
	i.LastEventAt = event.EventAt
	i.LastEventId = event.Id
	return i
}

func (i *ItemState) EnhancingItemSuccess(event *ItemEvent) *ItemState {
	i.Data = *event.Requests.GetItemData()
	i.Status = Onchain
	i.Lock = true
	i.LastEventAt = event.EventAt
	i.LastEventId = event.Id
	return i
}

func (i *ItemState) BurningItemRequest(event *ItemEvent) *ItemState {
	i.Status = Burning
	i.Lock = true
	i.LastEventAt = event.EventAt
	i.LastEventId = event.Id
	return i
}

func (i *ItemState) BurningItemFailure(event *ItemEvent) *ItemState {
	i.Status = Onchain
	i.Lock = false
	i.LastEventAt = event.EventAt
	i.LastEventId = event.Id
	return i
}

func (i *ItemState) BurningItemSuccess(event *ItemEvent) *ItemState {
	i.Status = Burned
	i.Lock = true
	i.LastEventAt = event.EventAt
	i.LastEventId = event.Id
	return i
}
