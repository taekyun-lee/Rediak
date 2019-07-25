// https://github.com/alash3al/redix
// Thanks
package main

import "github.com/tidwall/redcon"

type CmdInterface func(c CmdContext)

type CmdContext struct {
	redcon.Conn
	db     KVinterface
	cmd    string
	args   []string
	retval interface{}
}

type KVinterface interface {
	Get(key string) (Item, error)
	Set(key string, value interface{}, ttl int64)
	MGet(key []string) ([]Item, error)
	MSet(key []string, value []interface{}, ttl int64)
	AtomicIncr(key string, delta int64) (int64, error)
	Delete(key string) error
	IsExists(key string) bool
	GC()
	Close()

	/*


		getShard(key string) *mapwithmutex
		Get(key string) (Item, error)
		Set(key string, value interface{}, ttl int64)
		MGet(key []string) ([]Item, error)
		MSet(key []string, value []interface{}, ttl int64)
		AtomicIncr(key string, delta int64) (int64, error)
		Delete(key string) error
		IsExists(key string) bool
		GC()
		Close()

	*/

}

var (
	CMDLIST = map[string]CmdInterface{

		//
		//// strings
		"set":    string_set,
		"mset":   string_mset,
		"get":    string_get,
		"mget":   string_mget,
		"del":    string_del,
		"exists": string_exists,
		"incr":   string_incr,
		"ttl":    string_ttl,

		// hashes
		"hset":    hashmap_hset,
		"hget":    hashmap_hget,
		"hdel":    hashmap_hdel,
		"hgetall": hashmap_hgetall,
		"hkeys":   hashmap_hkeys,
		"hmset":   hashmap_hmset,
		"hexists": hashmap_hexists,
		"hincr":   hashmap_hincr,
		"httl":    hashmap_httl,
		"hlen":    hashmap_hlen,

		//
		//// lists
		//"lpush":      lpushCommand,
		//"lpushu":     lpushuCommand,
		//"lrange":     lrangeCommand,
		//"lrem":       lremCommand,
		//"lcount":     lcountCommand,
		//"lcard":      lcountCommand,
		//"lsrch":      lsearchCommand,
		//"lsrchcount": lsearchcountCommand,
		//

		//
		//// utils
		//"gc":       gcCommand,
		//"info":     infoCommand,
		//"echo":     echoCommand,
		//"flushdb":  flushdbCommand,
		//"flushall": flushallCommand,
	}
)
