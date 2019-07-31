// Original code : https://github.com/alash3al/redix
// Modifications copyright (C) 2019 Taekyun Lee
package main

import (
	"fmt"
	"strings"

	"github.com/tidwall/redcon"
)

func initRespServer(db *Bucket) error {

	return redcon.ListenAndServe(
		fmt.Sprintf(":%d",*respport),
		func(conn redcon.Conn, cmd redcon.Command) {
			// handles any panic
			if *evictionInterval != 0{
				go db.activeEviction()
			}

			defer (func() {
				if err := recover(); err != nil {
					logger.Panic(fmt.Sprintf("fatal error: %s", (err.(error)).Error()))
					conn.WriteError(fmt.Sprintf("fatal error: %s", (err.(error)).Error()))
				}
			})()

			//ctx := (conn.Context()).(map[string]interface{})
			todo := strings.TrimSpace(strings.ToLower(string(cmd.Args[0])))
			args := []string{}
			for _, v := range cmd.Args[1:] {
				v := strings.TrimSpace(string(v))
				args = append(args, v)
			}



			// internal ping-pong
			if todo == "ping" {
				conn.WriteString("PONG")
				return
			}

			// close the connection
			if todo == "quit" {
				conn.WriteString("OK")
				err := conn.Close()
				if err != nil {
					conn.WriteError(fmt.Sprintf("close error [%s]", todo))
					return
				}
				return
			}



			// find the required command in our registry
			fn := rediak_cmds[todo]
			if nil == fn {
				conn.WriteError(fmt.Sprintf("unknown commands [%s]", todo))
				return
			}
			// dispatch the command and catch its errors
			fn(db, RESPContext{
				Conn: conn,
				cmd:  todo,
				args: args,
			})
		},
		func(conn redcon.Conn) bool {
			//accept or denied
			// use for auth
			conn.SetContext(map[string]interface{}{})
			return true
		},
		func(conn redcon.Conn, err error) {
			//close
			// use for closing db
			conn.SetContext(map[string]interface{}{})
		},
	)
}
