package main

import (
	"fmt"
	"log"
	"sync"
)
var (
	testbucket = Bucket{
		RWMutex:          sync.RWMutex{},
		Map:              sync.Map{},
		changedNum:       0,
		closechan:        nil,
		evictionCallback: nil,
		modifyCallback:   nil,
	}
)

func init(){
	s := &Data{
		D:       "I wanna go home",
		TTL:     0,
		dtype:   0,
		expired: false,
	}
	m := &Data{
		D:       map[string]string{
			"a":"I",
			"b":"Wanna",
			"c":"Go",
			"d":"home",
		},
		TTL:     0,
		dtype:   0,
		expired: false,
	}
	l := &Data{
		D:       []string{"i","wanna","go","home"},
		TTL:     0,
		dtype:   0,
		expired: false,
	}
	testbucket.Map.Store("s",s)
	testbucket.Map.Store("l",l)
	testbucket.Map.Store("m",m)



}

func main(){
	err := testbucket.SaveSnapshot("./savesnap.json")
	if err !=nil{
		log.Fatal(err.Error())

	}
	err = testbucket.RestoreSnapshot("./savesnap.json")
	if err !=nil{
		log.Fatal(err.Error())

	}

	fmt.Println(testbucket.Map.Load("s"))
	fmt.Println(testbucket.Map.Load("l"))
	fmt.Println(testbucket.Map.Load("m"))

	a, _ := testbucket.Map.Load("l")
	fmt.Println(a.(Data).D)

}