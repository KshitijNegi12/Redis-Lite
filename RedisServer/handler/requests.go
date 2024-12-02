package handler

import (
	"Redis/implementation"
	"Redis/myConfig"
	"Redis/resp"
	"net"
	"strings"
)

func RequestHandler(conn net.Conn, data []byte, config *myConfig.Config) [] string{
	rawMsg := string(data)
	splitMsg := strings.Split(rawMsg, "\r\n")
	var interf []interface{}
	for _, val := range splitMsg {
		interf = append(interf, val)
	} 

	parsedCmds := resp.ParseMessage(interf)

	cmd := parsedCmds.Cmd
	args := parsedCmds.Args

	switch cmd {
		case "PING":
			if len(args) == 0{
				return implementation.HandlePing()
			}
			return resp.HandleErrors()
		
		case "ECHO":
			if len(args) > 0{
				str,ok := args[0].(string)
				if ok {
					return implementation.HandleEcho(str)
				}
			}
			return resp.HandleErrors()
	}

	return resp.HandleErrors()
}