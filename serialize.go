package main

import (
	"bytes"
	"encoding/gob"
)

type Serializer interface{

	Marshal(d interface{}) ([]byte, error) //data to []byte

	Unmarshal(src []byte, dst interface{}) error // []byte to data inteface and save to dst

}


type GobSerializer struct{}

func(g GobSerializer)Marshal(d interface{}) ([]byte, error){
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(d)

	// Toss error control responsibility to others ~~~
	return buf.Bytes(), err
}


func(g GobSerializer)UnMarshal(src []byte, dst interface{})error{
	buf := bytes.NewBuffer(src)
	dec := gob.NewDecoder(buf)
	err := dec.Decode(dst)
	return err

}



