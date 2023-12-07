package node

import "time"

type Task struct {
	Id          string
	Body        []byte
	DelayedTime time.Time
}
