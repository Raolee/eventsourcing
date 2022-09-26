Storage Requirements
=============================
-----------------------------
### 1. Event Storage

| 중요도 | 요구사항                                 |
|:---:|--------------------------------------|
| 필수  | Partitioning                         |
| 필수  | Sortable Event ID (or Event No)      |
| 필요  | PK의 가장 최근 Event 조회 (Get Latest Once) |
| 선택  | Event 생성일 조회                         |

- Partitioning
  - 분산 저장이 될 수 있어야, 확장성과 고가용성을 확보할 수 있음
  - PK 기준
    - PK 를 공유하는 Event list 를 손쉽게 조회
    - Event ID 중복 가능, 하지만 PK 하위에 존재하므로 이슈 없음
    - PK 를 모르는 상태에서 Event 를 조회할 때 별도 인덱스 필요
  - Event ID 기준
    - Event ID 로 Event 를 쉽게 조회
    - Event ID 중복 방지 보장
    - PK 를 공유하는 Event List 를 조회할 때 별도 인덱스 필요


- Sortable Event ID (or Event No)
  - Snapshot + Event List 로 State 를 만들 때, Event List 를 snapshot 이후로 조회하려면 Event 의 정렬 기준이 필요
  - Event ID 가 Sortable 하게 만들어진다면 Event ID를 그대로 사용하면됨
  - 하지만, Sotable 하지 않은 Event ID 라면 Event No 를 따로 만들어야함
  - Event No 는 PK 안에서만 증가하면 되기에, Globally Atomic 한 처리가 필요하지는 않음
    - ex) pk=1의 이벤트번호 1,2,3,4,5..., pk=2의 이벤트번호 1,2,3,4,5...
  - Event ID가 정렬가능하려면, UUID 와 같은 방법은 사용이 불가능하고, xid 라이브러리를 변경하여 사용해야함


- PK의 가장 최근 Event 조회 (Get Latest Once)
  - 가장 최근에 저장한 Event 로 Event command validate, policy 를 판단하고자 할 때 필요
  - Full scan 은 지양
  - 시간복잡도 O(1) 으로 조회할 수 있는 방안이 Best


- Event 생성일 조회
  - Pk 를 모르는 상태에서, 특정 생성일 기준 이후의 모든 Event 를 알고자 할 때 필요
  - 일괄적으로 특정 생성일 기준 Event 들의 PK 리스트를 알아내고자 할 때 사용
  - 사실, Snapshot 을 다시 만들고자 한다면 Snapshot 자체를 clear 시키고 최신 State 를 조회하는 요청이 있을 때 lazy 하게 만들어도 됨

### 2. Snapshot Storage
| 중요도 | 요구사항                    |
|:---:|-------------------------|
| 필수  | Partitioning            |
| 필수  | State Schema Versioning |
| 선택  | Snapshot 생성일 조회         |
 
- Partitioning
  - Event Storage 와 마찬가지로 확장성과 고가용성을 확보하기 위함


- State Schema Versioning
  - Snapshot 은 State schema 의 버전별로 따로 관리될 수 있음
  - 이 것이 가능하면, State 가 변해도 old/new 를 따로 관리할 수 있음
    - ex) 동일한 이벤트를 재생해서 계좌v1 State와 계좌2 State 를 만들어낼 수 있음

- Snapshot 생성일 조회
  - PK 없이, 특정 생성일 기준 이후의 모든 Snapshot 을 알고자 할 때 필요
  - Snapshot 을 일괄로 삭제하고자 할 때 필요