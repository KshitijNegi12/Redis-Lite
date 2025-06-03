package implementation

import (
	"Redis/myConfig"
	"Redis/resp"
	"strings"
	"net"
)

func HandleInfo(args []interface{}, config *myConfig.Config) []string {
	arg, ok := args[0].(string)
	if !ok {
		return resp.HandleErrors()
	}
	arg = strings.ToLower(arg)
	if(arg == "replication"){
		return resp.ToSimpleString(config.Role)
	}
	return resp.HandleErrors()
}

func HandleReplconf() []string {
	return  resp.ToSimpleString("OK")
}

func HandlePsync(conn net.Conn, config *myConfig.Config) []string {
	config.ConnectedSlaves[conn] = true
	return []string {resp.ToRESP([]interface{}{"FULLRESYNC", config.MasterReplID, config.MasterReplOffset}) }
}
