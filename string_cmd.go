package KangDB

import (
	"strconv"
	_"time"
)

/*

Strings
SET <key> <value> [<TTL "millisecond">]
MSET <key1> <value1> [<key2> <value2> ...]
GET <key> [<default value>]
MGET <key1> [<key2> ...]
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
		c.WriteError("SET command requires at least two arguments: SET <key> <value> [TTL Millisecond]")
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
	// if *flagACK {
	c.db.Set(k,v,ttlval)
	// } else {
	// 	kvjobs <- func() {
	// 		c.db.Set(k, v, ttlVal)
	// 	}
	// }

	c.WriteString("OK")

}

func string_get(db *DBInstance, key string)(string, error) {
	i, err := db.Get(key)
	if err != nil{
		return "", err
	}
	return i.v.(string), nil
}

