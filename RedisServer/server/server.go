package server

import (
	"Redis/handler"
	"Redis/myConfig"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

func GetSplitArray(data []byte) []interface{}{
	rawMsg := string(data)
	splitMsg := strings.Split(rawMsg, "\r\n")
	var interf []interface{}
	for _, val := range splitMsg {
		interf = append(interf, val)
	} 
	return interf
}

func createServer(config *myConfig.Config){
	listenAddr := fmt.Sprintf("%s:%v", config.Host, config.Port)
	server, err := net.Listen("tcp", listenAddr)
	if err != nil {
		log.Println("Error when creating Server, ",err)
		os.Exit(1)
	}
	defer server.Close()

	fmt.Printf("Redis %s is listening on port: %v\n", config.Role, config.Port)

	for {
		conn, err := server.Accept()
		if err != nil {
			log.Println("Error accepting connection, ",err)
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
				log.Println("Client disconnected abruptly.")
				break
			} else {
				log.Println("Socket error:", err)
				break
			}
		}

		data := buf[:bufLen]
		interf := GetSplitArray(data)
		fmt.Println(interf)
		dataToSendBack := handler.RequestHandler(conn, interf, config)
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
	if config.Role == "slave" {
		go handShakeWithMaster(config)
	}
	createServer(config)
}