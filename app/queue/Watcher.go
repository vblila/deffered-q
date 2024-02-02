package queue

import (
	"dq/config"
	"log"
	"time"
)

type Watcher struct {
	q *Queue
}

func (w *Watcher) SetQueue(q *Queue) {
	w.q = q
}

func (w *Watcher) WatchFor(taskId string, stuckAttempts uint8) {
	go func() {
		time.Sleep(time.Second * time.Duration(config.ReservedTaskStuckTimeSec))

		if stuckAttempts < uint8(config.ReservedTaskStuckMaxAttempts) {
			returnResult := w.q.Return(taskId, uint32(config.ReservedTaskStuckTimeSec)*1000, true)

			if config.ProfilerEnabled {
				if returnResult {
					log.Printf("Value %s is returned by watcher", taskId)
				} else {
					log.Printf("Watcher can't return task %s", taskId)
				}
			}
		} else {
			deleteResult := w.q.Delete(taskId)

			if config.ProfilerEnabled {
				if deleteResult {
					log.Printf("Value %s is deleted by watcher", taskId)
				} else {
					log.Printf("Watcher can't delete task %s", taskId)
				}
			}
		}
	}()
}
