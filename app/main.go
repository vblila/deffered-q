package main

import (
	"dq/node"
	"dq/tcp"
	"log"
)

func main() {
	// 1. Queue init
	queue := &node.Queue{}
	queue.Init()

	// 2. TCP Server init
	tcpServer := &tcp.Server{Parser: &tcp.Parser{}, Queue: queue, Watcher: &node.Watcher{Queue: queue}}
	tcpServer.Init()

	err := tcpServer.Start("", "12000")
	if err != nil {
		log.Fatalln("Error listening:", err.Error())
	}

	defer tcpServer.StopAndClose()

	log.Println("Listening...")
	tcpServer.Listen()
}
