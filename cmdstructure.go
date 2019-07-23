package KangDB

const (
	// For all and (string, bitmap) dataset
	NUMCMD_GET = iota
	NUMCMD_SET
	NUMCMD_DELETE
	NUMCMD_ISEXIST

	// TODO: For numbers only
	NUMCMD_INCR
	NUMCMD_DECR
	NUMCMD_INCRBY
	NUMCMD_DECRBY

	// For list ops
	NUMCMD_LPUSH
	NUMCMD_LPOP

	// For hashmap ops
	NUMCMD_HMSET
	NUMCMD_HGET
	NUMCMD_HSET
	NUMCMD_HDEL

	// Internal command

	// TODO: For server command

	//
)
