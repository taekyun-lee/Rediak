// Copyright 2018 The Redix Authors. All rights reserved.
// Use of this source code is governed by a Apache 2.0
// license that can be found in the LICENSE file.
// Modify by Taekyun Lee 2019
package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/tidwall/redcon"
)

func initRespServer() error {
	db := New(false, time.Duration(10)*time.Second)
	return redcon.ListenAndServe(
		*respport,
		func(conn redcon.Conn, cmd redcon.Command) {
			// handles any panic
			defer (func() {
				if err := recover(); err != nil {
					conn.WriteError(fmt.Sprintf("fatal error: %s", (err.(error)).Error()))
				}
			})()
			ctx := (conn.Context()).(map[string]interface{})
			todo := strings.TrimSpace(strings.ToLower(string(cmd.Args[0])))
			args := []string{}
			for _, v := range cmd.Args[1:] {
				v := strings.TrimSpace(string(v))
				args = append(args, v)
			}

			ctx["db"] = db

			// internal ping-pong
			if todo == "ping" {
				conn.WriteString("PONG")
				return
			}

			// close the connection
			if todo == "quit" {
				conn.WriteString("OK")
				err := conn.Close()
				if err != nil{
					conn.WriteError(fmt.Sprintf("close error [%s]", todo))
					return
				}
				return
			}

			// find the required command in our registry
			fn := CMDLIST[todo]
			if nil == fn {
				conn.WriteError(fmt.Sprintf("unknown commands [%s]", todo))
				return
			}
			fmt.Println(todo ,args)
			// dispatch the command and catch its errors
			fn(CmdContext{
				Conn:   conn,
				cmd: todo,
				args:   args,
				db:     &db,
			})
		},
		func(conn redcon.Conn) bool {
			//accept or denied
			// use for auth
			conn.SetContext(map[string]interface{}{})
			return true
		},
		func(conn redcon.Conn, err error)  {
			//close
			// use for closing db
			conn.SetContext(map[string]interface{}{})
		},
	)
}

func main() {
	err := make(chan error)

	go (func() {
		err <- initRespServer()
	})()


	if err := <-err; err != nil {
		fmt.Println(err.Error())
	}
}