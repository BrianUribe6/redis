package command

import (
	"net"
)

type Executor interface {
	Execute(conn net.Conn)
}

type Command struct {
	label string
	args  []string
}

type PingCommand Command
type EchoCommand Command
type NotImplementedCommand Command

func New(label string, params []string) Executor {
	switch label {
	case "ping":
		return &PingCommand{label: label, args: params}
	case "echo":
		return &EchoCommand{label: label, args: params}
	}
	return &NotImplementedCommand{}
}

func (cmd *PingCommand) Execute(con net.Conn) {
	if len(cmd.args) == 0 {
		ReplySimpleString(con, "PONG")
	} else if len(cmd.args) == 1 {
		ReplyBulkString(con, cmd.args[0])
	} else {
		ReplySimpleError(con, "wrong number of arguments for 'ping' command.")
	}
}

func (cmd *EchoCommand) Execute(con net.Conn) {
	if len(cmd.args) != 1 {
		ReplySimpleError(con, "wrong number of arguments for 'echo' command.")
		return
	}
	ReplyBulkString(con, cmd.args[0])
}

func (cmd *NotImplementedCommand) Execute(con net.Conn) {
	ReplySimpleError(con, "Unknown command, may not be implemented yet")
}
