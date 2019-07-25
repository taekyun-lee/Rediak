package main

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

// TODO : Read configuration from files
const(
	SHARDNUM = 32
	DEFAULTLISTSIZE = 10
)


//type mapwithmutex struct{
//	d map[string]Item
//	sync.RWMutex
//}

type Item struct{
	v interface{}
	ttl int64
}

type DBInstance struct{
	mu sync.RWMutex

	bucket sync.Map
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




}



//func (db *DBInstance)LockAndReturnMutex(key string) (*sync.RWMutex, error){
//	db.mu.Lock()
//	defer db.mu.Unlock()
//	if db.IsExists(key){
//		shardmap := db.getShard(key)
//		shardmap.RWMutex.Lock()
//		return &shardmap.RWMutex, nil
//	}
//	return nil, fmt.Errorf("CannotAcquireMutexerror")
//}

func New(isactive bool,interval time.Duration) DBInstance{

	db:= DBInstance{

		// TODO: Lots of config files like eviction config
		IsActiveEviction:isactive, // for config, default to passive(false)


	}

	if db.IsActiveEviction{
		db.EvictionInterval=interval
		go activeEviction(&db)

	}


	return db
}



// GET-done, SET-done , DELETE-done, ISEXIST, CLOSE


func (db *DBInstance)Get(key string) (Item, error){


	iv, ok := db.bucket.Load(key)
	if !ok {

		return Item{}, fmt.Errorf("KeyNotExists")
	}
	v :=iv.(Item)
	fmt.Println("Internal key value ", key, v)


	// Passive eviction
	if !db.IsActiveEviction && v.ttl >0 && v.ttl < time.Now().UnixNano(){

		return Item{}, fmt.Errorf("KeyItemExpired")

	}

	return v, nil

}


func (db *DBInstance)Set(key string, value interface{}, ttl int64) {


	d := Item{
		v:value,
		ttl:ttl,
	}

	db.bucket.Store(key, d)
	fmt.Println(" SET Internal key value ", key, d)


}

func (db *DBInstance)Delete(key string) error{

	if _, ok := db.bucket.Load(key);!ok{
		return fmt.Errorf("KeyItemNotExists")
	}


	db.bucket.Delete(key)
	return nil
}

func (db *DBInstance)IsExists(key string) bool{


	hv, ok := db.bucket.Load(key)
	if !ok{
		return ok
	}
	v := hv.(Item)
	if (v.ttl > 0){
		return ok && (v.ttl > time.Now().UnixNano())
	}


	return ok
}

func activeEviction(db *DBInstance){
	ticker := time.NewTicker(db.EvictionInterval)

	for{
		select{
		case <-ticker.C:
			//for _,shardmap := range db.bucket{
			//
			//	for k,v := range shardmap.d{
			//		if v.ttl > 0 &&  v.ttl < time.Now().UnixNano(){
			//			delete(shardmap.d, k)
			//		}
			//	}
			//
			//	shardmap.Unlock()
			//
			//
			//
			//}

			db.bucket.Range(func(key,value interface{}) bool{

						v := value.(Item)
						if v.ttl > 0 &&  v.ttl < time.Now().UnixNano(){
							db.bucket.Delete(key)
							fmt.Println("asdfsadfsd")
							return true

						}
						return false
			},
			)
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

		hv, ok := db.bucket.Load(k)
		v := hv.(Item)
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
	}


	return ret, nil

}


func (db *DBInstance)MSet(key []string, value []interface{}, ttl int64) {
	for i:=0;i<len(key);i++{


		d := Item{
			v:value[i],
			ttl:ttl,
		}
		db.bucket.Store(key[i],d)

	}

}

func (db *DBInstance)AtomicIncr(key string, delta int64) (int64, error){
	db.mu.Lock()
	defer db.mu.Unlock()
	hv,ok := db.bucket.Load(key)
	v := hv.(Item)
	if !ok{
		return 0, fmt.Errorf("Keynotexists")
	}
	if val, castok := v.v.(int64); castok {
		val +=1
		v.v = val
		return v.v.(int64), nil
	}
	return 0, fmt.Errorf("unexpectederror")

}


func (db *DBInstance)Close() {

	db.mu.Lock()
	db.activeeviction<- struct{}{}
	db.bucket = sync.Map{}
	// TODO : Log db close
	// TODO : Other db close func
	db.mu.Unlock()

}






