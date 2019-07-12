package KangDB

import (
	"bytes"
	"io/ioutil"
	"sync"
	"time"
	"github.com/klauspost/compress/zstd"
	"encoding/gob"
)

type snapshot interface{

	LoadSnapshot(path string) (interface{}, error)
	SaveSnapshot(path string, d interface{}) error
	LatestRecord() (path string, when time.Time)
	IsCompressed() bool
	IsSuccessed() bool

}


type ZippedSnapshot struct{
	mu sync.RWMutex

	lastSuccessTime time.Time
	lastSuccessflag bool
	lastSuccesspath string

	isCompressed bool
	compressAlgorithm string
}

func (z *ZippedSnapshot) SaveSnapshot(path string, d interface{})error{
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(d)
	if err != nil{
		return err
	}
	w,err := zstd.NewWriter(nil)
	if err != nil{
		return err
	}
	src := buf.Bytes()
	res := w.EncodeAll(src,make([]byte,0,len(src)))

	err = ioutil.WriteFile("./tmp/tmpfile", res, 0644)
	if err != nil{
		return err
	}
	return nil


}

