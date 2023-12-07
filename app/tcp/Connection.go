package tcp

import (
	"net"
	"time"
)

type Connection struct {
	Id      string
	NetConn net.Conn
	Start   time.Time
}
