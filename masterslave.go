package main

import (
	"fmt"
	"time"
)

type Msg struct{
	msg string
	msgtime int64
}


type Master struct{
	name string
	addr string
	//db DBInstance

	slavelist []Slave


	mastercmd chan struct{}
}


type Slave struct{
	name string
	addr string
	//db DBInstance

	master string
	cmdchan chan struct{}
	stopchan chan struct{}

}

func NewSlave(name,addr, master string) Slave{
	return Slave{
		name:name,
		addr:addr,
		master:master,
		cmdchan:make(chan struct{}),
		stopchan:make(chan struct{}),

	}
}

func (s Slave) SlaveMain(){
	ticker := time.NewTicker(1*time.Second)
	var msgbuf Msg

	for {
		select{
		case msgbuf = <- s.cmdchan:
			fmt.Println(msgbuf.msg, msgbuf.msgtime)

		case <-ticker.C:
			fmt.Println("I'm Whitehand man! with master :",s.master)

		case <-s.stopchan:
			fmt.Println("I'm dead man! with master :",s.master)
			return
		}


	}
}

func NewMaster(name,addr string, slavelist []slave)Master{
	return Master{

	}
}