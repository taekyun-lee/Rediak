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

//SET <key> <value> [<TTL "millisecond">] .
//SETEX key seconds value
func (b *Bucket) SET(c RESPContext){
	// stronglock enabled
	if *Stronglock{
		b.Lock()
		defer b.Unlock()
	}
	lenargs := len(c.args)
	var ttl time.Duration
	if lenargs <2 {
		c.WriteError(fmt.Sprintf("%s SET/SETEX has at least 2 arguments SET/SETEX  key value [ttl in milisecond]",ErrArgsLen))
		return
	}

	k,v := c.args[0], c.args[1]

	if lenargs > 2 {
		t, ok := strconv.Atoi(c.args[2])
		if ok != nil {
			c.WriteError(fmt.Sprintf("%s ttl is positive integer [ttl in milisecond]",ErrArgsLen))
			return
		}else{
			ttl = time.Duration(t) * time.Millisecond
		}

	}else{
		ttl = 0
	}

	b.Store(k,&Data{
		D:       v,
		TTL:     time.Now().Add(ttl).UnixNano(),
		dtype:   0,
		expired: false,
	})
	//
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
		//c.WriteError(fmt.Sprintf("%s key  %s not exists ",ErrNotExists,k))
		// Not exists
		c.WriteNull()
		return
	}
	dv := v.(*Data)
	if *evictionInterval == 0 && dv.TTL < time.Now().UnixNano(){
		//passive key eviction
		b.Delete(k)
		// key expired error -> not exists
		c.WriteNull()
		return
	}
	c.WriteString(dv.D.(string))
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
		_,deleted = b.LoadOrStore(c.args[i],nil)
		if deleted{
			delval+=1
		}
		// b.Delete(c.args[i]) if memory exists
	}
	atomic.AddInt32(&b.changedNum,1)
	if b.modifyCallback != nil{
		b.modifyCallback(c)
	}
	c.WriteInt(delval)
	return
}

//EXISTS <key>
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
		t, ok := strconv.Atoi(c.args[2])
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
	if *evictionInterval == 0 && dv.TTL < time.Now().UnixNano(){
		//passive key eviction
		b.Delete(k)
		c.WriteError(fmt.Sprintf("%s key expired ",ErrExpired))
		return

	}

	atomic.AddInt64(dv.D.(*int64),incrby)

	atomic.AddInt32(&b.changedNum,1)
	if b.modifyCallback != nil{
		b.modifyCallback(c)
	}


	c.WriteInt(int(dv.D.(int64) + incrby))
	return
}



//MSET <key1> <value1> [<key2> <value2> ...].

//MGET <key1> [<key2> ...].






