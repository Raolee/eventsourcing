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
		OwnerKey: "raol",
	})
	req.SetReq(&OnchainLink{
		TokenId:      "token",
		MintingNo:    "00000001",
		ContractAddr: "contractAddr",
	})
	req.SetReq(&Data{
		Data: "data",
	})
	event := NewEvent(CreateItemEvent, "v1", assetKey, req)
	json, err := json.Marshal(event)
	if err != nil {
		t.Error(err)
	}
	t.Log(string(json))
}
