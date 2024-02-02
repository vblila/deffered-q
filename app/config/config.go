package config

var Host string
var Port string
var InactiveConnectionTimeSec uint
var ReservedTaskStuckTimeSec uint
var ReservedTaskStuckMaxAttempts uint
var ProfilerEnabled bool

func SetData(host string, port string, profilerEnabled bool, inactiveConnectionTimeSec uint, reservedTaskStuckTimeSec uint, reservedTaskStuckMaxAttempts uint) {
	Host = host
	Port = port
	ProfilerEnabled = profilerEnabled
	InactiveConnectionTimeSec = inactiveConnectionTimeSec
	ReservedTaskStuckTimeSec = reservedTaskStuckTimeSec
	ReservedTaskStuckMaxAttempts = reservedTaskStuckMaxAttempts
}
