package KangDB

import (
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

	if len(c.args) != 2{
		c.WriteError("[CmdArgError]GET cmd has exact 2 elements (GET key )")
		return
	}
	key := c.args[0]
	v, err := c.db.Get(key)

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
	kl := make([]string, int(lenargs/2))
	vl := make([]interface{}, int(lenargs/2))

	for i :=0;i<lenargs;i+=2{
		kl[i] = c.args[i]
		vl[i] = c.args[i+1]
	}


	c.db.MSet(kl,vl, 0)


	c.WriteString("OK")

}

func string_mget(c CmdContext) {
	lenargs := len(c.args)

	if lenargs < 2{
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


