package main

import (
	"flag"
	"fmt"
	"github.com/sirupsen/logrus"
	"runtime"
)

func init(){
	flag.Parse()
	runtime.GOMAXPROCS(*numCore)

	logger.SetLevel(logrus.DebugLevel)


	// TODO : SWIM and consistent hashing
	//modifydb := new(sync.Map)


}

func main() {

	db := NewBucket()

	err := make(chan error)
	// TODO : consistent hashring
	//hashring := make(chan struct{})
	//torediak := make(chan struct{})

	fmt.Println(rediaklogo)
	go (func() {
		err <- initRespServer(&db)
	})()

	go(func() {
		// TODO : rdb snapshot implementation
	})()

	go(func() {
		// TODO : consistent hashring
	})()

	if err := <-err; err != nil {
		logger.Println(err.Error())
		fmt.Println(err.Error())
	}
}
