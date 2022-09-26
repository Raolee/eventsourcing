package example

import (
	"eventsourcing/example/currency"
	"eventsourcing/example/storage"
	es "eventsourcing/manager"
)

var (
	CurrencyEsManager es.Manager[currency.State, currency.Request]
)

func init() {
	CurrencyEsManager = es.NewBaseManager[currency.State, currency.Request](
		currency.Rule,
		currency.Processor,
		currency.Validator,
		storage.NewCurrencyEventStorage(),
		storage.NewCurrencySnapshotStorage(),
	)
}
