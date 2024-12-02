package implementation

import "Redis/resp"

func HandlePing() []string{
	return []string{resp.ToSimpleString("PONG")}
}

func HandleEcho(str string) []string{
	return []string{resp.ToSimpleString(str)}
}