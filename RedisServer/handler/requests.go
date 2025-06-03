package handler

import (
	"Redis/implementation"
	"Redis/myConfig"
	"Redis/resp"
	"Redis/store"
	"fmt"
	"net"
	"strings"
)

func RequestHandler(conn net.Conn, data []interface{}, config *myConfig.Config) [] string{
	parsedCmds := resp.ParseMessage(data)
	cmd := strings.ToUpper(parsedCmds.Cmd)
	if _, exists := store.MultiQueue[conn]; exists && cmd != "EXEC" && cmd != "DISCARD" {
		store.AddConnCmdsToQueue(conn, data)
		return resp.ToSimpleString("Queued")
	} 
	args := parsedCmds.Args

	switch cmd {
		case "PING":
			if len(args) == 0{
				return implementation.HandlePing()
			}
			return resp.HandleErrors()
		
		case "ECHO":
			if len(args) > 0{
				return implementation.HandleEcho(args)
			}
			return resp.HandleErrors()

		case "SET":
			if len(args) == 2 || len(args) == 4{
				return implementation.HandleSet(args, config)
			}
			return resp.HandleErrors()
		
		case "GET":
			if len(args) == 1 {
				return implementation.HandleGet(args)
			}
			return resp.HandleErrors()

		case "DEL":
			if len(args) == 1 {
				return implementation.HandleDel(args)
			}
			return resp.HandleErrors()

		case "INFO":
			if len(args) == 1 {
				return implementation.HandleInfo(args, config)
			}
			return resp.HandleErrors()

		case "REPLCONF":
			if len(args) == 2{
				return implementation.HandleReplconf()
			}
			return resp.HandleErrors()

		case "PSYNC":
			if len(args) == 2{
				return implementation.HandlePsync(conn, config)
			}
			return resp.HandleErrors()
		
		case "TYPE":
			if len(args) == 1{
				return implementation.HandleType(args)
			}
			return resp.HandleErrors()

		case "XADD":
			if len(args) >= 4 && len(args) % 2 == 0{
				return implementation.HandleXadd(args)
			}
			return resp.HandleErrors()
		
		case "XREAD":
			if len(args) == 3 {
				return implementation.HandleXread(args)
			}
			return resp.HandleErrors()

		case "INCR":
			if len(args) == 1{
				return implementation.HandleIncr(args)
			}
			return resp.HandleErrors()

		case "MULTI":
			if len(args) == 0{
				return implementation.HandleMulti(conn)
			}
			return resp.HandleErrors()

		case "EXEC":
			if !store.CheckConnInQueue(conn){
				return resp.ToSimpleError("ERR EXEC without MULTI")
			}

			commands := store.GetQueuedCmds(conn)
			fmt.Println(commands)
			for _, cmd := range commands {
				iCmd, ok := cmd.([]interface{})
				if !ok {
					conn.Write([]byte(fmt.Sprintf("%v", "Invalid command format")))
					continue
				}
				fmt.Println("cmd: ",iCmd)
				data := RequestHandler(conn, iCmd, config)
				fmt.Println("data: ",data)
				for _, parts := range data {
					conn.Write([]byte(fmt.Sprintf("%v", parts)))
				}
			}
			return resp.ToSimpleString("Done")

		case "DISCARD":
			if len(args) == 0{
				return implementation.HandleDiscard(conn)
			}
			return resp.HandleErrors()
		
		default:
			return resp.HandleErrors()

	}
}