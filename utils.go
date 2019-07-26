package main

import (
	"hash/crc32"
	"hash/fnv"
)

type Hashfunc func(key string) uint32


func crc32Hash(key string) uint32{
	v := crc32.ChecksumIEEE([]byte(key))
	return v

}


func fnv32Hash(key string) uint32{
	fnv32 := fnv.New32a()
	_,err :=fnv32.Write([]byte(key))
	if err!=nil{
		// TODO : Logger Panic w/ error [fnv32HashFailed]
		panic("[fnv32HashFailed]\n")
	}
	return fnv32.Sum32()


}


