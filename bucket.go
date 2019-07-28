package main

import (
	"runtime"
	"sync"
	"time"
)

//type Bucket struct{
//	sync.RWMutex
//	// typical map version
//	//d map[string]*Data
//
//	sync.Map //map[string]*Data
//
//	// How many times changed before take snapshot
//	// use atomic func when put/delete things
//	changedNum uint
//	closechan chan struct{}
//
//	// TODO : CALL THIS FUNCTION IF DATA IS EVICT
//	evictionCallback func(key string)interface{}
//	modifyCallback func(c interface{}) interface{}
//
//
//}
//
//type Data struct {
//	// POINTER OF DATA
//	D interface{}
//	TTL time.Duration
//	dtype byte
//	expired bool
//}



func NewBucket() Bucket{
	return Bucket{
		RWMutex:          sync.RWMutex{},
		Map:              sync.Map{},
		changedNum:       0,
		closechan:        make(chan struct{}),
		evictionCallback: nil,
		modifyCallback:   nil,
	}
}

func (b *Bucket)activeEviction(){
	ticker := time.NewTicker(time.Duration(*evictionInterval) * time.Second)

	for{
		select{
		case <-ticker.C:
			now := time.Now().UnixNano()
			b.Range(func(key,value interface{}) bool{
				v := value.(*Data)

				if  v.TTL > 0 && v.TTL < now {
					b.Delete(key)
				}
				return true
			},
			)
		case <-b.closechan:
			return
		}
	}

}

func (b *Bucket) GCExec(c RESPContext) {
	// USE WITH CAUTION, IT BLOCKS ENTIRE DB!!!!!
	c.WriteString("GC starts. IT BLOCKS ENTIRE DB!!!!!  ")
	runtime.GC()
}
