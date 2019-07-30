package main

import (
	"flag"
	"github.com/sirupsen/logrus"
	"runtime"
)

var (

	respaddr = flag.String("rediak-addr","127.0.0.1","the address of rediak server")
	respport = flag.Int("rediak-port", 6380, "the port of rediak server")
	evictionInterval = flag.Int("evict-interval",0,"Default interval of eviction, 0 means no active eviction")
	Stronglock = flag.Bool("Strong-lock",true, "use mutex to all modification command,  ")
	numCore = flag.Int("num-core", runtime.NumCPU(),"number of cores using this instances")
	restoreSnapshot = flag.String("restore-snapshot","","if restore needed, write absolute path of files( like /path/of/folder/file.rdb ) ")
	storageDir = flag.String("storage-dir","./","Default persistent storage location /path/to/snapshotfolder")
	snapshotInterval = flag.Int("snapshot-interval",0,"Default time interval of take snapshot, 0 means no active snapshot")
	snapshotmodifyInterval = flag.Int("snapshot-modify-interval",1,"Default time interval of take snapshot, 0 means no active eviction")
	printInfoInterval = flag.Int("printinfo-interval",1,"Default time interval of print interval")

	// TODO:consistent ring and SWIM Protocol
	//peeraddr = flag.String("peer", "", "peer to connect ring")

	// INTERNAL USAGE, TODO: RUNTIME CHANGE WHEN NEEDED
	DEFAULTHASHSIZE = 32
	DEFAULTSTRINGSIZE = 32


)


var (
	rediak_cmds = map[string]func(*Bucket,RESPContext){

		// string, numerical value
		"get":      (*Bucket).GET,
		"set":     (*Bucket).SET,
		"setex":   (*Bucket).SET,
		"del":     (*Bucket).DELETE,
		"exists":  (*Bucket).EXISTS,
		"incr":    (*Bucket).INCR,
		"incrby":  (*Bucket).INCR,

		// list
		"lpush": (*Bucket).LPUSH,
		"lpop": (*Bucket).LPOP,
		"lindex": (*Bucket).LINDEX,
		"llen": (*Bucket).LLEN,
		"lrange": (*Bucket).LRANGE,
		//"lrem": (*Bucket).LREM,

		// hashmap

		"hget":  (*Bucket).HGET,
		"hset": (*Bucket).HSET,
		"hdel": (*Bucket).HDELETE,
		"hexists":(*Bucket).HEXISTS,
		"expire":(*Bucket).EXPIRE,
		"ttl":(*Bucket).TTL,
		// set

		// sortedset

		// util
		"gc" :  (*Bucket).GCExec,



	}
)

var (
	ErrArgsLen     = "(Error with arguments lengths)"
	ErrInvalidArgs = "(Error with invalid argument)"
	ErrNotExists = "(Error Not exists)"
	ErrExpired = "(Error key expired)"

)

var baselogger *logrus.Logger
var logger *logrus.Entry

const (
	versionNum = 0.1
	introstring = "Rediak  started"
	rediaklogo=`

#####  ###### #####  #   ##   #    # 
#    # #      #    # #  #  #  #   #  
#    # #####  #    # # #    # ####   
#####  #      #    # # ###### #  #   
#   #  #      #    # # #    # #   #  
#    # ###### #####  # #    # #    # 
                                     
`


)




