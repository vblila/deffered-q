package main

import (
	"dq/config"
	"dq/tcp"
	"flag"
	"fmt"
	"log"
	"os"
)

func initByCmdLine() {
	host := flag.String("h", "127.0.0.1", "TCP server host")
	port := flag.String("p", "12000", "TCP server port")
	ict := flag.Uint("ict", 0, "Inactive connection time (in seconds), 0 - without limit")
	rtt := flag.Uint("rtt", 0, "Reserved task life time (in seconds) after which the watcher will delete the reserved task or add it back to the queue, 0 - disable watcher")
	rta := flag.Uint("rta", 0, "The number of attempts after which the watcher will delete the reserved task from queue, 0 - watcher delete the reserved task when life time expires")
	debug := flag.Uint("debug", 0, "Debug profiler, 1 - enable, 0 - disable")

	if len(os.Args) > 1 && os.Args[1] == "--help" {
		fmt.Printf("DefferedQ is a simple and fast work queue.\n\n")
	}

	flag.Parse()
	config.SetData(*host, *port, *debug != 0, *ict, *rtt, *rta)
}

func main() {
	initByCmdLine()

	tcpServer := tcp.Server{}
	tcpServer.Init()

	err := tcpServer.Start(config.Host, config.Port)
	if err != nil {
		log.Fatalln("Error listening:", err.Error())
	}

	defer tcpServer.StopAndClose()

	log.Printf(
		"Listening %s:%s with options -ict=%d -rtt=%d -rta=%d -debug=%t",
		config.Host,
		config.Port,
		config.InactiveConnectionTimeSec,
		config.ReservedTaskStuckTimeSec,
		config.ReservedTaskStuckMaxAttempts,
		config.ProfilerEnabled,
	)

	tcpServer.Listen()
}
