// Copyright 2018 The Redix Authors. All rights reserved.
// Use of this source code is governed by a Apache 2.0
// license that can be found in the LICENSE file.
// Modify by Taekyun Lee 2019
package main

import (
	"fmt"
	"strconv"
	_"time"
)

/*

Strings
SET <key> <value> [<TTL "millisecond">] .
MSET <key1> <value1> [<key2> <value2> ...].
GET <key> [<default value>] .
MGET <key1> [<key2> ...].
DEL <key1> [<key2> ...]
EXISTS <key>
INCR <key> [<by>]
TTL <key> returns -1 if key will never expire, -2 if it doesn't exists (expired), otherwise will returns the seconds remain before the key will expire.
KEYS [<regexp-pattern>]

*/

func string_set(c CmdContext) {

	var ttl int
	var err error
	if len(c.args) < 2 {
		c.WriteError("[CmdArgError]SET cmd has at least 2 elements (SET key value [ttl])")
		return
	}

	k, v := c.args[0], c.args[1]

	if len(c.args) > 2 {
		ttl, err  = strconv.Atoi(c.args[2])
		if err!=nil{
			c.WriteError(err.Error())
			return
		}
	}else{
		ttl = 0
	}


	ttlval:= int64(ttl)
	c.db.Set(k,v,ttlval)


	c.WriteString("OK")

}

func string_get(c CmdContext) {

	if len(c.args) != 1{
		c.WriteError("[CmdArgError]GET cmd has exact 2 elements (GET key )")
		return
	}
	key := c.args[0]
	v, err := c.db.Get(key)
	fmt.Println(v)
	if err != nil{
		c.WriteError(err.Error())
		return
	}
	d, ok := v.v.(string)
	if !ok{
		c.WriteError(err.Error())
		return
	}

	c.WriteString(d)

}


func string_mset(c CmdContext) {

	lenargs := len(c.args)
	if lenargs %2 != 0 {
		c.WriteError("[CmdArgError]MSET cmd has even argument (k1 v1 k2 v2...) (MSET key1 value1 [key valuen])")
		return
	}
	var kl []string
	var vl []interface{}
	for i :=0;i<lenargs/2;i++{
		kl = append(kl,c.args[i*2])
		vl = append(vl,c.args[i*2+1])
	}


	c.db.MSet(kl,vl, 0)


	c.WriteString("OK")

}

func string_mget(c CmdContext) {
	lenargs := len(c.args)

	if lenargs < 1{
		c.WriteError("[CmdArgError]MGET cmd at least 1 elements (GET key [key2] )")
		return
	}
	//vl := make([]string, lenargs)
	iteml,err := c.db.MGet(c.args)

	if err != nil{
		c.WriteError(err.Error())
		return
	}

	c.WriteArray(lenargs)
	for _,d := range iteml{
		if sd, ok := d.v.(string);ok{
			c.WriteBulkString(sd)
		}else{
			c.WriteNull()

		}

	}

}


func string_del(c CmdContext){
	lenargs := len(c.args)

	if lenargs < 1{
		c.WriteError("[CmdArgError]DEL cmd at least 1 elements (DEL key [key2] )")
		return
	}
	for _,k := range c.args{
		if err := c.db.Delete(k);err!=nil{
			c.WriteError(err.Error())
			return
		}
	}
	c.WriteString("OK")
}

func string_exists(c CmdContext){
	lenargs := len(c.args)

	if lenargs != 1{
		c.WriteError("[CmdArgError]EXIST cmd exact 1 elements (EXISTS key )")
		return
	}

	if ok:=c.db.IsExists(c.args[0]);ok{
		c.WriteString("1")
		return
	}

	c.WriteString("0")
	return

}


func string_incr(c CmdContext){
	var err error
	var delta int
	lenargs := len(c.args)

	if lenargs < 1{
		c.WriteError("[CmdArgError]INCR cmd at least 1 elements (INCR key [value] )")
		return
	}
	key := c.args[0]
	if lenargs > 1{
		delta,err = strconv.Atoi(c.args[1])
		if err != nil{
			c.WriteError(err.Error())
			return
		}
	}else{
		delta=10
	}

	v,err := c.db.AtomicIncr(key,int64(delta))
	if err != nil{
		c.WriteError(err.Error())
		return
	}

	c.WriteInt64(v)

}


func string_ttl(c CmdContext) {
	if len(c.args) != 1 {
		c.WriteError("[CmdArgError]TTL cmd has exact 1 elements (ttl key )")
		return
	}
	key := c.args[0]
	v, err := c.db.Get(key)

	if err != nil {
		c.WriteError(err.Error())
		return
	}
	c.WriteInt64(v.ttl)
}




