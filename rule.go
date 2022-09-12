package eventsourcing

import (
	"github.com/aws/smithy-go/ptr"
	"time"
)

// Rule | EventSourcing 의 규칙을 정의하는 구조체
type Rule struct {
	// snapshot 저장 규칙
	AlwaysSnapshot  *bool          // default false, 항상 snapshot 을 최신으로 유지하는지 여부
	MinSnapshotTerm *time.Duration // default 1 min, 최근 eventAt 과 snapshot eventAt 의  최소 시간 차이. 이 값을 넘어가면 snapshot 을 저장한다.
	MinEventNoTerm  *int           // default 5, 최근 eventNo 와 snapshot 의 eventNo 와 최소 차이. 이 값을 넘어가면 snapshot 을 저장한다.
}

func (r *Rule) Merge(rule *Rule) {
	if rule != nil {
		if rule.AlwaysSnapshot != nil {
			r.AlwaysSnapshot = rule.AlwaysSnapshot
		}
		if rule.MinSnapshotTerm != nil {
			r.MinSnapshotTerm = rule.MinSnapshotTerm
		}
		if rule.MinEventNoTerm != nil {
			r.MinEventNoTerm = rule.MinEventNoTerm
		}
	}
}

func newDefaultRule() *Rule {
	return &Rule{
		AlwaysSnapshot:  ptr.Bool(false),
		MinSnapshotTerm: ptr.Duration(1 * time.Minute),
		MinEventNoTerm:  ptr.Int(5),
	}
}
