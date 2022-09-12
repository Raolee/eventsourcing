package model

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
