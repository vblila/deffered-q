package queue

import "time"

type Task struct {
	Id            string
	Body          []byte
	DelayedTime   time.Time
	StuckAttempts uint8
}
