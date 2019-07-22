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

//type DB interface {
//	Incr(k string, by int64) (int64, error)
//	Set(k, v string, ttl int) error .
//	MSet(data map[string]string) error.
//	Get(k string) (string, error) .
//	MGet(keys []string) []string.
//	TTL(key string) int64
//	Del(keys []string) error .
//	Scan(ScannerOpt ScannerOptions) error n
//	Size() int64
//	GC() error.
//	Close().
//}
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