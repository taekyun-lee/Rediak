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

	err := make(chan error)
	fmt.Println(rediaklogo)
	go (func() {
		err <- initRespServer()
	})()


	if err := <-err; err != nil {
		fmt.Println(err.Error())
	}
}
