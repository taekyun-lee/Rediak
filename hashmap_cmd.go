package main

import (
	"strconv"
)

//
// hsetCommand - HSET <HASHMAP> <KEY> <VALUE> <TTL>
// hgetCommand - HGET <HASHMAP> <KEY>
// hdelCommand - HDEL <HASHMAP> [<key1> <key2> ...]
// hgetallCommand - HGETALL <HASHMAP>
// hkeysCommand - HKEYS <hashmap>
// hmsetCommand - HMSET <HASHMAP> <key1> <val1> [<key2> <val2> ...]
// hexistsCommand - HEXISTS <HASHMAP> [<key>]
// hlencommand - HLEN hashmap
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

func hashmap_hset(c CmdContext) {
	var ttl int
	var err error
	lenargs := len(c.args)
	if lenargs <3 {
		c.WriteError("[CmdArgError]HSET cmd has at least 3 arguments (HSET hkey k v [ttl] ) ")
		return
	}
	hk, k, v := c.args[0],c.args[1],c.args[2]
	if lenargs >3 {
		ttl,err = strconv.Atoi(c.args[3])
		if err != nil{
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
func hashmap_hget(c CmdContext) {

}
func hashmap_hdel(c CmdContext) {

}
func hashmap_hgetall(c CmdContext) {

}
func hashmap_hkeys(c CmdContext) {

}
func hashmap_hmset(c CmdContext) {

}
func hashmap_hexists(c CmdContext) {

}
func hashmap_hincr(c CmdContext) {

}
func hashmap_httl(c CmdContext) {

}
func hashmap_hlen(c CmdContext) {

}
