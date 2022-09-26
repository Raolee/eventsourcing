Storage Comparison
=============================
-----------------------------
## Event Storage
### 1. Dynamo DB
    
| 평가  | 요구사항                            | 방법                                                                              |
|:---:|---------------------------------|---------------------------------------------------------------------------------|
| 만족  | Partitioning                    | PK 또는 Event ID 파티셔닝 가능, GSI 로 보완                                                |
| 만족  | Sortable Event ID (or Event No) | Event No 는 따로 Event No 만 관리하는 Table 필요. 반면 Event ID는 assetid lib를 개선하여 이용 가능할 듯 |
| 만족  | PK의 가장 최근 Event 조회              | 가장 최근 event 만 저장하는 table을 따로 두어 여기서 조회하여 처리 가능                                  |
| 만족  | Event 생성일 조회                    | 생성일을 GSI 설정하면 가능                                                                |

- 예상되는 Tables
  - event_history_table
    - pk에 쌓인 이벤트를 저장
  - latest_event_table
    - pk의 최근 이벤트만 저장
  - (pk_event_no_table)
    - event no 만 저장


- 장/단점
  - 장점
    - 사실상 무한한 읽기/쓰기 성능
    - 설계된 table 이 단단하다면, LSI, GSI 추가로 조회 요구사항 만족 가능
    - Dynamodb 의 제약조건인 끊어 읽기가 메모리에 로드되는 양을 강제하여 안정성을 확보할 수 있음
      - 예를들어 event가 100000개인 것을 replay 해야하는데, 이를 제한된 메모리에 올리면 OOM이 발생할 수 있음
  - 단점
    - 무한한 읽기/쓰기 성능에 맞춰 비용이 클 것으로 예상 (실제로...얼만큼의 비용인지는 감이 안옴)
    - Get Latest Once 와 같은 조회 방식을 제공하지 않음
    - GSI 에 성능 한계가 있다고 알고 있음.
    - 혹, Event 의 크기가 400kb 를 넘어가면 그때부터 사용 불가능


### 2. S3

| 평가  | 요구사항                         | 방법                                                                        |
|:---:|------------------------------|---------------------------------------------------------------------------|
| 만족  | Partitioning                 | prefix를 PK 로 설정 혹은 Event ID 를 저장하고, tag로 PK 설정                            |
| 만족  | Event No or Ordered Event ID | Event No 는 PK Event No 버킷을 하나 만들어서 버전 개수로 처리, Event ID는 기존 asset_id 이용 가능 |
| 만족  | PK의 가장 최근 Event 조회           | PK Latest Event 버킷을 이용하여, 항상 최신만 update 하도록하여 O(1) 가능                     |
| 만족  | Event 생성일 조회                 | S3의 버킷 오브젝트의 최근 수정일로 조회 가능                                                |

- Buckets 
    - event_history_bucket
      - {pk}/{event_id}.json 으로 저장
    - latest_event_bucket
      - {pk}.json 으로 저장
    - (pk_event_no_bucket)
      - bucket object 에 tag 를 달고 select 조회(?) 


- 장/단점
    - 장점
        - PK 가 다르다면 거의 무한한 읽기/쓰기 성능
        - 고가용성 확보
        - 관리하기로한 prefix 나 object naming 이 단단하다면 스키마가 변해도 큰 이슈가 없음
        - 담는 Object 의 크기가 커도 수용 가능 (400kb 이상도 가능)
        - 코드와 버킷오브젝트의 스키마가 일치
    - 단점
        - Pk 하나에 초당 3000쓰기/5000읽기의 성능 제약사항 (사실 단점이 아닐수도?) 
        - 예상치 못한 조회 기능을 구현하기 어려울 수 있음
          - 구현을 하기 위한 bucket 파편화가 불가피
        - Get Latest Once 와 같은 조회 방식을 제공하지 않음

----------------------------------
## Snapshot Storage
### 1. Dynamodb

| 평가  | 요구사항                    | 방법                             |
|:---:|-------------------------|--------------------------------|
| 만족  | Partitioning            | PK 로 가능                        |
| 만족  | State Schema Versioning | Schema Version 을 Sort key 로 활용 |
| 만족  | Snapshot 생성일 조회         | 생성일로 GSI 를 설정하여 조회             |

- 예상 Tables
    - snapshot_table
      - DDB PartitionKey : State PK (string)
      - DDB SortKey : Schema version (숫자 혹은 스키마명)


- 장/단점
    - 장점
        - Event Storage의 DDB 장점과 동일
        - State 자체를 DDB 로 표현하면, 일부 필드만 update 하여 비용 절감 가능
    - 단점
        - Event Storage 의 DDB 단점과 동일
        

### 2. S3

| 평가  | 요구사항                    | 방법                                        |
|:---:|-------------------------|-------------------------------------------|
| 만족  | Partitioning            | prefix 를 PK 로 설정                          |
| 만족  | State Schema Versioning | prefix 하위에 Schema Version 을 Object 명으로 사용 |
| 만족  | Snapshot 생성일 조회         | 버킷 전체에서 생성일로 조회                           |

- 에상 Buckets
    - snapshot_bucket
      - {pk}/{schema_ver}.json 으로 저장


- 장/단점
    - 장점
        - Event Storage의 S3 장점과 동일
        - 쉽게 추출 가능, 외부 노출 가능
    - 단점
        - Event Storage 의 S3 단점과 동일
        - Snapshot 을 무조건 통짜로 저장해야함

### 3. Mongodb (=DocumentDb) 

| 평가  | 요구사항                    | 방법                                                         |
|:---:|-------------------------|------------------------------------------------------------|
| 불만족 | Partitioning            | Sharding 을 쓰려면 Mongodb, 샤딩없이 Replica 구성으로 진행하려면 DocumentDB |
| 만족  | State Schema Versioning | document 의 field 로 schema field 로 schema version 을 사용      |
| 만족  | Snapshot 생성일 조회         | 컬렉션에서 생성일에 대해 ascending 인덱스 추가                             |

- 예상 Collections
    - snapshot_collection
      - PK field 를 추가하고 unique index 생성
      - Schema field 를 추가하고 PK + Schema 를 compound index 생성
      - CreatedAt field 를 추가하고, 이 필드에 asc index 생성


- 장/단점
    - 장점
        - State 의 스키마를 바로 Document 로 저장 가능
        - 조회 요구 사항에 DDB 보다 자유로움
        - 특정 Field 만 업데이트 가능
    - 단점
        - 제약된 쓰기 성능
        - 읽기가 많다면 Read Replica 를 늘려야함
        - 고가용성 확보는 DDB, S3 보다 못함