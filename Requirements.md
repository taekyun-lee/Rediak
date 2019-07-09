#2. Distribute in-memory key-value store 만들기 
- Distributed System을 구성하는데 필요한 지식들을 습득하고 각각의 trade-off를 비교하며 시스템을 설계해보기

- [필수 구현 사항] 
    - standalone in-memory kv store의 OOD
    - kv store의 최소한의 프로토콜 정의 (get, set, delete...), CLI로 시작해도 좋으나 networking을 지원

-  [구현 시 우대 사항] 
    - value type이 collection인 entry 지원(list, hashmap 등) 많을수록 좋음 (easy)
    - entry에 TTL 지원 (medium)
    - primitive key-value entry에 대해 concurrent operation 지원 (test and set) (medium)
    - persistent 지원 (medium)
    - multithreaded 모델로 변경 시도 (redis는 싱글 스레드), 같은 키에 대해서는 serializability 지원 (medium~hard)
    - distributed로 확장, fail-over 고려 (hard)

- [심사 기준] 
    - 각 클래스가 객체지향 원칙에 따라서 잘 설계되어있는가?
    - 필수 사항을 만족하는가?
    - 우대 사항을 만족하는가? (각각 난이도에 따른 점수 부여, 설계에 따른 점수 부여)
    - trade-off를 고려하는 과정에 대한 문서
    - 완성된 시스템의 performance test 문서 (테스트 전에 용량 산정도 하면 추가 점수, cpu, mem, network io, r/w ops에 대해서)

- [참고 자료]
    - https://redis.io