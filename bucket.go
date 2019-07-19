package KangDB

import (
	"fmt"
	"sync"
	"time"
)

// TODO : Read configuration from files
const(
	SHARDNUM = 32
)



type shardmap []*mapwithmutex



type mapwithmutex struct{
	d map[string]Item
	sync.RWMutex
}

type Item struct{
	v interface{}
	ttl int64
}

type DBInstance struct{
	mu sync.RWMutex

	bucket shardmap
	shardNum int

	//// TESTING
	//wg sync.WaitGroup
	//pool sync.Pool

	// TODO: Eviction config

	IsActiveEviction bool
	activeeviction chan struct{}
	EvictionInterval time.Duration


	// TODO: Snapshot config
	//IsAsyncSnapshot bool
	//SnapshotInterval time.Duration


	// TODO:Consistent hashing config
	hf Hashfunc

}

func newShardmap (shardNum int) shardmap{
	m := make(shardmap,shardNum)
	for i:=0;i<shardNum;i++{
		m[i] = &mapwithmutex{
			d:make(map[string]Item),

		}
	}
	return m
}



func New(isactive bool,interval time.Duration) DBInstance{

	db:= DBInstance{
		bucket:newShardmap(SHARDNUM),
		hf:crc32Hash,
		// TODO: Lots of config files like eviction config
		IsActiveEviction:isactive, // for config, default to passive(false)


	}

	if db.IsActiveEviction{
		db.EvictionInterval=interval
		go activeEviction(&db)

	}


	return db
}

func (db *DBInstance)GetShard(key string) *mapwithmutex{
	return db.bucket[db.hf(key)%uint32(SHARDNUM)]
}

// GET-done, SET-done , DELETE-done, ISEXIST, CLOSE


func (db *DBInstance)Get(key string) (Item, error){
	shardmap := db.GetShard(key)
	shardmap.RLock()
	v, ok := shardmap.d[key]
	if !ok {
		shardmap.RUnlock()
		return Item{}, fmt.Errorf("KeyNotExists")
	}

	// Passive eviction
	if !db.IsActiveEviction && v.ttl >0 && v.ttl < time.Now().UnixNano(){
		shardmap.RUnlock()
		return Item{}, fmt.Errorf("KeyItemExpired")

	}
	shardmap.RUnlock()
	return v, nil

}


func (db *DBInstance)Set(key string, value interface{}, ttl int64) {
	shardmap := db.GetShard(key)

	d := Item{
		v:value,
		ttl:ttl,
	}
	shardmap.Lock()
	shardmap.d[key] = d

	shardmap.Unlock()

}

func (db *DBInstance)Delete(key string) error{
	shardmap := db.GetShard(key)
	if _, ok := shardmap.d[key];!ok{
		return fmt.Errorf("KeyItemNotExists")
	}

	shardmap.Lock()
	delete(shardmap.d,key)
	shardmap.Unlock()
	return nil
}

func (db *DBInstance)IsExists(key string) bool{
	shardmap := db.GetShard(key)

	_, ok := shardmap.d[key]
	return ok
}

func activeEviction(db *DBInstance){
	ticker := time.NewTicker(db.EvictionInterval)

	for{
		select{
		case <-ticker.C:
			for _,shardmap := range db.bucket{
				shardmap.Lock()
				for k,v := range shardmap.d{
					if v.ttl > 0 &&  v.ttl < time.Now().UnixNano(){
						delete(shardmap.d, k)
					}
				}

				shardmap.Unlock()



			}
			case <-db.activeeviction:
				return
		}
	}

}

func (db *DBInstance)GracefulCloseDB() {

	db.mu.Lock()
	db.activeeviction<- struct{}{}
	db.bucket = shardmap{}
	// TODO : Log db close
	// TODO : Other db close func
	db.mu.Unlock()

}






