package main

import (
	"github.com/go-redis/redis"
	"reflect"
)

func main(){
	redisdb := redis.NewRing(&redis.RingOptions{
		Addrs: map[string]string{
			"shard1": ":7000",
			"shard2": ":7001",
			"shard3": ":7002",
		},
	})

}
