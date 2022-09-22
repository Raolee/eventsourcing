비동기적으로 Event Sourcing 을 하는 시퀀스 다이어그램
========================================
* * *
### Purpose
1. 빠른 응답이 중요한 도메인
2. State 및 Status 단순하고, 최신 이벤트만으로 Validating 이 가능한 도메인
3. Query view 구성이 되기까지 시간이 상대적으로 여유있는 도메인
### Domain 설명
- Provider
  - Event Sourcing 요청자
- EventService : 하위 도메인을 노출하는 서비스
  - Domains
    - EventMarker : Pk 마다 최신 이벤트를 마킹(marking)하는 서비스/스토리지
    - EventWorker : Event 메세지 스트림을 소비하는 워커
- Storage
  - MessageStream : Event 메세지 스트림
  - Event Storage : Event 저장소
  - Snapshot Storage : Replay 된 State Snapshot 저장소
* * *
### Core Rules
- Validating
  - 마킹된 이벤트만 사용하여 validating 한다
- Command 요청 시
  - Message Stream에 잘 써지면 return
- Snapshot 저장
  - Snapshot Rule 에 따라, 주기적으로 저장
  - Query 요청 중, 최신 State 를 조회할 때 저장

* * *
### Request Event Sequence Diagram
```mermaid
sequenceDiagram
    Provider ->>+ EventService: request event
    EventService ->> EventService:get validator
    alt has validator
        EventService ->>+ EventMarker: get marked eventType (with expire)
        EventMarker ->>- EventService: eventType
        EventService ->> EventService: validate
    end
    EventService ->>+ EventMarker: mark eventType
    EventMarker ->> EventMarker: overwrite marked eventType
    EventMarker ->>- EventService: ok
    EventService ->> EventService: make event
    EventService ->>+ MessageStream: produce event
    EventService ->>- Provider : result
```
* * *
### Event Consuming Sequence Diagram
```mermaid
sequenceDiagram
    MessageStream ->> MessageStream: storing events
    loop
        EventWorker ->>+ EventWorker: start
        EventWorker ->>+ MessageStream: consume event
        MessageStream ->>- EventWorker: events
        loop
            EventWorker ->>+ EventStorage: add event
            EventStorage ->>- EventWorker: ok
            EventWorker ->> EventWorker: check rule
            alt make snapshot
                par using goroutine
                    EventWorker ->>+ SnapshotStorage: get snapshot
                    SnapshotStorage ->>- EventWorker: snapshot
                    EventWorker ->>+ EventStroage: get events after snapshot
                    EventStroage ->>- EventWorker: events
                    EventWorker ->> EventWorker: replay events
                    EventWorker ->>+ SnapshotStorage: update snapshot
                    SnapshotStorage ->>- EventWorker: ok
                end
            end
            EventWorker ->> MessageStream: produce event
        end
        EventWorker ->>- EventWorker: continue
    end
    EventWorker ->> MessageStream: commit events
```
* * *
### Event/Snapshot/State Query Sequence Diagram
```mermaid
sequenceDiagram
    Provider ->>+ EventService: get events
    EventService ->>+ EventStorage: query
    EventStorage ->>- EventService: events
    EventService ->>- Provider: events
    
    Provider ->>+ EventService: get snapshot
    EventService ->>+ SnapshotStorage: query
    SnapshotStorage ->>- EventService: snapshot
    EventService ->>- Provider: snapshot
    
    Provider ->>+ EventService: get latest state 
    EventService ->>+ SnapshotStorage: query
    SnapshotStorage ->>- EventService: snapshot
    EventService ->>+ EventStorage: get events after snapshot
    EventStorage ->>- EventService: events
    EventService ->> EventService: replay
    EventService ->>- Provider: latest state
```