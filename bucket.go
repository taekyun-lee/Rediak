package main

import (
	"fmt"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

// TODO : Read configuration from files
const(
	SHARDNUM = 32
	DEFAULTLISTSIZE = 10
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

func (db *DBInstance)getShard(key string) *mapwithmutex{
	return db.bucket[db.hf(key)%uint32(SHARDNUM)]
}

// GET-done, SET-done , DELETE-done, ISEXIST, CLOSE


func (db *DBInstance)Get(key string) (Item, error){
	shardmap := db.getShard(key)
	shardmap.RLock()
	v, ok := shardmap.d[key]
	fmt.Println("Internal key value ", key, v)
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
	shardmap := db.getShard(key)

	d := Item{
		v:value,
		ttl:ttl,
	}
	shardmap.Lock()
	shardmap.d[key] = d
	fmt.Println(" SET Internal key value ", key, d)

	shardmap.Unlock()

}

func (db *DBInstance)Delete(key string) error{
	shardmap := db.getShard(key)
	if _, ok := shardmap.d[key];!ok{
		return fmt.Errorf("KeyItemNotExists")
	}

	shardmap.Lock()
	delete(shardmap.d,key)
	shardmap.Unlock()
	return nil
}

func (db *DBInstance)IsExists(key string) bool{
	shardmap := db.getShard(key)

	v, ok := shardmap.d[key]
	if v.ttl > 0{
		return ok && (v.ttl > time.Now().UnixNano())
	}
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

func (db *DBInstance)GC() {
	// USE WITH CAUTION, IT BLOCKS ENTIRE DB!!!!!
	runtime.GC()
}

func (db *DBInstance)MGet(key []string) ([]Item, error){
	ret := make([]Item,len(key))
	for i,k :=range key{
		shardmap := db.getShard(k)
		shardmap.RLock()
		v, ok := shardmap.d[k]

		if !ok {
			ret[i] = Item{}
			//shardmap.RUnlock()
			//return Item{}, fmt.Errorf("KeyNotExists")
		}

		// Passive eviction
		if !db.IsActiveEviction && v.ttl >0 && v.ttl < time.Now().UnixNano(){
			ret[i] = Item{}
			//shardmap.RUnlock()
			//return Item{}, fmt.Errorf("KeyItemExpired")

		}
		ret[i] = v
		shardmap.RUnlock()
	}


	return ret, nil

}


func (db *DBInstance)MSet(key []string, value []interface{}, ttl int64) {
	for i:=0;i<len(key);i++{
		shardmap := db.getShard(key[i])

		d := Item{
			v:value[i],
			ttl:ttl,
		}
		shardmap.Lock()
		shardmap.d[key[i]] = d

		shardmap.Unlock()

	}

}

func (db *DBInstance)AtomicIncr(key string, delta int64) (int64, error){
	shardmap := db.getShard(key)
	v,ok := shardmap.d[key]
	if !ok{
		return 0, fmt.Errorf("Keynotexists")
	}
	if _, castok := v.v.(int64); castok {
		atomic.AddInt64(v.v.(*int64),delta)
		return v.v.(int64), nil
	}
	return 0, fmt.Errorf("Keynotexists")

}


func (db *DBInstance)Close() {

	db.mu.Lock()
	db.activeeviction<- struct{}{}
	db.bucket = shardmap{}
	// TODO : Log db close
	// TODO : Other db close func
	db.mu.Unlock()

}






