package KangDB

import (
	"sync"
	"time"
	"lz"
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

