// Event Sourcing Storage Interfaces
//
// 이벤트 소싱에서 필수적으로 구현해야하는 스토리지의 기능을 인터페이스로 정의한다.
//
// [Event Storage 인터페이스]
// Event Storage 에서 Event 는 저장만 가능하고, 수정하거나 삭제할 수 없다는 원칙을 지키도록 구현한다.
// AddEvent 가 Command 영역이 되고, 그외 필요에 따라 Event 를 Get 하는 방법이 더 늘어날 수 있다.
//
// [CommonState Snapshot Storage 인터페이스]
// Snapshot 은 CommonState 의 특정 상태를 의미한다. 즉, Snapshot 은 최신일 수도 있지만 과거의 CommonState 상태일 수도 있다.
// 구현 시 CommonState 가 변경될 때 마다 Snapshot 에 저장하면 Material View 가 된다. 이 경우 Snapshot 을 조회하면 최신 CommonState 를 알 수 있다.
// 반면, Snapshot 을 특정 시점이나, 룰에 따라 만들고 수정하면 과거의 CommonState 를 보게 된다. 혹 최신 상태를 얻고자 한다면,
// Snapshot + Event replay 를 합쳐서 최신 CommonState 를 알아낼 수 있다.
//
// written by Raol
//

package eventsourcing

// EventStorage | Event 저장소의 인터페이스
type EventStorage[R any] interface {
	IncreaseEventNo(pk PartitionKey) (eventNo int, err error)                // atomic 하게 event 번호를 증가시켜 가져온다. pk가 처음 들어오는 것이면 1을 리턴
	AddEvent(event *Event[R]) error                                          // event 를 저장
	GetEvent(id EventId) (*Event[R], error)                                  // event 를 조회
	GetEvents(pk PartitionKey) ([]*Event[R], error)                          // partition key 의 전체 event list 를 조회
	GetEventsAfterEventNo(pk PartitionKey, eventNo int) ([]*Event[R], error) // partition key 의 eventNo 보다 큰 events 를 조회
	GetLastEvent(pk PartitionKey) (*Event[R], error)                         // partition key 의 마지막 event 를 조회
	GetLock(pk PartitionKey) (bool, error)                                   // partition key 에 락이 걸려있는지 확인
	Lock(pk PartitionKey) (already bool, err error)                          // partition key 을 잠금. 이미 잠금이 되어있으면 already=true
	Unlock(pk PartitionKey) (already bool, err error)                        // partition key 를 잠금 해제. 이미 잠금 해제가 되어있으면 already=false
}

// StateSnapshotStorage | Event Snapshot 저장소의 인터페이스
type StateSnapshotStorage[S CommonState[R], R any] interface {
	SaveSnapshot(pk PartitionKey, state *State[S, R]) error      // PartitionKey 의 snapshot 저장
	GetSnapshot(pk PartitionKey) (state *State[S, R], err error) // PartitionKey 로 검색하여 현재 Snapshot 조회
}
