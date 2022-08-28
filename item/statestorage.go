package item

type ItemStateStorage interface {
	CreateItemState(event *Event) error
	SaveItemState(event *Event) error
	GetItemState(assetKey string) (*State, error)
}
