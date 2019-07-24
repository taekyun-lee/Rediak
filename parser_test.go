package main

import (
	"fmt"
	"reflect"
	"testing"
)

func TestParseArgs(t *testing.T) {
	setoneQ :="SET key1 value1"
	setoneA :=[]string{"SET", "key1","value1"}

	v,ok := ParseArgs(setoneQ)
	fmt.Println(setoneA,v)
	if ok != nil || !reflect.DeepEqual(setoneA,v){
		t.Errorf("SET 1  parse failed %v %v\n",v,setoneA)
	}

	getoneQ :="GET keyasdf"
	getoneA :=[]string{"GET", "keyasdf"}

	v,ok = ParseArgs(getoneQ)
	fmt.Println(getoneA,v)
	if ok != nil || !reflect.DeepEqual(getoneA,v){
		t.Errorf("SET 1  parse failed %v %v\n",v,getoneA)
	}


}
