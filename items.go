package KangDB

import (
	"fmt"
	"sync"
	"time"
)


type Bucket struct {

	items sync.Map

	cleanup chan struct{}
	cleaningInterval time.Duration

	//TODO: NOT YET IMPLEMENTED -> Master-slave mechanism
	//isMaster bool
	//slaveList []*Bucket
	//slaveNameList []string

	// Expiration mechanism selector
	isActivelyExpire bool


}

type item struct {

	value interface{}


	//TODO: NOT YET IMPLEMENTED -> metadata
	//createdAt   time.Time
	//accessedAt  time.Time
	//accessCount int64

	//Additional value
	ttl int64
	//dtype byte // 0 for string 1 for bytestream 2 for list 4 for hashmap 8 for sortedmap
	//closefunc []func(v interface{})
}


func (b *Bucket) New(Tinterval time.Duration) *Bucket{

	newBucket := &Bucket{
		cleanup : make(chan struct{}),
		cleaningInterval: Tinterval,
	}

	if newBucket.isActivelyExpire{

		go func() {
			ticker := time.NewTicker(newBucket.cleaningInterval)
			defer ticker.Stop()
			for {
				select{
				case <- ticker.C:
					now := time.Now().UnixNano()
					newBucket.items.Range( func(key,value interface{}) bool{
						item := value.(item)
						if item.ttl < now{
							newBucket.items.Delete(key)
							fmt.Printf("%s : Expired data with key: %s with ttl : %d/n",time.Now().Format("2006-01-02T15:04:05.999999-07:00"),key,item.ttl)

						}
						return true
					})


				case <- newBucket.cleanup:
					return
				}

			}
		} ()

	}
	return newBucket

}


func (b *Bucket) Get(key interface{}) (interface{}, error){

	if v, ok := b.items.Load(key); ok {
		fmt.Printf("%s : Load data with key: %s is : %s/n",time.Now().Format("2006-01-02T15:04:05.999999-07:00"),key,v)
		if b.isActivelyExpire == false{
			item := v.(item)
			if item.ttl >0 && item.ttl < time.Now().UnixNano() {
				return nil, fmt.Errorf("ItemExpiredError")
			}
		}
		return v, nil
	}

	// load failed, not found
	return nil, fmt.Errorf("NotFoundError")

}

func (b *Bucket) Set(key, value interface{}, ttl int64){

	b.items.Store(key, item{
		value:value,
		ttl: time.Now().Add(time.Duration(ttl)).UnixNano(),
	})
	fmt.Printf("%s : Store k:%s v:%s with ttl %d ",time.Now().Format("2006-01-02T15:04:05.999999-07:00"),key,value,ttl)

}

func (b *Bucket) Delete(key interface{}) {
	b.items.Delete(key)
	fmt.Printf("%s : Delete k:%s  \n",time.Now().Format("2006-01-02T15:04:05.999999-07:00"),key)
}

func (b *Bucket) Close(){
	b.cleanup <- struct{}{}
	b.items = sync.Map{}
}



