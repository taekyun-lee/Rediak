package KangDB

import (
	"fmt"
	"log"
	"sync"
	"time"
)

type DB struct{

	mu sync.RWMutex
	bucks *Bucket // TODO : split bucket for lock delay

	logger *log.Logger // TODO : For logging algorithm


	// For Persistence
	snapshot interface{} // TODO :  For snapshot(persistence) implementation Interface struct
	snapshotInverval time.Duration



	// TTL Config
	IsTTLActive bool // TTL algorithm active/passive
	ActiveExpireInterval time.Duration
	CloseExpire chan struct{}

	//SplitBuckNum int // TODO : split bucket for lock delay


	// Serializer
	s Serializer



	// Consistence hashing part TODO: Consistence hashing
	// Consistence hashing part

}



type Bucket struct{

	mu sync.RWMutex

	d sync.Map // Use as map[string]item{}

}

type item struct{
	Value []byte
	TTL int64
	dtype int // string string  0 bytestream []byte 1 numbers int64 2 list_of_string []string  4 hashmap map[string]string 8 and etc

}
// get(raw) put(raw) delete isexists atomic-incr/decr close_bucket


func (db *DB) Get(key string) (*item, error){
	v, ok := db.bucks.d.Load(key)
	var d item
	if ok != true || v != nil {
		// TODO: Logging [KeyNotExistError]
		return nil, fmt.Errorf("[KeyNotExistError] data with key %s not exists.\n ",key)
	}


	err := db.s.Unmarshal(v.([]byte),&d)
	if err != nil{
		// TODO: Logging [UnmarshalFailedError]
		return nil, fmt.Errorf("[UnmarshalFailedError] data with key %s Failed unmarshalling.\n ",key)

	}

	// Passive eviction
	if !db.IsTTLActive && d.TTL < time.Now().UnixNano() {
		// Delete key
		db.bucks.d.Delete(key)
		// TODO: Logging [Passive_KeyExpires]
		return nil, fmt.Errorf("[Passive_KeyExpires] data with key %s expired and cannot use.\n ",key)
	}

	return &d, nil

}


func (db *DB) Set(key string, v interface{}, ttl time.Duration, dtype int)  error{
	marv, ok := db.s.Marshal(v)
	if ok != nil{
		// TODO: Logging [marshalFailedError]
		return fmt.Errorf("[marshalFailedError] data with key %s Failed marshalling.\n ",key)
	}

	d := item{
		Value:marv,
		TTL:time.Now().Add(ttl).UnixNano(),
		dtype:dtype,

	}
	db.bucks.d.Store(key,d)

	return nil

}

func (db *DB) Delete(key string) {
	db.bucks.d.Delete(key)
}

func (db *DB) IsExist(key string) bool{
	_,ok := db.bucks.d.Load(key)
	return ok
}

func (db *DB) AtomicIncr(key string, delta int) error{

}


func (db *DB) AtomicDecr(key string, delta int) error{

}

func (db *DB) KickBucket() error{
	db.bucks = &Bucket{}
	return nil
}

func activeExpire(db DB) {
	ticker := time.NewTicker(db.ActiveExpireInterval)
	defer ticker.Stop()

	for {
		select {

		case <-ticker.C:
			now := time.Now().UnixNano()

			db.bucks.d.Range(func(key interface{}, value interface{}) bool {
				ttlitem := value.(item).TTL
				if ttlitem < now {
					db.bucks.d.Delete(key)
				}
				return true
			})
		case <-db.CloseExpire:
			return
		}
	}
}

