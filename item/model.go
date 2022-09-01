package item

/*
model.go 에서 정의하는 것들은 Event 와 State 모두 활용된다.
*/

// PartitionKey | string 을 type 으로 재정의 했다. 명시적으로 보이고, 추가 func 을 더 붙일 수 있기 때문임
type PartitionKey string

// BodyModel | Item RequestBody 를 하나의 type 으로 받기 위해 Trick 을 사용, 이 interface 는 아무런 액션도 하지 않음
type BodyModel interface {
	bodyModel()
}

// bodyModel interface 를 구현했기에, RequestBody 에 다같이 들어갈 수 있는 것
func (i *Owner) bodyModel() {
	return
}

// RequestBody | Item 도메인의 Request Body, 계속 추가 될 수 있다.
type RequestBody struct {
	*Owner              `json:"owner,omitempty"`
	*ItemOnchainLink    `json:"itemOnchainLink,omitempty"`
	*CharOnchainLink    `json:"charOnchainLink,omitempty"`
	*CatalogOnchainLink `json:"catalogOnchainLink,omitempty"`
	*Data               `json:"data,omitempty"`
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
	case *ItemOnchainLink:
		i.ItemOnchainLink = req.(*ItemOnchainLink)
	case *CharOnchainLink:
		i.CharOnchainLink = req.(*CharOnchainLink)
	case *CatalogOnchainLink:
		i.CatalogOnchainLink = req.(*CatalogOnchainLink)
	case *Data:
		i.Data = req.(*Data)
	}
}
func (i *RequestBody) GetOwner() *Owner {
	if i == nil || i.Owner == nil {
		return nil
	}
	return i.Owner
}
func (i *RequestBody) GetItemContract() *ItemOnchainLink {
	if i == nil || i.ItemOnchainLink == nil {
		return nil
	}
	return i.ItemOnchainLink
}
func (i *RequestBody) GetCharContract() *CharOnchainLink {
	if i == nil || i.CharOnchainLink == nil {
		return nil
	}
	return i.CharOnchainLink
}
func (i *RequestBody) GetCatalogContract() *CatalogOnchainLink {
	if i == nil || i.CatalogOnchainLink == nil {
		return nil
	}
	return i.CatalogOnchainLink
}
func (i *RequestBody) GetData() *Data {
	if i == nil || i.Data == nil {
		return nil
	}
	return i.Data
}

// Owner | State 에서 다루는 소유주 정보
type Owner struct {
	AccountKey string  `json:"accountKey"`
	WalletAddr *string `json:"walletAddr"`
}

// ERC721Contract | EVM 의 ERC721 컨트랙트를 의미
type ERC721Contract struct {
	ContractAddr string `json:"contractAddr"`
	TokenId      string `json:"tokenId"`
	TxHash       string `json:"txHash"`
}

func (i *ERC721Contract) bodyModel() {
	return
}

// ItemOnchainLink | ERC 721 로 구현된 아이템 컨트랙트와의 링크
type ItemOnchainLink struct {
	*ERC721Contract
	MintingNo string `json:"mintingNo"`
}

// CharOnchainLink | ERC 721 로 구현된 캐릭터 컨트랙트와의 링크
type CharOnchainLink struct {
	*ERC721Contract
}

// CatalogOnchainLink | ERC 721 로 구현된 카탈로그 컨트랙트와의 링크
type CatalogOnchainLink struct {
	*ERC721Contract
}

// Data | State 에 저장된 data
type Data struct {
	Data string `json:"data"`
}

func (i *Data) bodyModel() {
	return
}
