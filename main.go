package main

import (
	"flag"
	"fmt"
	"github.com/sirupsen/logrus"
	"runtime"
	"sync/atomic"
	"time"
)

func init(){
	flag.Parse()
	runtime.GOMAXPROCS(*numCore)

	logger.SetLevel(logrus.DebugLevel)


	// TODO : SWIM and consistent hashing
	//modifydb := new(sync.Map)


}

func main() {

	db := NewBucket()

	if *restoreSnapshot != ""{
		restoreerror := db.RestoreSnapshot(*restoreSnapshot)
		if restoreerror != nil{
			logger.Log(logrus.ErrorLevel,restoreerror.Error())
			return
		}
	}

	err := make(chan error)


	fmt.Println(rediaklogo)
	go (func() {
		err <- initRespServer(&db)
	})()

	// TODO : rdb snapshot implementation -> complete
	// TODO : rdb snapshot testing
	if *snapshotInterval != 0{
		go(func() {

			t := time.Tick(time.Duration(*snapshotInterval) * time.Second)
			select {
				case <- t:
					if atomic.LoadInt32(&db.changedNum) >= int32(*snapshotmodifyInterval){
						name := fmt.Sprintf("%s/snapshot_%s.rdb",*storageDir,time.Now().Format("2006-01-02T15:04:05.999999-07:00"))
						saveerr :=db.SaveSnapshot(name)
						if saveerr != nil{
							logger.Log(logrus.ErrorLevel,saveerr.Error())
							return
						}
					}

			}
		})()
	}

	// TODO : consistent hashring
	//hashring := make(chan struct{})
	//torediak := make(chan struct{})
	go(func() {

	})()

	if err := <-err; err != nil {
		logger.Println(err.Error())
		fmt.Println(err.Error())
	}
}
