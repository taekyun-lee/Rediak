// Original code : https://github.com/alash3al/redix
// Modifications copyright (C) 2019 Taekyun Lee

package main

import (
	"fmt"
	"sync/atomic"
	"time"
)


/*

Strings


TTL <key> returns -1 if key will never expire, -2 if it doesn't exists (expired), otherwise will returns the seconds remain before the key will expire.
KEYS [<regexp-pattern>]

*/

//HSET <key> <field> <value>

func (b *Bucket) HSET(c RESPContext){
	// stronglock enabled
	if *Stronglock{
		b.Lock()
		defer b.Unlock()
	}
	lenargs := len(c.args)
	var ttl int64
	if lenargs <3 {
		c.WriteError(fmt.Sprintf("%s HSET has exast 3 arguments HSET hkey key value ",ErrArgsLen))
		return
	}

	hk, k,v := c.args[0], c.args[1], c.args[2]

	ttl = 0
	d := make(map[string]string,DEFAULTHASHSIZE)
	d[k] = v
	already, loaded := b.LoadOrStore(hk, &Data{
		D:       d,
		TTL:     ttl,
		dtype:   21,
		expired: false,
	})


	if loaded && already.(*Data).dtype==21  { // already has hashmap
		newv := already.(*Data)
		newv.D.(map[string]string)[k] = v

		//newv.(map[string]string)[k] = v
		//
		// b.Store(hk, &Data{
		//	D:       newv,
		//	TTL:     ttl,
		//	dtype:   1,
		//	expired: false,
		//})
		atomic.AddInt32(&b.changedNum,1)
		c.WriteInt(1)
		if b.modifyCallback != nil{
			b.modifyCallback(c)
		}
		return
	}else{
		b.Store(hk, &Data{
			D:       d,
			TTL:     ttl,
			dtype:   21,
			expired: false,
		})

		//
		atomic.AddInt32(&b.changedNum,1)
		c.WriteInt(1)
		if b.modifyCallback != nil{
			b.modifyCallback(c)
		}
		return
	}

}
//HGET <key> <field> [<default value>] .
func (b *Bucket) HGET (c RESPContext){
	lenargs := len(c.args)
	if lenargs != 2 {
		c.WriteError(fmt.Sprintf("%s HGET has exact 2 arguments HGET key field ",ErrArgsLen))
		return
	}
	hk := c.args[0]
	k := c.args[1]

	v, ok := b.Load(hk)
	if !ok {
		c.WriteError(fmt.Sprintf("%s key  %s not exists ",ErrNotExists,hk))

		return
	}
	if  v.(*Data).dtype != 21{
		c.WriteNull()
		return
	}
	dv := v.(*Data)
	if *evictionInterval == 0 && dv.TTL != 0  && dv.TTL < time.Now().UnixNano(){
		//passive key eviction
		c.WriteError(fmt.Sprintf("%s key  %s  expired ",ErrExpired,k))
		b.Delete(hk)
		// key expired error -> not exists
		//c.WriteNull()
		return
	}
	if v,ok:=dv.D.(map[string]string)[k] ;!ok{
		c.WriteNull()
		return
	}else{
		c.WriteBulkString(v)
		return
	}

}
// HDEL key field [field ...]
func (b *Bucket) HDELETE (c RESPContext){

	// stronglock enabled
	if *Stronglock{
		b.Lock()
		defer b.Unlock()
	}

	lenargs := len(c.args)

	delval :=0
	if lenargs < 2 {
		c.WriteError(fmt.Sprintf("%s HDEL has at least 2 arguments HDEL key field [field ...] ",ErrArgsLen))
		return
	}

	hk := c.args[0]


	v, ok := b.Load(hk)
	if !ok {
		c.WriteError(fmt.Sprintf("%s key  %s not exists ",ErrNotExists,hk))

		return
	}
	if  v.(*Data).dtype != 21{
		c.WriteNull()
		return
	}
	dv := v.(*Data)
	hashmap :=dv.D.(map[string]string)


	for i:=1;i<lenargs;i++{
		_,ok := hashmap[c.args[i]]
		if ok{
			delete(hashmap,c.args[i])
			delval+=1
		}

	}
	atomic.AddInt32(&b.changedNum,1)
	if b.modifyCallback != nil{
		b.modifyCallback(c)
	}
	c.WriteInt(delval)
	return
}

// HEXISTS key field
func (b *Bucket) HEXISTS (c RESPContext){
	lenargs := len(c.args)
	if lenargs < 2 {
		c.WriteError(fmt.Sprintf("%s HEXISTS has exists arguments HEXISTS key field",ErrArgsLen))
		return
	}

	hk := c.args[0]


	v, ok := b.Load(hk)
	if !ok {
		//c.WriteError(fmt.Sprintf("%s key  %s not exists ",ErrNotExists,hk))
		c.WriteInt(0)
		return
	}
	if  v.(*Data).dtype != 21{
		//c.WriteNull()
		c.WriteInt(0)
		return
	}
	dv := v.(*Data)
	hashmap :=dv.D.(map[string]string)
	_, ok = hashmap[c.args[1]]
	if ok{
		//delete(hashmap, c.args[1])
		c.WriteInt(1)
		return
	}

	c.WriteInt(0)
	return

}
//
////INCR <key> [<by>]
////INCRBY key increment
//
//func (b *Bucket) HINCR (c RESPContext){
//
//	// stronglock enabled
//	if *Stronglock{
//		b.Lock()
//		defer b.Unlock()
//	}
//
//	lenargs := len(c.args)
//	var incrby int64
//	incrby = 1
//	if lenargs < 1 {
//		c.WriteError(fmt.Sprintf("%s INCR/INCRBY has at least 1 arguments INCR/INCRBY key [increment]",ErrArgsLen))
//		return
//	}
//
//	k := c.args[0]
//	if lenargs == 2 {
//		t, ok := strconv.Atoi(c.args[1])
//		if ok != nil {
//			c.WriteError(fmt.Sprintf("%s increment is numeric value ",ErrArgsLen))
//			return
//		}else{
//			incrby = int64(t)
//		}
//	}
//	v, ok := b.Load(k)
//	if !ok {
//		c.WriteError(fmt.Sprintf("%s key  %s not exists ",ErrNotExists,k))
//		return
//	}
//	dv := v.(*Data)
//	if *evictionInterval == 0 && dv.TTL != 0  && dv.TTL < time.Now().UnixNano(){
//		//passive key eviction
//		b.Delete(k)
//		c.WriteError(fmt.Sprintf("%s key expired ",ErrExpired))
//		return
//
//	}
//	orig,castok :=strconv.ParseInt(dv.D.(string), 10,64)
//
//	if castok != nil{
//		c.WriteError(fmt.Sprintf("key has no numerical value "))
//		return
//	}
//
//	newd := orig+incrby
//	newds := strconv.FormatInt(newd,10)
//	b.Store(k, &Data{
//		D:       newds,
//		TTL:     dv.TTL,
//		dtype:   0,
//		expired: false,
//	})
//	atomic.AddInt32(&b.changedNum,1)
//	if b.modifyCallback != nil{
//		b.modifyCallback(c)
//	}
//
//
//	c.WriteInt(int(newd))
//	return
//}
//
//
//
////MSET <key1> <value1> [<key2> <value2> ...].
//
////MGET <key1> [<key2> ...].
//
//
//
//
//
//
