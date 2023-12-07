package main

import (
	"dq/node"
	"dq/tcp"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func main() {
	var err error

	// 1. Logger init
	f, err := os.OpenFile(filepath.Dir(os.Args[0])+"/logs/server.log", os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		fmt.Printf("Error opening file: %v", err)
	}
	log.SetOutput(f)
	defer f.Close()

	// 2. Queue init
	queue := &node.Queue{}
	queue.Init()

	// 3. TCP Server init
	tcpServer := &tcp.Server{Parser: &tcp.Parser{}, Queue: queue}
	tcpServer.Init()

	err = tcpServer.Start("", "12000")
	if err != nil {
		log.Fatalln("Error listening:", err.Error())
	}

	defer tcpServer.StopAndClose()

	log.Println("Listening...")
	tcpServer.Listen()
}
