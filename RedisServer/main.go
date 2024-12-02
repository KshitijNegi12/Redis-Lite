package main

import (
	"Redis/myConfig"
	"Redis/server"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
)

func main(){
	config := &myConfig.Config{
		Host:            "127.0.0.1",
		Port:            6379,
		Role:            "master",
		Connections:     make(map[net.Conn]bool),
		ConnectedSlaves: make(map[net.Conn]bool),
		MasterReplID:    "8371b4fb1155b71f4a04d3e1bc3e18c4a990aeeb",
		MasterReplOffset: 0,
		AckCount:         0,
		AckNeeded:        0,
		PropagationCount: 0,
		WaitingForAck:    false,
		Cargs:            make(map[string]string),	 
		MasterHost: 	  "",
		MasterPort: 	  0,
	}

	args := os.Args[1:]
	for i:= 0; i<len(args); i++ {
		arg := args[i]
		if strings.HasPrefix(arg, "--") {
			key := strings.TrimPrefix(arg, "--")
			if i+1 < len(args) {
				config.Cargs[key] = args[i+1]
				i++
			}
		}
	}
	
	port, isThere := config.Cargs["port"]
	if isThere {
		port, err := strconv.Atoi(port)
		if err != nil {
			fmt.Println("Error Wrong Port Specified, ", err)
			os.Exit(1)
		}
		config.Port = port
		delete(config.Cargs, "port")
	}

	id, isThere := config.Cargs["master_replid"]
	if isThere {
		config.MasterReplID = id
		delete(config.Cargs, "master_replid")
	}

	masterInfo, isThere := config.Cargs["replicaof"]
	if isThere {
		config.Role = "slave"
		parts := strings.Split(masterInfo, " ")
		config.MasterHost = parts[0]
		port, err := strconv.Atoi(parts[1])
		if err != nil {
			fmt.Println("Error Wrong Port Specified of Master, ", err)
			os.Exit(1)
		}
		config.MasterPort = port
		delete(config.Cargs, "replicaof")
	}

	// fmt.Println(config)
	server.Start(config)
}