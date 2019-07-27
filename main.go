package main

import (
	"flag"
	"fmt"
)

func main() {
	flag.Parse()
	err := make(chan error)

	go (func() {
		err <- initRespServer()
	})()


	if err := <-err; err != nil {
		fmt.Println(err.Error())
	}
}
