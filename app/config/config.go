package config

import (
	"log"
	"os"
	"strconv"
)

var ProfilerEnabled bool
var InactiveConnectionTimeSec uint
var ReservedTaskStuckTimeSec uint
var ReservedTaskStuckMaxAttempts uint

func init() {
	ProfilerEnabled = readBoolEnv("DQ_PROFILER_ENABLED")
	InactiveConnectionTimeSec = readUintEnv("DQ_INACTIVE_CONNECTION_TIME_SECONDS")
	ReservedTaskStuckTimeSec = readUintEnv("DQ_RESERVED_TASK_STUCK_TIME_SECONDS")
	ReservedTaskStuckMaxAttempts = readUintEnv("DQ_RESERVED_TASK_STUCK_MAX_ATTEMPTS")
}

func readBoolEnv(envKey string) bool {
	value := os.Getenv(envKey)
	if value == "" {
		log.Fatalf("Env %s is required\n", envKey)
	}

	if value == "1" {
		return true
	}

	return false
}

func readUintEnv(envKey string) uint {
	value, err := strconv.Atoi(os.Getenv(envKey))
	if err != nil {
		log.Fatalf("Env %s integer value is required\n", envKey)
	}

	if value < 0 {
		log.Fatalf("Env %s positive value is required\n", envKey)
	}

	return uint(value)
}
