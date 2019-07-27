// Copyright 2018 The Redix Authors. All rights reserved.
// Use of this source code is governed by a Apache 2.0
// license that can be found in the LICENSE file.
// Modify by Taekyun Lee 2019

package main

import (
	"fmt"
	"strconv"
)

//
// hincrCommand - HINCR <hash> <key> [<number>]
// httlCommand - HTTL <HASH> <KEY>
//

//// hashes
//"hset":    hashmap_hset,
//"hget":    hashmap_hget,
//"hdel":    hashmap_hdel,
//"hgetall": hashmap_hgetall,
//"hkeys":   hashmap_hkeys,
//"hmset":   hashmap_hmset,
//"hexists": hashmap_hexists,
//"hincr":   hashmap_hincr,
//"httl":    hashmap_httl,
//"hlen":    hashmap_hlen,



// hsetCommand - HSET <HASHMAP> <KEY> <VALUE> <TTL>
func hashmap_hset(c CmdContext) {
	var ttl int
	var err error
	lenargs := len(c.args)
	if lenargs <3 {
		c.WriteError("[CmdArgError]HSET cmd has at least 2 arguments (HSET hkey k v [ttl] ) ")
		return
	}
	hk, k, v := c.args[0],c.args[1],c.args[2]
	if lenargs >3 {
		ttl,err = strconv.Atoi(c.args[3])
		if err != nil{
			c.WriteError("asdfdsafsdfsadf")
			c.WriteError(err.Error())
			return
		}
	}else{
		ttl = 0
	}

	newhash := make(map[string]string,DEFAULTHASHSIZE)
	newhash[k] = v
	c.db.Set(hk, newhash,int64(ttl))
	c.WriteInt(1)

}

// hgetCommand - HGET <HASHMAP> <KEY>

func hashmap_hget(c CmdContext) {

	lenargs := len(c.args)
	fmt.Println(lenargs)

	if lenargs !=2 {
		c.WriteError("[CmdArgError]HGET cmd has exact  2 arguments (HGET hkey k  ) ")
		return
	}

	hv, err := c.db.Get(c.args[0])
	if err != nil{
		c.WriteError("dadadadada")

		c.WriteError(err.Error())
		return
	}

	hmap,ok := hv.v.(map[string]string)
	if !ok{
		c.WriteError("WrongType (Required hash)")
		return
	}
	v,hok := hmap[c.args[1]]
	if !hok{
		c.WriteError("Not exist in this hashmap\n")
		return
	}

	c.WriteString(v)


}

// hdelCommand - HDEL <HASHMAP> [<key1> <key2> ...]
func hashmap_hdel(c CmdContext) {

	lenargs := len(c.args)
	if lenargs <2 {
		c.WriteError("[CmdArgError]HDEL cmd has at least  2 arguments (HDEL <HASHMAP> [<key1> <key2> ...] ) ")
		return
	}

	hk := c.args[0]
	keys := c.args[1:]

	hv, err := c.db.Get(hk)
	if err != nil{
		c.WriteError(err.Error())
		return
	}

	hmap,ok := hv.v.(map[string]string)
	if !ok{
		c.WriteError("WrongType (Required hash)")
		return
	}
	for _,v :=range keys{
		delete(hmap,v)
	}

	c.WriteString("OK")


}

// hgetallCommand - HGETALL <HASHMAP>
func hashmap_hgetall(c CmdContext) {

	lenargs := len(c.args)
	if lenargs != 1 {
		c.WriteError("[CmdArgError]HGETALL cmd has exact 1 arguments ( HGETALL <HASHMAP>) ")
		return
	}

	hk := c.args[0]

	hv, err := c.db.Get(hk)
	if err != nil{
		c.WriteError(err.Error())
		return
	}

	hmap,ok := hv.v.(map[string]string)
	if !ok{
		c.WriteError("WrongType (Required hash)")
		return
	}
	c.WriteArray(len(hmap)*2)
	for k,v :=range hmap{
		c.WriteBulkString(k)
		c.WriteBulkString(v)
	}

}
// hkeysCommand - HKEYS <hashmap>
func hashmap_hkeys(c CmdContext) {
	lenargs := len(c.args)
	if lenargs <1 {
		c.WriteError("[CmdArgError]HKEYS cmd has exact 1  arguments (HKEYS <hashmap> ) ")
		return
	}

}
// hmsetCommand - HMSET <HASHMAP> <key1> <val1> [<key2> <val2> ...]

func hashmap_hmset(c CmdContext) {
	lenargs := len(c.args)
	if lenargs <3   {
		c.WriteError("[CmdArgError]HMSET cmd has at least 3 arguments (HMSET hkey k1 v1 k2 v2...) ")
		return
	}
	hk := c.args[0]
	hmsetargs := c.args[1:]
	if len(hmsetargs) %2 != 0 {
		c.WriteError("[CmdArgError]HMSET cmd has even arguments (HMSET hkey k1 v1 k2 v2...)")
		return
	}

	newhm := make(map[string]string)
	for i :=0;i<len(hmsetargs)/2;i++{
		newhm[hmsetargs[i*2]] = hmsetargs[i*2+1]
	}


	c.db.Set(hk,newhm, DEFAULTTTLVALUE)
	c.WriteString("OK")
}

// hexistsCommand - HEXISTS <HASHMAP> [<key>]

func hashmap_hexists(c CmdContext) {
	lenargs := len(c.args)

	if lenargs != 2{
		c.WriteError("[CmdArgError]HEXIST cmd exact 2 elements (HEXISTS hkey key )")
		return
	}

	v, ok:= c.db.Get(c.args[0])
	if ok != nil{
		c.WriteString("0")
		return
	}else{
		hv,cok :=v.v.(map[string]string)
		if !cok{
			c.WriteString("0")
			return
		}
		_,isexist := hv[c.args[1]]
		if isexist {
			c.WriteString("1")
			return
		}

	}
	c.WriteString("0")
	return

}

// hlencommand - HLEN hashmap

func hashmap_hlen(c CmdContext) {


	if len(c.args) != 1{
		c.WriteError("[CmdArgError]HLEN cmd exact 1 elements (HLEN hkey key )")
		return
	}

	v,ok := c.db.Get(c.args[0])
	if ok !=nil{
		c.WriteInt(0)
		return
	}
	d,conok := v.v.(map[string]string)
	if conok{
		c.WriteInt(0)
		return
	}else{
		c.WriteInt(len(d))
		return
	}

}
// hincrCommand - HINCR <hash> <key> [<number>]
