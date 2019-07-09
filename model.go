package KangDB

import (
	"reflect"
	"sync"
	"time"
)

type Bucket struct {
	data map[string]Nodevalue

	mu sync.Mutex

}

//func (b *Bucket) GetData(key string) (v value, err error){
//
//}
//


type Nodevalue struct {


	// Payload
	v interface{}

	//Metadata
	createdAt  time.Time
	modifiedAt time.Time
	ttl        uint64
	size       uint64
	owner      uint64
	dataType   uint8

}




type pString struct{//datatype 0
	v []byte
}
type pList struct {//datatype 1
	v [][]byte
}

type pHashmap struct{//datatype 2
	v map[string]byte
}


func getType(myvar interface{}) string {
	return reflect.TypeOf(myvar).Name()
}



