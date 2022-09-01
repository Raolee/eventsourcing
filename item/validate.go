package item

type Validate interface {
	ValidateCreateItem(*Event) (bool, error)
	ValidateSaveItemData(*Event) (bool, error)
	ValidateRemoveItem(*Event) (bool, error)
	ValidateMintingItemRequest(*Event) (bool, error)
	ValidateMintingItemFailure(*Event) (bool, error)
	ValidateMintingItemSuccess(*Event) (bool, error)
	ValidateRegisterMarketItem(*Event) (bool, error)
	ValidateCancelMarketItem(*Event) (bool, error)
	ValidateChangeItemOwner(*Event) (bool, error)
	ValidateEnhanceItemRequest(*Event) (bool, error)
	ValidateEnhanceItemFailure(*Event) (bool, error)
	ValidateEnhanceItemSuccess(*Event) (bool, error)
	ValidateBurnItemRequest(*Event) (bool, error)
	ValidateBurnItemFailure(*Event) (bool, error)
	ValidateBurnItemSuccess(*Event) (bool, error)
}

type MockValidator struct {
	EventStorage
	StateSnapshotStorage
}

func NewMockValidator(eventStorage EventStorage, stateSnapshotStorage StateSnapshotStorage) Validate {
	return &MockValidator{
		EventStorage:         eventStorage,
		StateSnapshotStorage: stateSnapshotStorage,
	}
}

func (m *MockValidator) ValidateCreateItem(event *Event) (bool, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MockValidator) ValidateSaveItemData(event *Event) (bool, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MockValidator) ValidateRemoveItem(event *Event) (bool, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MockValidator) ValidateMintingItemRequest(event *Event) (bool, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MockValidator) ValidateMintingItemFailure(event *Event) (bool, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MockValidator) ValidateMintingItemSuccess(event *Event) (bool, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MockValidator) ValidateRegisterMarketItem(event *Event) (bool, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MockValidator) ValidateCancelMarketItem(event *Event) (bool, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MockValidator) ValidateChangeItemOwner(event *Event) (bool, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MockValidator) ValidateEnhanceItemRequest(event *Event) (bool, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MockValidator) ValidateEnhanceItemFailure(event *Event) (bool, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MockValidator) ValidateEnhanceItemSuccess(event *Event) (bool, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MockValidator) ValidateBurnItemRequest(event *Event) (bool, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MockValidator) ValidateBurnItemFailure(event *Event) (bool, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MockValidator) ValidateBurnItemSuccess(event *Event) (bool, error) {
	//TODO implement me
	panic("implement me")
}
