// Copyright 2018 The Redix Authors. All rights reserved.
// Use of this source code is governed by a Apache 2.0
// license that can be found in the LICENSE file.
package KangDB

import (
	"fmt"
	"log"
	"strings"

	"github.com/tidwall/redcon"
)

func initRespServer() error {
	return redcon.ListenAndServe(
		respaddr,
		func(conn redcon.Conn, cmd redcon.Command) {
			// handles any panic
			defer (func() {
				if err := recover(); err != nil {
					conn.WriteError(fmt.Sprintf("fatal error: %s", (err.(error)).Error()))
				}
			})()

			// fetch the connection context
			// normalize the todo action "command"
			// normalize the command arguments
			ctx := (conn.Context()).(map[string]interface{})
			todo := strings.TrimSpace(strings.ToLower(string(cmd.Args[0])))
			args := []string{}
			for _, v := range cmd.Args[1:] {
				v := strings.TrimSpace(string(v))
				args = append(args, v)
			}




			// set the default db if there is no db selected
			if ctx["db"] == nil || ctx["db"].(string) == "" {
				ctx["db"] = "0"
			}



			// internal ping-pong
			if todo == "ping" {
				conn.WriteString("PONG")
				return
			}

			// close the connection
			if todo == "quit" {
				conn.WriteString("OK")
				conn.Close()
				return
			}

			// find the required command in our registry
			fn := CMDLIST[todo]
			if nil == fn {
				conn.WriteError(fmt.Sprintf("unknown commands [%s]", todo))
				return
			}

			// dispatch the command and catch its errors
			fn(CmdContext{
				Conn:   conn,
				cmd: todo,
				args:   args,
				db:     db,
			})
		},
		func(conn redcon.Conn) bool {
			conn.SetContext(map[string]interface{}{})
			return true
		},
		nil,
	)
}