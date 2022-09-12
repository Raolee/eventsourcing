package item

import (
	"eventsourcing/item/command"
	"github.com/rs/xid"
	"testing"
)

func TestCommand(t *testing.T) {
	createCommand := command.Command.Create

	state := &command.State{}
	assetKey := PartitionKey(xid.New().String())
	req := NewRequests(nil)
	req.SetReq(&Owner{
		AccountKey: "raol",
	})
	req.SetReq(&ItemOnchainLink{
		ERC721Contract: &ERC721Contract{
			ContractAddr: "contractAddr",
			TokenId:      "token",
		},
		MintingNo: "00000001",
	})
	req.SetReq(&Data{
		Data: "data",
	})
	createEvent := NewEvent(CreateEvent, "v1", assetKey, req)

	state = createCommand(state, createEvent)
	t.Log(state)

	saveReq := NewRequests(&Data{Data: "saved data"})
	saveEvent := NewEvent(SaveDataEvent, "v1", assetKey, saveReq)

	state = command.Command.SaveData(state, saveEvent)
	t.Log(state)
}
