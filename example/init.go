package example

import (
	es "eventsourcing"
	"eventsourcing/example/currency"
	"eventsourcing/example/storage"
)

var (
	CurrencyEsManager es.Manager[currency.State, currency.Request]
)

func init() {
	CurrencyEsManager = es.NewBaseManager[currency.State, currency.Request](
		currency.Rule,
		currency.Commander,
		currency.Validator,
		storage.NewCurrencyEventStorage(),
		storage.NewCurrencySnapshotStorage(),
	)
}
