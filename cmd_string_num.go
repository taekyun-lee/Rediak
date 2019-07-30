
// Original code : https://github.com/alash3al/redix
// Modifications copyright (C) 2019 Taekyun Lee

package main

import (
	"fmt"
	"strconv"
	"sync/atomic"
	"time"
)


/*

Strings


TTL <key> returns -1 if key will never expire, -2 if it doesn't exists (expired), otherwise will returns the seconds remain before the key will expire.
KEYS [<regexp-pattern>]

*/

//SET key value [expiration EX seconds|PX milliseconds]
//SETEX key seconds value
func (b *Bucket) SET(c RESPContext){
	// stronglock enabled
	if *Stronglock{
		b.Lock()
		defer b.Unlock()
	}
	lenargs := len(c.args)
	var ttl int64
	if lenargs <2 {
		c.WriteError(fmt.Sprintf("%s SET/SETEX has at least 2 arguments SET/SETEX  key value [ttl in milisecond]",ErrArgsLen))
		return
	}

	k,v := c.args[0], c.args[1]

	if lenargs >2 {
		if lenargs !=4 {
			c.WriteError(fmt.Sprintf("%s SET key value [expiration EX seconds|PX milliseconds]",ErrArgsLen))
			return
		}
		xmode := c.args[2]
		if xmode =="EX"{
			t, ok := strconv.Atoi(c.args[3])
			if ok != nil {
				logger.Infoln(fmt.Sprintf("%s ttl is positive integer [ttl in milisecond]",ErrArgsLen))
				c.WriteNull()

				return
			}else{
				ttl = time.Now().Add(time.Duration(t) * time.Second).UnixNano()
			}
		}else if xmode=="PX"{
			t, ok := strconv.Atoi(c.args[3])
			if ok != nil {
				logger.Infoln(fmt.Sprintf("%s ttl is positive integer [ttl in milisecond]",ErrArgsLen))
				c.WriteNull()

				return
			}else{
				ttl = time.Now().Add(time.Duration(t) * time.Millisecond).UnixNano()
			}
		}else{
			c.WriteError(fmt.Sprintf("%s SET key value [expiration EX seconds|PX milliseconds]",ErrArgsLen))
			return
		}


	}else{
		ttl = 0
	}

	b.Store(k,&Data{
		D:       v,
		TTL:     ttl,
		dtype:   0,
		expired: false,
	})



	atomic.AddInt32(&b.changedNum,1)
	c.WriteString("OK")
	if b.modifyCallback != nil{
		b.modifyCallback(c)
	}
	return

}
//GET <key> [<default value>] .
func (b *Bucket) GET (c RESPContext){
	lenargs := len(c.args)
	if lenargs != 1 {

		c.WriteError(fmt.Sprintf("%s GET has exact 1 arguments GET key ",ErrArgsLen))
		return
	}

	k := c.args[0]

	v, ok := b.Load(k)
	if !ok {
		//logger.Infoln(fmt.Sprintf("%s key  %s not exists ",ErrNotExists,k))
		// Not exists
		c.WriteNull()
		return
	}
	dv := v.(*Data)
	if *evictionInterval == 0 && dv.TTL != 0  && dv.TTL < time.Now().UnixNano(){
		//passive key eviction

		//logger.Infoln(fmt.Sprintf("%s key  %s  expired ",ErrExpired,k))
		b.Delete(k)
		// key expired error -> not exists
		c.WriteNull()
		return
	}
	c.WriteBulkString(dv.D.(string))
	return


}
// DEL <key1> [<key2> ...]
func (b *Bucket) DELETE (c RESPContext){

	// stronglock enabled
	if *Stronglock{
		b.Lock()
		defer b.Unlock()
	}

	lenargs := len(c.args)
	var deleted bool
	delval :=0
	if lenargs < 1 {
		c.WriteError(fmt.Sprintf("%s GET has at least 1 arguments DEL key [key] ",ErrArgsLen))
		return
	}

	for i:=0;i<lenargs;i++{
		_,deleted = b.Load(c.args[i])
		if deleted{
			b.Delete(c.args[i])
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

//EXISTS key [key ...]
func (b *Bucket) EXISTS (c RESPContext){
	lenargs := len(c.args)
	existval :=0
	if lenargs < 1 {
		c.WriteError(fmt.Sprintf("%s EXISTS has at least 1 arguments EXISTS key [key]",ErrArgsLen))
		return
	}
	for i:=0;i<lenargs;i++{
		_, ok := b.Load(c.args[i])
		if ok{
			existval+=1
		}
	}
	c.WriteInt(existval)
	return

}

//INCR <key> [<by>]
//INCRBY key increment

func (b *Bucket) INCR (c RESPContext){

	// stronglock enabled
	if *Stronglock{
		b.Lock()
		defer b.Unlock()
	}

	lenargs := len(c.args)
	var incrby int64
	incrby = 1
	if lenargs < 1 {
		c.WriteError(fmt.Sprintf("%s INCR/INCRBY has at least 1 arguments INCR/INCRBY key [increment]",ErrArgsLen))
		return
	}

	k := c.args[0]
	if lenargs == 2 {
		t, ok := strconv.Atoi(c.args[1])
		if ok != nil {
			c.WriteError(fmt.Sprintf("%s increment is numeric value ",ErrArgsLen))
			return
		}else{
			incrby = int64(t)
		}
	}
	v, ok := b.Load(k)
	if !ok {
		c.WriteError(fmt.Sprintf("%s key  %s not exists ",ErrNotExists,k))
		return
	}
	dv := v.(*Data)
	if *evictionInterval == 0 && dv.TTL != 0  && dv.TTL < time.Now().UnixNano(){
		//passive key eviction
		b.Delete(k)
		c.WriteError(fmt.Sprintf("%s key expired ",ErrExpired))
		return

	}
	orig,castok :=strconv.ParseInt(dv.D.(string), 10,64)

	if castok != nil{
		c.WriteError(fmt.Sprintf("key has no numerical value "))
		return
	}

	newd := orig+incrby
	newds := strconv.FormatInt(newd,10)
	b.Store(k, &Data{
		D:       newds,
		TTL:     dv.TTL,
		dtype:   0,
		expired: false,
	})
	atomic.AddInt32(&b.changedNum,1)
	if b.modifyCallback != nil{
		b.modifyCallback(c)
	}


	c.WriteInt(int(newd))
	return
}
// EXPIRE key second
func (b *Bucket)EXPIRE(c RESPContext){
	lenargs := len(c.args)

	if lenargs < 2 {
		c.WriteError(fmt.Sprintf("%s EXPIRE has at least 2 arguments EXPIRE key second",ErrArgsLen))
		return
	}

	v, ok := b.Load(c.args[0])
	if !ok{
		c.WriteInt(1)
		return
	}
	ttl, castok := strconv.Atoi(c.args[1])
	if castok!=nil{
		c.WriteInt(0)
		return
	}
	v.(*Data).TTL = time.Now().Add(time.Duration(ttl)*time.Second).UnixNano()
	c.WriteInt(1)
	return
}
// TTL key
func (b *Bucket)TTL(c RESPContext){
	lenargs := len(c.args)

	if lenargs < 1 {
		c.WriteError(fmt.Sprintf("%s TTL has 1 arguments TTL key",ErrArgsLen))
		return
	}

	v,ok := b.Load(c.args[0])
	if !ok{
		c.WriteInt(-2)
		return
	}
	if v.(*Data).TTL ==0{
		c.WriteInt(-1)
		return
	}
	realttl :=v.(*Data).TTL
	remains := time.Unix(realttl-time.Now().UnixNano(),0).Second()
	c.WriteInt(int(remains))
	return
}






