package main

import (
	"sync"
	"time"
)

var DB DBInstance


type PersistentStore struct {
	sync.RWMutex

	targetDB       *DBInstance
	savePath       string
	filepermission int
	lastSaved      time.Time

	isSaveSucceed    bool
	isAsync        bool

	numChanged int32

}

func (p *PersistentStore) SaveSnapshot(){

	p.Lock()
	p.isSaveSucceed = false
	p.Unlock()





}

