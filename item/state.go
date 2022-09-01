package item

import (
	"encoding/json"
	"time"
)

type Command interface {
	CreateItem(*Event) *State
	SaveItemData(*Event) *State
	RemoveItem(*Event) *State

	/** Minting Events **/
	ApplyRequestMintingItem(*Event) *State
	ApplyFailedMintingItem(*Event) *State
	ApplySuccessMintingItem(*Event) *State

	/** Market Events **/
	RegisterMarketItem(*Event) *State
	CancelMarketItem(*Event) *State

	ChangeItemOwner(*Event) *State
	ChangeCharOwner(*Event) *State

	/** Enhance Events **/
	ApplyRequestEnhanceItem(*Event) *State
	ApplyFailedEnhanceItem(*Event) *State
	ApplySuccessEnhanceItem(*Event) *State

	/** Character Events **/
	ApplyRequestBindItemToCharacter(*Event) *State
	ApplyFailedBindItemToCharacter(*Event) *State
	ApplySuccessBindItemToCharacter(*Event) *State
	ApplyRequestUnbindItemToCharacter(*Event) *State
	ApplyFailedUnbindItemToCharacter(*Event) *State
	ApplySuccessUnbindItemToCharacter(*Event) *State

	/** Catalog Events **/
	ApplyRequestBindItemToCatalog(*Event) *State
	ApplyFailedBindItemToCatalog(*Event) *State
	ApplySuccessBindItemToCatalog(*Event) *State

	/** Burn Events **/
	ApplyRequestBurnItem(*Event) *State
	ApplyFailedBurnItem(*Event) *State
	ApplySuccessBurnItem(*Event) *State
}

// Status | Item State 가 가지는 상태 값
type Status string

const (
	Created          = Status("created")            // 처음 상태
	Removed          = Status("removed")            // 삭제한 상태 (Created 일 대만 가능)
	Minting          = Status("minting")            // 민팅중 상태
	Onchain          = Status("onchain")            // 온체인에 올라간 상태, 민팅이 완료되면 Onchain 상태가 된다
	Trading          = Status("trading")            // 거래중 상태, 거래 취소/성공/실패하면 Onchain 상태로 변경된다
	Enhancing        = Status("enhancing")          // 강화중 상태, 강화 성공/실패하면 Onchain 상태로 변경된다
	BindingToChar    = Status("binding_to_char")    // 캐릭터 NFT 에 연결되는 상태, 성공하면 BoundCharNFT 가 되고 실패하면 Onchain 상태로 되돌아 간다
	UnbindingToChar  = Status("unbinding_to_char")  // 캐릭터 NFT 에서 연결을 해제하는 상태, 성공하면 Onchain 이 되고 실패하면 BoundChar 상태로 되돌아 간다
	BoundChar        = Status("bound_char")         // 캐릭터 NFT 에 연결된 상태, 연결을 해제하여 다시 Onchain 상태가 될 수 있다
	BindingToCatalog = Status("binding_to_catalog") // 도감에 등록하는 상태, 성공하면 Cataloged 가 되고 실패하면 Onchain 상태로 되돌아 간다
	BoundCatalog     = Status("bound_catalog")      // 도감에 등록된 상태, 다른 상태로 바꿀 수 없다
	Burning          = Status("burning")            // 소각중 상태, 소각이 성공하면 Burned 상태가 되고 실패하면 Onchain 상태로 되돌아 간다
	Burned           = Status("burned")             // 소각된 상태, 다른 상태로 바꿀 수 없다
)

var (
	_ Partition = &State{}
	_ Command   = &State{}
)

// Partition | State 는 Partition interface 를 구현한다. 즉, State 는 Partition key 를 꼭 가지고 있어야 한다는 의미
type Partition interface {
	PartitionKey() PartitionKey
}

// State | Item Domain 의 State, Command interface 를 구현하여 replay 할 수 있게 만든다.
type State struct {
	AssetKey           PartitionKey        `json:"assetKey"`
	Status             Status              `json:"status"`
	Owner              Owner               `json:"owner"`
	ItemOnchainLink    *ItemOnchainLink    `json:"ItemOnchainLink,omitempty"`
	CharOnchainLink    *CharOnchainLink    `json:"CharOnchainLink,omitempty"`
	CatalogOnchainLink *CatalogOnchainLink `json:"CatalogOnchainLink,omitempty"`
	Data               Data                `json:"data"`
	Lock               bool                `json:"lock"`
	CreatedAt          time.Time           `json:"createdAt"`
	LastEventAt        time.Time           `json:"lastEventAt"`
	LastEventId        EventId             `json:"lastEventId"`   // 현 State 가 만들어진 마지막 EventId
	LastEventName      EventName           `json:"lastEventName"` // 현 STate 가 만들어진 마지막 EventName
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
		AssetKey:           event.PartitionKey,
		Status:             Created,
		Owner:              *event.Requests.GetOwner(),
		ItemOnchainLink:    event.Requests.GetItemContract(),
		CharOnchainLink:    event.Requests.GetCharContract(),
		CatalogOnchainLink: event.Requests.GetCatalogContract(),
		Data:               *event.Requests.GetData(),
		Lock:               false,
		CreatedAt:          time.Now(),
		LastEventAt:        event.EventAt,
		LastEventId:        event.Id,
		LastEventName:      event.Name,
	}
}

