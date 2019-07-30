package main

import (
	"flag"
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"runtime"
	"sync/atomic"
	"time"
)

func init() {
	flag.Parse()
	runtime.GOMAXPROCS(*numCore)
	baselogger = logrus.New()
	baselogger.SetLevel(logrus.InfoLevel)

	// TODO : SWIM and consistent hashing
	//modifydb := new(sync.Map)


}

func main() {

	db := NewBucket()
	baselogger.SetOutput(os.Stdout)
	logger = baselogger.WithTime(time.Now())
	snapshotstop := make(chan struct{})
	if *restoreSnapshot != "" {
		restoreerror := db.RestoreSnapshot(*restoreSnapshot)
		if restoreerror != nil {
			logger.Log(logrus.ErrorLevel, restoreerror.Error())
			return
		}
	}

	//debug.SetGCPercent(-1)
	err := make(chan error)

	// TODO : consistent hashring

	fmt.Println(rediaklogo)
	flag.CommandLine.PrintDefaults()
	go (func() {
		err <- initRespServer(&db)
	})()

	// TODO : rdb snapshot implementation -> complete
	// TODO : rdb snapshot testing
	if *snapshotInterval != 0 {
		go (func() {

			t := time.NewTicker(time.Duration(*snapshotmodifyInterval) * time.Second)
			for {
				select {
				case <-t.C:
					db.RLock()
					defer db.RUnlock()
					if atomic.LoadInt32(&db.changedNum) >= int32(*snapshotmodifyInterval) {

						logger.Printf("[%s] =======save snapshot==============\n", time.Now().Format("2006-01-02T15:04:05.999999-07:00"))
						name := fmt.Sprintf("%s/snapshot_%s.rdb", *storageDir, time.Now().Format("2006-01-02T15:04:05.999999-07:00"))
						logger.Printf("Save complete, : %s\n", name)
						saveerr := db.SaveSnapshot(name)

						if saveerr != nil {
							logger.Log(logrus.ErrorLevel, saveerr.Error())
							return
						}
						atomic.StoreInt32(&db.changedNum, 0)
						runtime.GC()
					}
				case <-snapshotstop:
					logger.Printf("Save Stopped, : %s\n")

					return
				default:
					runtime.Gosched()
				}
			}
		})()
	}



	// periodically write info
	go func() {
		t := time.NewTicker(time.Second * time.Duration(*printInfoInterval))
		for {
			select {
			case <-t.C:
				var m runtime.MemStats
				runtime.ReadMemStats(&m)

				logger.Printf("\n\n[%s] =======info check==============\n", time.Now().Format("2006-01-02T15:04:05.999999-07:00"))
				logger.Printf("# of go processes : %d\n", runtime.GOMAXPROCS(0))
				logger.Printf("# of goroutines : %d\n", runtime.NumGoroutine())
				logger.Printf("Alloc = %v MiB\n", bToMb(m.Alloc))
				logger.Printf("\tTotalAlloc = %v MiB\n", bToMb(m.TotalAlloc))
				logger.Printf("\tSys = %v MiB\n", bToMb(m.Sys))
				logger.Printf("\tNumGC = %v\n", m.NumGC)
				logger.Printf("\tNumofModified = %v\n", atomic.LoadInt32(&db.changedNum))

			}
		}

	}()

	if err := <-err; err != nil {
		logger.Println(err.Error())
		fmt.Println(err.Error())
	}
}
