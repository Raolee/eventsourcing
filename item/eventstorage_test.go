package item

import (
	"github.com/rs/xid"
	"testing"
)

var (
	mockStorage  = NewMockEventStorage()
	partitionKey = PartitionKey(xid.New().String())
)

func TestMockEventStorage(t *testing.T) {
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
	event := NewEvent(CreateItemEvent, "v1", partitionKey, req)
	err := mockStorage.SetEvent(event)
	if err != nil {
		t.Error(err)
	}

	// get event
	getEvent, err := mockStorage.GetEvent(event.Id)
	if err != nil {
		t.Error(err)
	}

	t.Log(getEvent)

	// set another event
	saveEvent := NewEvent(SaveItemDataEvent, "v1", partitionKey, NewRequests(&Data{Data: "data"}))
	err = mockStorage.SetEvent(saveEvent)
	if err != nil {
		t.Error(err)
	}
	mintingReqEvent := NewEvent(RequestMintingItemEvent, "v1", partitionKey, nil)
	err = mockStorage.SetEvent(mintingReqEvent)
	if err != nil {
		t.Error(err)
	}

	// get events
	events, err := mockStorage.GetEvents(event.PartitionKey)
	if err != nil {
		t.Error(err)
	}

	t.Log(events)
}
