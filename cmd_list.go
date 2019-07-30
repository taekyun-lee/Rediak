package main

import (
	"encoding/hex"
	"fmt"
	"strconv"
	"sync/atomic"
)

//LPUSH key value [value ...]
func (b *Bucket) LPUSH(c RESPContext){
	// stronglock enabled
	if *Stronglock{
		b.Lock()
		defer b.Unlock()
	}
	lenargs := len(c.args)
	if lenargs <2 {
		c.WriteError(fmt.Sprintf("%s lpush has at least 3 arguments LPUSH key value [value ...]",ErrArgsLen))
		return
	}
	k := c.args[0]
	s := make([]string,0)
	for i:=lenargs-1;i>0;i--{
		s = append(s,c.args[i])
	}

	v, ok := b.LoadOrStore(k, &Data{
		D:       s,
		TTL:     0,
		dtype:   11,
		expired: false,
	})
	if ok{
		if v.(*Data).dtype == 11{
			v.(*Data).D = append(s,v.(*Data).D.([]string)...)
			atomic.AddInt32(&b.changedNum,1)
			if b.modifyCallback != nil{
				b.modifyCallback(c)
			}
			c.WriteInt(len(v.(*Data).D.([]string)))
			return
		}
		c.WriteError("not list type")
		return
	}else{
		c.WriteInt(len(s))
		return
	}


}
//LPOP key
func (b *Bucket) LPOP(c RESPContext){
	// stronglock enabled
	if *Stronglock{
		b.Lock()
		defer b.Unlock()
	}
	lenargs := len(c.args)
	if lenargs <1 {
		c.WriteError(fmt.Sprintf("%s lpop has exact 1 arguments lpop key",ErrArgsLen))
		return
	}
	k := c.args[0]
	v, ok := b.Load(k)
	if ok{
		if v.(*Data).dtype == 11{
			var x string
			if len(v.(*Data).D.([]string)) >1{

				x, v.(*Data).D = v.(*Data).D.([]string)[0],v.(*Data).D.([]string)[1:]
			}else if len(v.(*Data).D.([]string)) ==1 {
				x , v.(*Data).D = v.(*Data).D.([]string)[0],[]string{}

			}else{
				c.WriteNull()
				return
			}

			atomic.AddInt32(&b.changedNum,1)
			if b.modifyCallback != nil{
				b.modifyCallback(c)
			}
			c.WriteString(x)
			return
		}
		c.WriteError("not list type")
		return
	}else{
		c.WriteError("not exists")
		return
	}


}

//LINDEX key index
func (b *Bucket) LINDEX(c RESPContext){
	// stronglock enabled
	if *Stronglock{
		b.Lock()
		defer b.Unlock()
	}
	lenargs := len(c.args)
	if lenargs <1 {
		c.WriteError(fmt.Sprintf("%s lindex has at least 3 arguments lindex key index",ErrArgsLen))
		return
	}
	k:= c.args[0]
	i, castok := strconv.Atoi(c.args[0])
	if castok != nil{
		c.WriteError(castok.Error())
		return
	}

	v, ok := b.Load(k)

	if ok{
		if v.(*Data).dtype == 11{
			lenlist := len( v.(*Data).D.([]string))
			if lenlist-1 >= i{
				x := v.(*Data).D.([]string)[i]
				c.WriteString(x)
				return
			}else if i == -1 {
				x := v.(*Data).D.([]string)[lenlist-1]
				c.WriteString(x)
				return
			} else{
				c.WriteNull()
				return
			}



		}
		c.WriteError("not list type")
		return
	}else{
		c.WriteError("not exists")
		return
	}


}
//LLEN key
func (b *Bucket) LLEN(c RESPContext){
	// stronglock enabled
	if *Stronglock{
		b.Lock()
		defer b.Unlock()
	}
	lenargs := len(c.args)
	if lenargs <1 {
		c.WriteError(fmt.Sprintf("%s llen has at least 3 arguments llen key",ErrArgsLen))
		return
	}
	k:= c.args[0]
	v,ok := b.Load(k)
	if !ok {
		c.WriteInt(0)
		return
	}else{
		lenlist := len(v.(*Data).D.([]string))
		c.WriteInt(lenlist)
		atomic.AddInt32(&b.changedNum,1)
		if b.modifyCallback != nil{
			b.modifyCallback(c)
		}
		return
	}
}

//LRANGE key start stop
func (b *Bucket) LRANGE(c RESPContext){
	// stronglock enabled
	if *Stronglock{
		b.Lock()
		defer b.Unlock()
	}
	lenargs := len(c.args)
	if lenargs <3 {
		c.WriteError(fmt.Sprintf("%s lrange has at least 3 arguments lrange key start stop",ErrArgsLen))
		return
	}
	k:= c.args[0]
	st, stok := strconv.Atoi(c.args[1])
	if stok != nil{
		c.WriteError(stok.Error())
		return
	}

	stop, stopok := strconv.Atoi(c.args[1])
	if stopok != nil{
		c.WriteError(stok.Error())
		return
	}
	v ,ok := b.Load(k)
	if !ok{
		c.WriteNull()
		return
	}
	vlist := v.(*Data).D.([]string)
	if st < 0 || st > len(vlist){
		c.WriteError("List out of range")
		return
	}
	if stop < 0 || stop >len(vlist){
		c.WriteError("List out of range")
		return
	}
	if st> stop{
		c.WriteError("List out of range")
		return
	}

	c.WriteArray(stop-st)
	for i:=st; i<stop;i++{
		c.WriteBulkString(hex.EncodeToString([]byte(vlist[i])))
	}
	return

}
//LREM key count value
//func (b *Bucket) LREM(c RESPContext){
//
//}

