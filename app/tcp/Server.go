package tcp

import (
	"bufio"
	"dq/config"
	"dq/node"
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
	Connections map[string]*Connection
	Queue       *node.Queue
	Parser      *Parser
}

func (s *Server) Init() {
	s.Connections = make(map[string]*Connection)
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

		connectionId, err := utils.RandomString(8)
		if err != nil {
			log.Println("Connection id generation failed", err.Error())
			continue
		}

		connection := &Connection{NetConn: conn, Start: time.Now(), Id: connectionId}

		s.mu.Lock()
		s.Connections[connectionId] = connection
		s.mu.Unlock()

		go s.handleConnection(connection)
	}
}

func (s *Server) closeConnection(connectionId string) bool {
	connection, ok := s.Connections[connectionId]
	if !ok {
		return false
	}

	connection.NetConn.Close()
	delete(s.Connections, connectionId)

	return true
}

func (s *Server) handleConnection(c *Connection) {
	defer s.closeConnection(c.Id)

	conn := c.NetConn

	if config.ProfilerEnabled() {
		log.Printf("New connection %s\n", conn.RemoteAddr().String())
	}

	for {
		if config.InactiveConnectionTimeSec() > 0 {
			conn.SetReadDeadline(time.Now().Add(time.Duration(config.InactiveConnectionTimeSec()) * time.Second))
		}

		reader := bufio.NewReader(conn)
		lineBuffer, _, err := reader.ReadLine()
		if err != nil {
			if err == io.EOF {
				if config.ProfilerEnabled() {
					log.Printf("Disconnection %s\n", conn.RemoteAddr().String())
				}
			} else {
				log.Println("Error read buffer", err.Error())
			}

			s.closeConnection(c.Id)
			return
		}

		command, err := s.Parser.ParseCommand(lineBuffer)
		if err != nil {
			conn.Write([]byte(err.Error() + "\n"))
			continue
		}

		if command == commandADD {
			delayMs, err := s.Parser.ParseDelayMs(lineBuffer, 1)
			if err != nil {
				conn.Write([]byte(err.Error() + "\n"))
				continue
			}

			taskBody, err := s.Parser.ParseTaskBody(lineBuffer, 2)
			if err != nil {
				conn.Write([]byte(err.Error() + "\n"))
				continue
			}

			taskId, _ := s.Queue.Add(taskBody, delayMs)
			conn.Write([]byte(fmt.Sprintf("TASK %s DELAY %dms\n", taskId, delayMs)))

			if config.ProfilerEnabled() {
				log.Printf("New task, tasks %d, heap %.2fmb\n", s.Queue.TasksLength(), utils.HeapAllocMb())
			}

			continue
		}

		if command == commandRESERVE {
			task := s.Queue.Reserve()
			if task == nil {
				conn.Write([]byte("nil\n"))
				continue
			}

			conn.Write([]byte(fmt.Sprintf("TASK %s BODY %s\n", task.Id, task.Body)))
			if config.ProfilerEnabled() {
				log.Printf("Task reserved, tasks %d, reserved %d, heap %.2fmb\n", s.Queue.TasksLength(), s.Queue.ReservedTasksLength(), utils.HeapAllocMb())
			}
			continue
		}

		if command == commandDELETE {
			taskId, err := s.Parser.ParseTaskId(lineBuffer, 1)
			if err != nil {
				conn.Write([]byte(err.Error() + "\n"))
				continue
			}

			if !s.Queue.Delete(taskId) {
				conn.Write([]byte("unknown TASK_ID\n"))
				continue
			}

			conn.Write([]byte(fmt.Sprintf("ok\n")))
			if config.ProfilerEnabled() {
				log.Printf("Task deleted, tasks %d, reserved %d, heap %.2fmb\n", s.Queue.TasksLength(), s.Queue.ReservedTasksLength(), utils.HeapAllocMb())
			}
			continue
		}

		if command == commandRETURN {
			taskId, err := s.Parser.ParseTaskId(lineBuffer, 1)
			if err != nil {
				conn.Write([]byte(err.Error() + "\n"))
				continue
			}

			delayMs, err := s.Parser.ParseDelayMs(lineBuffer, 2)
			if err != nil {
				conn.Write([]byte(err.Error() + "\n"))
				continue
			}

			if !s.Queue.Return(taskId, delayMs) {
				conn.Write([]byte("unknown TASK_ID\n"))
				continue
			}

			conn.Write([]byte(fmt.Sprintf("ok\n")))
			if config.ProfilerEnabled() {
				log.Printf("Task returned, tasks %d, reserved %d, heap %.2fmb\n", s.Queue.TasksLength(), s.Queue.ReservedTasksLength(), utils.HeapAllocMb())
			}
			continue
		}

		if command == commandSTATS {
			conn.Write([]byte(fmt.Sprintf("TASKS %d RESERVED %d CONNECTIONS %d HEAP %.2fmb\n", s.Queue.TasksLength(), s.Queue.ReservedTasksLength(), len(s.Connections), utils.HeapAllocMb())))
			continue
		}

		conn.Write([]byte(fmt.Sprintf("unexpected message\n")))
	}
}
