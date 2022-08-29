package item

import (
	"github.com/rs/xid"
	"testing"
)

func TestCommand(t *testing.T) {
	createCommand := Command.CreateItem

	state := &State{}
	assetKey := PartitionKey(xid.New().String())
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
	createEvent := NewEvent(CreateItemEvent, "v1", assetKey, req)

	state = createCommand(state, createEvent)
	t.Log(state)

	saveReq := NewRequests(&Data{Data: "saved data"})
	saveEvent := NewEvent(SaveItemDataEvent, "v1", assetKey, saveReq)

	state = Command.SaveItemData(state, saveEvent)
	t.Log(state)
}
