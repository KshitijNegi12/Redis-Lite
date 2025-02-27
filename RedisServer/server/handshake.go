package server

import (
	"Redis/myConfig"
	"Redis/resp"
	"fmt"
	"log"
	"net"
)

func handShakeWithMaster(config *myConfig.Config) {
	address := fmt.Sprintf("%v:%v",config.MasterHost, config.MasterPort)
	connWithMaster, err := net.Dial("tcp", address)
	if err != nil {
		log.Println("Error in connecting to Master, ",err)
		return
	}
	
	_, err = connWithMaster.Write([]byte(resp.ToRESP("PING")))
	if err != nil {
		log.Println("Error in sending message to Master, ",err)
		return
	}
	
	go func (){
		defer connWithMaster.Close()
		toSendCapa  := true
		buf := make([]byte, 1024)

		for {
			n, err := connWithMaster.Read(buf)
			if err != nil {
				log.Println("Error reading from connection:", err)
				return
			}

			data := buf[:n]
			splitData := GetSplitArray(data)
			response := resp.ParseMessage(splitData)
			fmt.Println(response)
			cmd := response.Cmd

			switch cmd {
				case "PONG":
					_, err = connWithMaster.Write([]byte(resp.ToRESP([]interface{}{"REPLCONF", "listening-port", config.Port})))

				case "OK":
					if toSendCapa {
						_, err = connWithMaster.Write([]byte(resp.ToRESP([]interface{}{"REPLCONF", "capa", "psync2"})))
						toSendCapa = false
					} else {
						_, err = connWithMaster.Write([]byte(resp.ToRESP([]interface{}{"PSYNC", "?", -1})))
					}
				
				case "FULLRESYNC":
					if len(response.Args) == 2 {
						config.MasterReplID = response.Args[0].(string)
					} 
					
				default:
					log.Println("Unknown command received:", cmd)
			}

			if err != nil {
				log.Println("Error in sending message to Master:", err)
				return
			}
		}
	}()
}
