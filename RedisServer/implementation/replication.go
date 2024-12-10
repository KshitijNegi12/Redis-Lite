package implementation

import (
	"Redis/myConfig"
	"Redis/resp"
)

func HandleReplconf() []string {
	return  resp.ToSimpleString("OK")
}

func HandlePsync(config *myConfig.Config) []string {
	return []string {resp.ToRESP([]interface{}{"FULLRESYNC", config.MasterReplID, config.MasterReplOffset}), }
}