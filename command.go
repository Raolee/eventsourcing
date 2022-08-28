package eventsourcing

import (
	"time"
)

type ItemCommander struct {
}

func (c *ItemCommander) apply(eventCh <-chan ItemEvent) error {
	for {
		select {
		case <-eventCh:
			//TODO
		default:
			time.Sleep(100 * time.Millisecond)
		}
	}
}
