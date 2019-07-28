package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"sync"
	"time"
)

func Serializemap(m sync.Map)([]byte, error){
	tmpmap := make(map[string]interface{})
	m.Range(func(k, v interface{}) bool {
		v2,_ :=json.Marshal(v)
		tmpmap[k.(string)] = v2
		return true
	})
	return json.Marshal(tmpmap)

}


func UnSerializemap(data []byte)(sync.Map, error){
	var tmpmap sync.Map
	var a map[string][]byte
	var v2 Data

	e := json.Unmarshal(data,&a)
	if e!=nil{
		fmt.Println(e.Error())
		panic(e)
	}
	for k,v := range(a){

		e =json.Unmarshal(v,&v2)

		if e!=nil{
			fmt.Println(e.Error())
			panic(e)
		}
		if v2.TTL!= 0 && v2.TTL <time.Now().UnixNano(){
			// key expired
			continue
		}

		tmpmap.Store(k,v2)
	}
	return tmpmap, nil
}


func (b *Bucket) SaveSnapshot(filepath string) error{
	d, err := Serializemap(b.Map)

	if err != nil{
		fmt.Println(err.Error())
		return err
	}

	err = ioutil.WriteFile(filepath,d,0644)
	if err != nil{
		fmt.Println(err.Error())
		return err
	}
	return nil
}

func (b *Bucket)RestoreSnapshot(filepath string) error {
	d,err := ioutil.ReadFile(filepath)
	if err != nil{
		fmt.Println(err.Error())
		return err
	}
	m, err2:= UnSerializemap(d)
	if err2 != nil{
		fmt.Println(err2.Error())
		return err2
	}
	b.Map = m
	return nil


}