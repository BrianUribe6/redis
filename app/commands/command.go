package command

import (
	"net"

	"github.com/codecrafters-io/redis-starter-go/app/resp"
)

const (
	errWrongNumberOfArgs = "wrong number of arguments"
	errSyntax            = "syntax error"
)

type Executor interface {
	Execute(conn net.Conn)
}

type Command struct {
	label      string
	args       []string
	IsMutation bool
}

type NotImplementedCommand Command

func New(label string, params []string) Executor {
	switch label {
	case "ping":
		return &PingCommand{label: label, args: params}
	case "echo":
		return &EchoCommand{label: label, args: params}
	case "set":
		return &SetCommand{label: label, args: params, IsMutation: true}
	case "get":
		return &GetCommand{label: label, args: params}
	case "info":
		return &InfoCommand{label: label, args: params}
	case "replconf":
		return &ReplConfCommand{label: label, args: params}
	case "psync":
		return &PSYNCCommand{label: label, args: params}
	}
	return &NotImplementedCommand{}
}

func (cmd *NotImplementedCommand) Execute(con net.Conn) {
	resp.ReplySimpleError(con, "unknown command, may not be implemented yet")
}
