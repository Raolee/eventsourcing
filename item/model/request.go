package model

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
