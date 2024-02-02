package tcp

import (
	"bufio"
	"dq/config"
	"dq/queue"
	"dq/utils"
	"fmt"
	"io"
	"log"
	"net"
	"sync"
	"time"
)

type Server struct {
	mu          sync.Mutex
	listener    net.Listener
	connections map[string]*Connection
	q           *queue.Queue
	parser      Parser
	watcher     queue.Watcher
}

func (s *Server) Init() {
	s.connections = make(map[string]*Connection, 128)
	s.parser = Parser{}

	q := &queue.Queue{}
	q.Init()

	s.q = q

	s.watcher = queue.Watcher{}
	s.watcher.SetQueue(q)
}

func (s *Server) Start(host string, port string) error {
	var err error
	s.listener, err = net.Listen("tcp", host+":"+port)
	return err
}

func (s *Server) StopAndClose() {
	s.listener.Close()
}

func (s *Server) Listen() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			log.Println("Error connect accept:", err.Error())
			continue
		}

		connection := &Connection{NetConn: conn, Start: time.Now(), Id: utils.RandomId()}
		go s.handleConnection(connection)
	}
}

func (s *Server) addConnection(connection *Connection) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.connections[connection.Id] = connection
}

func (s *Server) closeConnection(connectionId string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	connection, ok := s.connections[connectionId]
	if !ok {
		return false
	}

	connection.NetConn.Close()
	delete(s.connections, connectionId)

	return true
}

func (s *Server) getConnectionsCount() int {
	s.mu.Lock()
	defer s.mu.Unlock()

	return len(s.connections)
}

func (s *Server) handleConnection(c *Connection) {
	s.addConnection(c)
	defer s.closeConnection(c.Id)

	conn := c.NetConn

	if config.ProfilerEnabled {
		log.Printf("New connection %s\n", conn.RemoteAddr().String())
	}

	for {
		if config.InactiveConnectionTimeSec > 0 {
			conn.SetReadDeadline(time.Now().Add(time.Duration(config.InactiveConnectionTimeSec) * time.Second))
		}

		reader := bufio.NewReader(conn)
		lineBuffer, _, err := reader.ReadLine()
		if err != nil {
			if err == io.EOF {
				if config.ProfilerEnabled {
					log.Printf("Disconnection %s\n", conn.RemoteAddr().String())
				}
			} else {
				log.Println("Error read buffer", err.Error())
			}

			s.closeConnection(c.Id)
			return
		}

		command, err := s.parser.ParseCommand(lineBuffer)
		if err != nil {
			conn.Write([]byte(err.Error() + "\n"))
			continue
		}

		if command == commandADD {
			delayMs, err := s.parser.ParseDelayMs(lineBuffer, 1)
			if err != nil {
				conn.Write([]byte(err.Error() + "\n"))
				continue
			}

			taskBody, err := s.parser.ParseTaskBody(lineBuffer, 2)
			if err != nil {
				conn.Write([]byte(err.Error() + "\n"))
				continue
			}

			taskId, _ := s.q.Add(taskBody, delayMs)
			conn.Write([]byte(fmt.Sprintf("TASK %s DELAY %dms\n", taskId, delayMs)))

			if config.ProfilerEnabled {
				log.Printf("New task, tasks %d, heap %.2fmb\n", s.q.TasksLength(), utils.HeapAllocMb())
			}

			continue
		}

		if command == commandRESERVE {
			taskId, taskBody, stuckAttempts, ok := s.q.Reserve()
			if ok == false {
				conn.Write([]byte("nil\n"))
				continue
			}

			if config.ReservedTaskStuckTimeSec > 0 {
				s.watcher.WatchFor(taskId, stuckAttempts)
			}

			conn.Write([]byte(fmt.Sprintf("TASK %s BODY %s\n", taskId, taskBody)))
			if config.ProfilerEnabled {
				log.Printf("Value reserved, tasks %d, reserved %d, heap %.2fmb\n", s.q.TasksLength(), s.q.ReservedTasksLength(), utils.HeapAllocMb())
			}
			continue
		}

		if command == commandDELETE {
			taskId, err := s.parser.ParseTaskId(lineBuffer, 1)
			if err != nil {
				conn.Write([]byte(err.Error() + "\n"))
				continue
			}

			if !s.q.Delete(taskId) {
				conn.Write([]byte("unknown TASK_ID\n"))
				continue
			}

			conn.Write([]byte(fmt.Sprintf("ok\n")))
			if config.ProfilerEnabled {
				log.Printf("Value deleted, tasks %d, reserved %d, heap %.2fmb\n", s.q.TasksLength(), s.q.ReservedTasksLength(), utils.HeapAllocMb())
			}
			continue
		}

		if command == commandRETURN {
			taskId, err := s.parser.ParseTaskId(lineBuffer, 1)
			if err != nil {
				conn.Write([]byte(err.Error() + "\n"))
				continue
			}

			delayMs, err := s.parser.ParseDelayMs(lineBuffer, 2)
			if err != nil {
				conn.Write([]byte(err.Error() + "\n"))
				continue
			}

			if !s.q.Return(taskId, delayMs, false) {
				conn.Write([]byte("unknown TASK_ID\n"))
				continue
			}

			conn.Write([]byte(fmt.Sprintf("ok\n")))
			if config.ProfilerEnabled {
				log.Printf("Value returned, tasks %d, reserved %d, heap %.2fmb\n", s.q.TasksLength(), s.q.ReservedTasksLength(), utils.HeapAllocMb())
			}
			continue
		}

		if command == commandSTATS {
			conn.Write([]byte(fmt.Sprintf("TASKS %d RESERVED %d CONNECTIONS %d HEAP %.2fm\n", s.q.TasksLength(), s.q.ReservedTasksLength(), s.getConnectionsCount(), utils.HeapAllocMb())))
			continue
		}

		conn.Write([]byte(fmt.Sprintf("unexpected message\n")))
	}
}
