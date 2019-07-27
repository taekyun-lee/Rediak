package main

import (
	"github.com/tidwall/redcon"
	"sync"
)

type Bucket struct {
	sync.RWMutex
	// typical map version
	//d map[string]*Data

	sync.Map //map[string]*Data

	// How many times changed before take snapshot
	// use atomic func when put/delete things
	changedNum int32
	closechan  chan struct{}

	// TODO : CALL THIS FUNCTION IF DATA IS EVICT
	evictionCallback func(key string) interface{}
	modifyCallback   func(c RESPContext)
}

type Data struct {
	// POINTER OF DATA
	D       interface{}
	TTL     int64
	dtype   byte //string =0 , integer = 1, float = 2 list = 11 , hash = 21...
	expired bool
}

type RESPContext struct {
	redcon.Conn
	cmd  string
	args []string
}

var (
	ErrArgsLen     = "(Error with arguments lengths)"
	ErrInvalidArgs = "(Error with invalid argument)"
	ErrNotExists = "(Error Not exists)"
	ErrExpired = "(Error key expired)"

)
