package config

import (
	"log"
	"os"
	"strconv"
)

var profilerEnabled string
var inactiveConnectionTimeSec int
var reservedTaskStuckTimeSec int
var reservedTaskStuckMaxAttempts int

func init() {
	inactiveConnectionTimeSec = -1
	reservedTaskStuckTimeSec = -1
	reservedTaskStuckMaxAttempts = -1
}

func ProfilerEnabled() bool {
	if profilerEnabled == "" {
		profilerEnabled = os.Getenv("DQ_PROFILER_ENABLED")
		if profilerEnabled == "" {
			log.Fatalln("Env DQ_PROFILER_ENABLED is required")
		}
	}

	if profilerEnabled == "1" {
		return true
	}

	return false
}

func InactiveConnectionTimeSec() int {
	var err error

	if inactiveConnectionTimeSec == -1 {
		inactiveConnectionTimeSec, err = strconv.Atoi(os.Getenv("DQ_INACTIVE_CONNECTION_TIME_SECONDS"))
		if err != nil {
			log.Fatalln("Env DQ_INACTIVE_CONNECTION_TIME_SECONDS integer value is required")
		}

		if inactiveConnectionTimeSec < 0 {
			log.Fatalln("Env DQ_INACTIVE_CONNECTION_TIME_SECONDS positive value is required")
		}
	}

	return inactiveConnectionTimeSec
}

func ReservedTaskStuckTimeSec() int {
	var err error

	if reservedTaskStuckTimeSec == -1 {
		reservedTaskStuckTimeSec, err = strconv.Atoi(os.Getenv("DQ_RESERVED_TASK_STUCK_TIME_SECONDS"))
		if err != nil {
			log.Fatalln("Env DQ_RESERVED_TASK_STUCK_TIME_SECONDS integer value is required")
		}

		if reservedTaskStuckTimeSec < 0 {
			log.Fatalln("Env DQ_RESERVED_TASK_STUCK_TIME_SECONDS positive value is required")
		}
	}

	return reservedTaskStuckTimeSec
}

func ReservedTaskStuckMaxAttempts() int {
	var err error

	if reservedTaskStuckMaxAttempts == -1 {
		reservedTaskStuckMaxAttempts, err = strconv.Atoi(os.Getenv("DQ_RESERVED_TASK_STUCK_MAX_ATTEMPTS"))
		if err != nil {
			log.Fatalln("Env DQ_RESERVED_TASK_STUCK_MAX_ATTEMPTS integer value is required")
		}

		if reservedTaskStuckMaxAttempts < 0 {
			log.Fatalln("Env DQ_RESERVED_TASK_STUCK_MAX_ATTEMPTS positive value is required")
		}
	}

	return reservedTaskStuckMaxAttempts
}
