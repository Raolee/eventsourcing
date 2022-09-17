package currency

import (
	es "eventsourcing"
	"github.com/aws/smithy-go/ptr"
	"time"
)

var (
	Rule *es.Rule
)

func init() {

	Rule = &es.Rule{
		AlwaysSnapshot:  ptr.Bool(false),
		MinSnapshotTerm: ptr.Duration(1 * time.Second),
		MinEventNoTerm:  ptr.Int(5),
	}
}
