package main

import "flag"

var (

	respaddr = "localhost"
	respport = flag.String("resp-addr", ":6380", "the address of resp server")

	DEFAULTHASHSIZE = 32
	DEFAULTTTLVALUE = int64(0) // Set TTL to infinite if not specified


)
