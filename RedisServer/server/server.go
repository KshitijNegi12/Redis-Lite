package server

import (
	"Redis/myConfig"
	"Redis/handler"
	"fmt"
	"net"
	"os"
)

func createServer(config *myConfig.Config){
	listenAddr := fmt.Sprintf("%s:%v", config.Host, config.Port)
	server, err := net.Listen("tcp", listenAddr)
	if err != nil {
		fmt.Println("Error when create Server, ",err)
		os.Exit(1)
	}
	defer server.Close()

	fmt.Printf("Redis %s is listening on port: %v\n", config.Role, config.Port)

	for {
		conn, err := server.Accept()
		if err != nil {
			fmt.Println("Error accepting connection, ",err)
			continue
		}

		go handleConnection(conn, config)
	}
}

func handleConnection(conn net.Conn, config *myConfig.Config){
	defer conn.Close()
	defer deleteConnFromServer(conn, config)

	if config.Role == "master" {
		config.Connections[conn] = true
	}

	buf := make([]byte, 1024)

	for {
		bufLen, err := conn.Read(buf)
		if err != nil {
			if err.Error() == "EOF" {
				fmt.Println("Client disconnected abruptly.", err)
				break
			} else {
				fmt.Println("Socket error:", err)
				break
			}
		}

		data := buf[:bufLen]
		dataToSendBack := handler.RequestHandler(conn, data, config)
		for _, parts := range dataToSendBack {
			conn.Write([]byte(fmt.Sprintf("%v", parts)))
		}
	}
}

func deleteConnFromServer(conn net.Conn, config *myConfig.Config) {
	if config.Role == "master" {
		delete(config.ConnectedSlaves, conn)
		delete(config.Connections, conn)
	}
}

func Start(config *myConfig.Config){
	createServer(config)
}