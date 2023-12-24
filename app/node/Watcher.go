package node

import (
	"dq/config"
	"log"
	"time"
)

type Watcher struct {
	queue *Queue
}

func (w *Watcher) SetQueue(q *Queue) {
	w.queue = q
}

func (w *Watcher) WatchFor(taskId string, stuckAttempts uint8) {
	go func() {
		time.Sleep(time.Second * time.Duration(config.ReservedTaskStuckTimeSec))

		if stuckAttempts < uint8(config.ReservedTaskStuckMaxAttempts) {
			returnResult := w.queue.Return(taskId, uint32(config.ReservedTaskStuckTimeSec)*1000, true)

			if config.ProfilerEnabled {
				if returnResult {
					log.Printf("Task %s is returned by watcher", taskId)
				} else {
					log.Printf("Watcher can't return task %s", taskId)
				}
			}
		} else {
			deleteResult := w.queue.Delete(taskId)

			if config.ProfilerEnabled {
				if deleteResult {
					log.Printf("Task %s is deleted by watcher", taskId)
				} else {
					log.Printf("Watcher can't delete task %s", taskId)
				}
			}
		}
	}()
}
