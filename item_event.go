package eventsourcing

import (
	"fmt"
	"github.com/rs/xid"
	"time"
)

type ItemEventName string

const (
	CreateItemEvent           = ItemEventName("createItem")
	SaveItemDataEvent         = ItemEventName("saveItemData")
	RemoveItemEvent           = ItemEventName("removeItem")
	MintingItemRequestEvent   = ItemEventName("mintingItemRequest")
	MintingItemFailureEvent   = ItemEventName("mintingItemFailure")
	MintingItemSuccessEvent   = ItemEventName("mintingItemSuccess")
	RegisterMarketItemEvent   = ItemEventName("registerMarketItem")
	CancelTradingItemEvent    = ItemEventName("cancelTradingItem")
	ChangeItemOwnerEvent      = ItemEventName("successTradingItem")
	EnhancingItemRequestEvent = ItemEventName("enhancingItemRequest")
	EnhancingItemFailureEvent = ItemEventName("enhancingItemFailure")
	EnhancingItemSuccessEvent = ItemEventName("enhancingItemSuccess")
	BurningItemRequestEvent   = ItemEventName("burningItemRequest")
	BurningItemFailureEvent   = ItemEventName("burningItemFailure")
	BurningItemSuccessEvent   = ItemEventName("burningItemSuccess")
)

var (
	ItemEvents = map[ItemEventName]func(command ItemCommand, event *ItemEvent) *ItemState{
		CreateItemEvent:           ItemCommand.CreateItem,
		SaveItemDataEvent:         ItemCommand.SaveItemData,
		RemoveItemEvent:           ItemCommand.RemoveItem,
		MintingItemRequestEvent:   ItemCommand.MintingItemRequest,
		MintingItemFailureEvent:   ItemCommand.MintingItemFailure,
		MintingItemSuccessEvent:   ItemCommand.MintingItemSuccess,
		RegisterMarketItemEvent:   ItemCommand.RegisterMarketItem,
		CancelTradingItemEvent:    ItemCommand.CancelMarketItem,
		ChangeItemOwnerEvent:      ItemCommand.ChangeItemOwner,
		EnhancingItemRequestEvent: ItemCommand.EnhancingItemRequest,
		EnhancingItemFailureEvent: ItemCommand.EnhancingItemFailure,
		EnhancingItemSuccessEvent: ItemCommand.EnhancingItemSuccess,
		BurningItemRequestEvent:   ItemCommand.BurningItemRequest,
		BurningItemFailureEvent:   ItemCommand.BurningItemFailure,
		BurningItemSuccessEvent:   ItemCommand.BurningItemSuccess,
	}
)

type ItemEvent struct {
	Id           string        `json:"id"`
	Name         ItemEventName `json:"name"`
	Version      string        `json:"version"`
	EventAt      time.Time     `json:"eventAt"`
	PartitionKey PartitionKey  `json:"partitionKey"` // partition key
	Requests     *ItemRequests `json:"requests"`
}

func NewItemEvent(name ItemEventName, version, assetKey string, reqs *ItemRequests) *ItemEvent {
	return &ItemEvent{
		Id:           xid.New().String(),
		Name:         name,
		Version:      version,
		EventAt:      time.Now(),
		PartitionKey: PartitionKey(assetKey),
		Requests:     reqs,
	}
}

func (e *ItemEvent) EventName() string {
	return fmt.Sprintf("%s_%s", e.Name, e.Version)
}

type PartitionKey string
type ItemRequests struct {
	*ItemOwner
	*ItemOnchainLink
	*ItemData
}

func NewItemRequests(req ItemRequestTypes) *ItemRequests {
	return &ItemRequests{}
}
func (i *ItemRequests) SetReq(req ItemRequestTypes) {
	switch req.(type) {
	case *ItemOwner:
		i.ItemOwner = req.(*ItemOwner)
	case *ItemOnchainLink:
		i.ItemOnchainLink = req.(*ItemOnchainLink)
	case *ItemData:
		i.ItemData = req.(*ItemData)
	}
}
func (i *ItemRequests) GetItemOwner() *ItemOwner {
	if i == nil || i.ItemOwner == nil {
		return nil
	}
	return i.ItemOwner
}
func (i *ItemRequests) GetItemOnchainLink() *ItemOnchainLink {
	if i == nil || i.ItemOnchainLink == nil {
		return nil
	}
	return i.ItemOnchainLink
}
func (i *ItemRequests) GetItemData() *ItemData {
	if i == nil || i.ItemData == nil {
		return nil
	}
	return i.ItemData
}

type ItemRequestTypes interface {
	itemReqType()
}
type ItemOwner struct {
	OwnerKey string `json:"ownerKey"`
}

func (i *ItemOwner) itemReqType() {
	return
}

type ItemOnchainLink struct {
	TokenId      string `json:"tokenId"`
	MintingNo    string `json:"mintingNo"`
	ContractAddr string `json:"contractAddr"`
}

func (i *ItemOnchainLink) itemReqType() {
	return
}

type ItemData struct {
	Data string `json:"data"`
}

func (i *ItemData) itemReqType() {
	return
}
