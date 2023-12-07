package tcp

const (
	commandADD     = "ADD"
	commandRESERVE = "RESERVE"
	commandDELETE  = "DELETE"
	commandRETURN  = "RETURN"
	commandSTATS   = "STATS"
)

func allCommands() [5]string {
	return [...]string{commandADD, commandRESERVE, commandDELETE, commandRETURN, commandSTATS}
}
