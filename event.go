package eventsourcing

import (
	"fmt"
	"github.com/rs/xid"
	"time"
)

type Domain int

const (
	ItemDomain = iota
)

type Event struct {
	Id      string    `json:"id"`
	Domain  Domain    `json:"domain"`
	Name    string    `json:"name"`
	Version string    `json:"version"`
	EventAt time.Time `json:"eventAt"`
}

func (e *Event) EventName() string {
	return fmt.Sprintf("%s_%s", e.Name, e.Version)
}

func NewEvent(domain Domain, name, version string) *Event {
	return &Event{
		Id:      xid.New().String(),
		Domain:  domain,
		Name:    name,
		Version: version,
		EventAt: time.Now(),
	}
}
