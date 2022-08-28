package eventsourcing

import (
	"eventsourcing/item"
	"time"
)

type ItemCommander struct {
}

func (c *ItemCommander) apply(eventCh <-chan item.Event) error {
	for {
		select {
		case <-eventCh:
			//TODO
		default:
			time.Sleep(100 * time.Millisecond)
		}
	}
}
