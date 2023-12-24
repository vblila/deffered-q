package main

import (
	"dq/tcp"
	"log"
)

func main() {
	tcpServer := tcp.Server{}
	tcpServer.Init()

	err := tcpServer.Start("", "12000")
	if err != nil {
		log.Fatalln("Error listening:", err.Error())
	}

	defer tcpServer.StopAndClose()

	log.Println("Listening...")
	tcpServer.Listen()
}
