package eventsourcing

type ItemStateStorage interface {
	CreateItemState(event *ItemEvent) error
	SaveItemState(event *ItemEvent) error
	GetItemState(assetKey string) (*ItemState, error)
}
