package item

import (
	"encoding/json"
	"github.com/rs/xid"
	"testing"
)

func TestNewItemEvent(t *testing.T) {
	assetKey := xid.New().String()
	req := NewRequests(nil)
	req.SetReq(&Owner{
		AccountKey: "raol",
	})
	req.SetReq(&ItemOnchainLink{
		MintingNo: "00000001",
		ERC721Contract: &ERC721Contract{
			ContractAddr: "contractAddr",
			TokenId:      "token",
		},
	})
	req.SetReq(&Data{
		Data: "data",
	})
	event := NewEvent(CreateEvent, "v1", PartitionKey(assetKey), req)
	bytes, err := json.Marshal(event)
	if err != nil {
		t.Error(err)
	}
	t.Log(string(bytes))
}
