package KangDB

import (
	"encoding/binary"
	"fmt"
	"sync"
	"time"
)

type Bucket struct {

	//Bucket name
	name string

	// Main hash table in this bucket : Where Items are contained, thread-safe by go sync package
	// Key and values are arbitrary value, and type as interface{} in Go.
	items sync.Map

	//Mutex for some reason
	mu sync.RWMutex

	// For cleanup channel for messaging
	closeExpire chan struct{}
	// Cleaning interval (when select active invalidation) in seconds
	// If this value is nil, passive invalidation(in GET/SET method) will be activated
	expiringInterval time.Duration

	//TODO: NOT YET IMPLEMENTED -> Master-slave mechanism
	//isMaster bool
	//slaveList []*Bucket
	//slaveNameList []string

}

type item struct {

	//Real value of
	value interface{}
	// 0 for string (string) 1 for byte list(binary) []byte  2 for list([]string) 4 for hashmap(map[string]string) 8 for sortedmap
	dtype byte
	//Time-To-Live in seconds
	ttl int64

	//TODO: NOT YET IMPLEMENTED -> time metadata (if needed)
	//createdAt   time.Time
	//accessedAt  time.Time
	//accessCount int64

	//TODO: Callback function when this item expired/deleted (if needed)
	//closefunc []func(v interface{})
}

func New(cinteval time.Duration, name string) *Bucket {

	var b *Bucket
	b = &Bucket{
		name:name,
		expiringInterval: cinteval,
		//TODO: Not implemented
	}

	if cinteval !=0{
		// Active expiration with items.ttl
		// Goroutine activeExpire

		go activeExpire(b)
	}

	return b

}

func activeExpire(b *Bucket) {
	ticker := time.NewTicker(b.expiringInterval)
	defer ticker.Stop()

	for {
		select {

		case <-ticker.C:
			now := time.Now().UnixNano()

			b.items.Range(func(key interface{}, value interface{}) bool {
				ttlitem := value.(item).ttl
				if ttlitem < now {
					b.items.Delete(key)
				}
				return true
			})
		case <-b.closeExpire:
			return
		}
	}
}

func (b *Bucket) GetInterval() time.Duration {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.expiringInterval
}

//func (b *Bucket) SetInterval(t time.Duration) {
//	b.mu.Lock()
//	defer b.mu.Unlock()
//	if b.expiringInterval !=0 {
//		if t ==0{
//			//Kill active expiration goroutine
//			b.closeExpire <- struct{}{}
//
//		}
//		b.expiringInterval = t
//	}else{
//		if t !=0{
//			b.expiringInterval = t
//			go activeExpire(b)
//		}
//	}
//
//}

func (b *Bucket) GET(key interface{}) (interface{},error){
	v,ok := b.items.Load(key)
	if !ok{
		// Value not found
		return nil, fmt.Errorf("NotfoundError")
	}
	it:=v.(item)

	if it.ttl >0 && it.ttl < time.Now().UnixNano(){
		return nil, fmt.Errorf("TimeExpired")
	}


	return v, nil

}

func (b *Bucket) SET(key, value interface{}, ttl int64, dtype byte){
	var exp int64

	if ttl >0{
		//Set with expiration time
		exp = time.Now().Add(time.Second * time.Duration(ttl)).UnixNano()
	}else{
		// Will Not be expired
		exp=0
	}

	v := item{
		value:value,
		dtype:dtype,
		ttl:exp,
	}
	b.items.Store(key,v)

}

func (b *Bucket) DELETE(key interface{}) error {
	b.items.Delete(key)
	return nil
}


func (b *Bucket) CLOSE()  {
	b.closeExpire <- struct{}{}
	b.items = sync.Map{}
}


func (b *Bucket) REFRESH(key interface{},ttl int64) error {
	v, ok := b.items.Load(key)
	var exp int64
	if !ok {
		// Value not found
		return fmt.Errorf("NotfoundError")
	}
	it := v.(item)

	if ttl >0{
		//Set with expiration time
		exp = time.Now().Add(time.Second * time.Duration(ttl)).UnixNano()
	}else{
		// Will Not be expired
		exp=0
	}
	// Refresh item's TTL with current time and parameter
	it.ttl = exp

	b.items.Store(key,it)
	return nil
}

func (b *Bucket) CONTAINS(key interface{}) bool {
	_,ok:=b.items.Load(key)
	return ok
}

func (b *Bucket) BUCKETSIZE() int{
	return binary.Size(b.items)
}



