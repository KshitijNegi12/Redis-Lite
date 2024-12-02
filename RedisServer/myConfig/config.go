package myConfig

import (
	"net"
)

type Config struct {
	Host              string
	Port              int
	Role              string
	Connections       map[net.Conn]bool 
	ConnectedSlaves   map[net.Conn]bool 
	MasterReplID      string
	MasterReplOffset  int
	AckCount          int
	AckNeeded         int
	PropagationCount  int
	WaitingForAck     bool
	Cargs             map[string]string 
	MasterHost		  string
	MasterPort		  int
}

type Content struct {
	Cmd			string
	Args 		[]interface{}
}