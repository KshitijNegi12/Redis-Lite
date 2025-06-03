package store

import (
	"net"
	"time"
)

var StoredKeys = make(map[interface{}]interface{})
var ExpiryKeys = make(map[interface{}]time.Time)
var MultiQueue = make(map[net.Conn][]interface{})

func AddConnToMultiQueue(conn net.Conn) bool {
	if _, exists := MultiQueue[conn]; exists {
		return false
	} 
	MultiQueue[conn] = []interface{}{}
	return true
}

func AddConnCmdsToQueue(conn net.Conn, data []interface{}) {
	MultiQueue[conn] = append(MultiQueue[conn], data)
}

func GetQueuedCmds(conn net.Conn) []interface{} {
	cmds := MultiQueue[conn]
	delete(MultiQueue, conn)
	return cmds
}

func DiscardQueueCmds(conn net.Conn) {
	delete(MultiQueue, conn)
}

func CheckConnInQueue(conn net.Conn) bool{
	if _, exists := MultiQueue[conn]; exists {
		return true
	} 
	return false
}