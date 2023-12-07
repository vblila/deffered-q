package config

import (
	"log"
	"os"
	"strconv"
)

var profilerEnabled string
var inactiveConnectionTimeSec int

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

	if inactiveConnectionTimeSec == 0 {
		inactiveConnectionTimeSec, err = strconv.Atoi(os.Getenv("DQ_INACTIVE_CONNECTION_TIME_SECONDS"))
		if err != nil {
			log.Fatalln("Env DQ_INACTIVE_CONNECTION_TIME_SECONDS integer value is required")
		}

		if inactiveConnectionTimeSec < 0 {
			log.Fatalln("Env PTM_TASK_EXECUTION_INTERVAL_MS positive value is required")
		}

		if inactiveConnectionTimeSec == 0 {
			inactiveConnectionTimeSec = -1
		}
	}

	return inactiveConnectionTimeSec
}
