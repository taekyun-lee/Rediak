// https://github.com/alash3al/redix
// Thanks
package KangDB


import "github.com/tidwall/redcon"

type CmdInterface func(c CmdContext)

type CmdContext struct{
	redcon.Conn
	db KVinterface
	cmd string
	args []string
	retval interface{}
}

type KVinterface interface {
	Get(key string) (Item, error)
	Set(key string, value interface{}, ttl int64)
	MGet(key []string) ([]Item, error)
	MSet(key []string, value []interface{}, ttl int64)
	AtomicIncr(key string, by int64) (int64, error)
	Delete(key string) error
	IsExists(key string) bool
	GC()
	Close()

}

var (
	CMDLIST= map[string]CmdInterface{
"set":    string_set,
"mset":   string_mset,
"get":    string_get,
"mget":   string_mget,
"del":    string_del,
"exists": string_exists,
"incr":   string_incr,
"ttl": string_ttl,

// lists
"lpush":      lpushCommand,
"lpushu":     lpushuCommand,
"lrange":     lrangeCommand,
"lrem":       lremCommand,
"lcount":     lcountCommand,
"lcard":      lcountCommand,
"lsrch":      lsearchCommand,
"lsrchcount": lsearchcountCommand,

// hashes
"hset":    hsetCommand,
"hget":    hgetCommand,
"hdel":    hdelCommand,
"hgetall": hgetallCommand,
"hkeys":   hkeysCommand,
"hmset":   hmsetCommand,
"hexists": hexistsCommand,
"hincr":   hincrCommand,
"httl":    httlCommand,
"hlen":    hlenCommand,

// utils
"gc":       gcCommand,
"info":     infoCommand,
"echo":     echoCommand,
"flushdb":  flushdbCommand,
"flushall": flushallCommand,
}
)