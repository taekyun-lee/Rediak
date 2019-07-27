package main

import (
	"flag"
	"time"
)

var (

	respaddr = flag.String("rediak-addr","127.0.0.1","the address of rediak server")
	respport = flag.String("rediak-addr", ":6380", "the port of rediak server")
	evictionInterval = flag.Duration("evict-time",10*time.Second,"Default interval of eviction, 0 means no active eviction")
	defaultTTL = flag.Duration("evict-time",0,"Default TTL, 0 means never expired")
	Stronglock = flag.Bool("Strong-lock",true, "use mutex to all modification command,  ")

	// INTERNAL USAGE, TODO: RUNTIME CHANGE WHEN NEEDED
	DEFAULTHASHSIZE = 32
	DEFAULTSTRINGSIZE = 32

)