func (i *State) setEventInfo(event *Event) {
	i.LastEventAt = event.EventAt
	i.LastEventId = event.Id
	i.LastEventName = event.Name
}

func (i *State) SaveItemData(event *Event) *State {
	i.Data = *event.Requests.GetData()
	i.setEventInfo(event)
	return i
}

func (i *State) RemoveItem(event *Event) *State {
	i.Status = Removed
	i.setEventInfo(event)
	return i
}

func (i *State) ApplyRequestMintingItem(event *Event) *State {
	i.Status = Minting
	i.Lock = true
	i.setEventInfo(event)
	return i
}

func (i *State) ApplyFailedMintingItem(event *Event) *State {
	i.Status = Created
	i.Lock = false
	i.setEventInfo(event)
	return i
}

func (i *State) ApplySuccessMintingItem(event *Event) *State {
	i.ItemOnchainLink = event.Requests.GetItemContract()
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
	i.Owner = *event.Requests.GetOwner()
	i.Status = Onchain
	i.Lock = false
	i.setEventInfo(event)
	return i
}

func (i *State) ChangeCharOwner(event *Event) *State {
	i.Owner = *event.Requests.GetOwner()
	i.Status = Onchain
	i.Lock = false
	i.setEventInfo(event)
	return i
}

func (i *State) ApplyRequestEnhanceItem(event *Event) *State {
	i.Status = Enhancing
	i.Lock = true
	i.setEventInfo(event)
	return i
}

func (i *State) ApplyFailedEnhanceItem(event *Event) *State {
	i.Status = Onchain
	i.Lock = false
	i.setEventInfo(event)
	return i
}

func (i *State) ApplySuccessEnhanceItem(event *Event) *State {
	i.Data = *event.Requests.GetData()
	i.Status = Onchain
	i.Lock = true
	i.setEventInfo(event)
	return i
}

func (i *State) ApplyRequestBindItemToCharacter(event *Event) *State {
	i.Status = BindingToChar
	i.Lock = true
	i.setEventInfo(event)
	return i
}

func (i *State) ApplyFailedBindItemToCharacter(event *Event) *State {
	i.Status = Onchain
	i.Lock = false
	i.setEventInfo(event)
	return i
}

func (i *State) ApplySuccessBindItemToCharacter(event *Event) *State {
	i.CharOnchainLink = event.Requests.GetCharContract()
	i.Status = BoundChar
	i.Lock = false
	i.setEventInfo(event)
	return i
}

func (i *State) ApplyRequestUnbindItemToCharacter(event *Event) *State {
	i.Status = UnbindingToChar
	i.Lock = true
	i.setEventInfo(event)
	return i
}

func (i *State) ApplyFailedUnbindItemToCharacter(event *Event) *State {
	i.Status = BoundChar
	i.Lock = false
	i.setEventInfo(event)
	return i
}

func (i *State) ApplySuccessUnbindItemToCharacter(event *Event) *State {
	i.CharOnchainLink = nil
	i.Status = Onchain
	i.Lock = false
	i.setEventInfo(event)
	return i
}

func (i *State) ApplyRequestBindItemToCatalog(event *Event) *State {
	i.Status = BindingToCatalog
	i.Lock = true
	i.setEventInfo(event)
	return i
}

func (i *State) ApplyFailedBindItemToCatalog(event *Event) *State {
	i.Status = Onchain
	i.Lock = true
	i.setEventInfo(event)
	return i
}

func (i *State) ApplySuccessBindItemToCatalog(event *Event) *State {
	i.CatalogOnchainLink = event.Requests.GetCatalogContract()
	i.Status = BoundCatalog
	i.Lock = true
	i.setEventInfo(event)
	return i
}

func (i *State) ApplyRequestBurnItem(event *Event) *State {
	i.Status = Burning
	i.Lock = true
	i.setEventInfo(event)
	return i
}

func (i *State) ApplyFailedBurnItem(event *Event) *State {
	i.Status = Onchain
	i.Lock = false
	i.setEventInfo(event)
	return i
}

func (i *State) ApplySuccessBurnItem(event *Event) *State {
	i.Status = Burned
	i.Owner.AccountKey = ""
	i.Lock = true
	i.setEventInfo(event)
	return i
}
