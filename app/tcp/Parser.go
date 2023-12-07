package tcp

import (
	"fmt"
	"strconv"
)

// Структура входящего сообщения:
//
// <ATTR_0> <ATTR_1> <ATTR_2> ... <ATTR_N>\n
// где ATTR_0 - команда, ATTR_1, ..., ATTR_N - параметры команды
//
// Команды:
//
// ADD <DELAY_MS> <TASK_BODY>
// RESERVE
// DELETE <TASK_ID>
// RETURN <TASK_ID> <DELAY_MS>
// STATS

type Parser struct {
}

func (p *Parser) ParseCommand(message []byte) (string, error) {
	commandBytes, err := p.parseAttr(message, 0)
	if err != nil {
		return "", err
	}

	commandStr := string(commandBytes)
	commands := allCommands()
	for i := 0; i < len(commands); i++ {
		if commandStr == commands[i] {
			return commandStr, nil
		}
	}

	return "", fmt.Errorf("unknown command")
}

func (p *Parser) ParseDelayMs(message []byte, n int) (uint32, error) {
	msBytes, err := p.parseAttr(message, n)
	if err != nil {
		return 0, fmt.Errorf("invalid DELAY_MS attr")
	}

	msInt64, err := strconv.ParseInt(string(msBytes), 10, 32)
	if err != nil {
		return 0, fmt.Errorf("invalid DELAY_MS attr")
	}

	return uint32(msInt64), nil
}

func (p *Parser) ParseTaskId(message []byte, n int) (string, error) {
	taskId, err := p.parseAttr(message, n)
	if err != nil {
		return "", fmt.Errorf("invalid TASK_ID attr")
	}

	return string(taskId), nil
}

func (p *Parser) ParseTaskBody(message []byte, n int) ([]byte, error) {
	body, err := p.parseAttr(message, n)
	if err != nil {
		return nil, fmt.Errorf("invalid TASK_BODY attr")
	}

	return body, nil
}

func (p *Parser) parseAttr(message []byte, n int) ([]byte, error) {
	attrN := 0
	buf := make([]byte, 0, 10)
	for i := 0; i < len(message); i++ {
		// line break
		if message[i] == 10 {
			break
		}

		// space
		if message[i] == 32 {
			if attrN >= n {
				break
			}
			attrN++

			continue
		}

		if attrN == n {
			buf = append(buf, message[i])
		}
	}

	if attrN == n {
		return buf, nil
	}

	return nil, fmt.Errorf("few attributes")
}
