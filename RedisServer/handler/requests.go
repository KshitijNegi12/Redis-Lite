package handler

import (
	"Redis/implementation"
	"Redis/myConfig"
	"Redis/resp"
	"strings"
	"net"
)

func RequestHandler(conn net.Conn, data []interface{}, config *myConfig.Config) [] string{
	parsedCmds := resp.ParseMessage(data)
	cmd := strings.ToUpper(parsedCmds.Cmd)
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
				return implementation.HandleSet(args)
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
				return implementation.HandlePsync(config)
			}
			return resp.HandleErrors()
	}

	return resp.HandleErrors()
}