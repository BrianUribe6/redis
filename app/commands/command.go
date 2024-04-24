package command

import (
	"github.com/codecrafters-io/redis-starter-go/app/resp"
	"github.com/codecrafters-io/redis-starter-go/app/resp/client"
)

const (
	errWrongNumberOfArgs = "wrong number of arguments"
	errSyntax            = "syntax error"
)

type Executor interface {
	Execute(client client.Client)
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

func (cmd *NotImplementedCommand) Execute(con client.Client) {
	resp.ReplySimpleError(con, "unknown command, may not be implemented yet")
}
