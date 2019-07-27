package main

type rediakfunc func(Bucket, RESPContext)

var (
	rediak_cmds = map[string]rediakfunc{

		// string, numerical value
		"get":    Bucket.GET,
		"set":    Bucket.SET,
		"setex":  Bucket.GET,
		"del":    Bucket.DELETE,
		"exists": Bucket.EXISTS,
		"incr":   Bucket.INCR,
		"incrby": Bucket.INCR,

		// list


		// hashmap


		// set

		// sortedset

		// util
		"gc" : Bucket.GC,



	}
)
