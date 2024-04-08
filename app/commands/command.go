package command

import (
	"net"

	"github.com/codecrafters-io/redis-starter-go/app/store"
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
type SetCommand Command
type GetCommand Command
type NotImplementedCommand Command

func New(label string, params []string) Executor {
	switch label {
	case "ping":
		return &PingCommand{label, params}
	case "echo":
		return &EchoCommand{label, params}
	case "set":
		return &SetCommand{label, params}
	case "get":
		return &GetCommand{label, params}
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

func (cmd *SetCommand) Execute(con net.Conn) {
	if len(cmd.args) != 2 {
		ReplySimpleError(con, "wrong number of arguments for 'set' command.")
		return
	}
	store.Set(cmd.args[0], cmd.args[1])
	ReplyBulkString(con, "OK")
}

func (cmd *GetCommand) Execute(con net.Conn) {
	if len(cmd.args) != 1 {
		ReplySimpleError(con, "wrong number of arguments for 'get' command.")
		return
	}
	value, exist := store.Get(cmd.args[0])
	if !exist {
		ReplyNullBulkString(con)
	} else {
		ReplyBulkString(con, value)
	}
}
