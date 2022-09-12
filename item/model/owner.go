package model

// Owner | State 에서 다루는 소유주 정보
type Owner struct {
	OwnerKey string `json:"ownerKey"`
}

type OwnerModel interface {
	ownerModel()
}

type OwnerType string

const (
	AccountType   = OwnerType("account")
	WalletType    = OwnerType("wallet")
	CharacterType = OwnerType("character")
)

func NewAccountOwner(owner *Owner) *AccountOwner {
	return &AccountOwner{
		Owner: owner,
		Type:  AccountType,
	}
}

type AccountOwner struct {
	*Owner
	Type OwnerType
}

func (o *AccountOwner) ownerModel() {
}

func NewWalletOwner(owner *Owner) *WalletOwner {
	return &WalletOwner{
		Owner: owner,
		Type:  WalletType,
	}
}

type WalletOwner struct {
	*Owner
	Type OwnerType
}

func (o *WalletOwner) ownerModel() {
}

func NewCharacterOwner(owner *Owner) *CharacterOwner {
	return &CharacterOwner{
		Owner: owner,
		Type:  CharacterType,
	}
}

type CharacterOwner struct {
	*Owner
	Type OwnerType
}

func (o *CharacterOwner) ownerModel() {
}
