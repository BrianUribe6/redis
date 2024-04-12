package command

import (
	"fmt"
	"net"
	"strings"

	"strconv"

	"github.com/codecrafters-io/redis-starter-go/app/store"
)

const (
	errWrongNumberOfArgs = "wrong number of arguments"
	errSyntax            = "syntax error"
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
type InfoCommand Command
type ReplConfCommand Command
type PSYNCCommand Command
type NotImplementedCommand Command

func New(label string, params []string) Executor {
	switch strings.ToLower(label) {
	case "ping":
		return &PingCommand{label, params}
	case "echo":
		return &EchoCommand{label, params}
	case "set":
		return &SetCommand{label, params}
	case "get":
		return &GetCommand{label, params}
	case "info":
		return &InfoCommand{label, params}
	case "replconf":
		return &ReplConfCommand{label, params}
	case "psync":
		return &PSYNCCommand{label, params}
	}
	return &NotImplementedCommand{}
}

func (cmd *PingCommand) Execute(con net.Conn) {
	if len(cmd.args) == 0 {
		ReplySimpleString(con, "PONG")
	} else if len(cmd.args) == 1 {
		ReplyBulkString(con, cmd.args[0])
	} else {
		ReplySimpleError(con, errWrongNumberOfArgs)
	}
}

func (cmd *EchoCommand) Execute(con net.Conn) {
	if len(cmd.args) != 1 {
		ReplySimpleError(con, errWrongNumberOfArgs)
		return
	}
	ReplyBulkString(con, cmd.args[0])
}

func (cmd *NotImplementedCommand) Execute(con net.Conn) {
	ReplySimpleError(con, "unknown command, may not be implemented yet")
}

func (cmd *SetCommand) Execute(con net.Conn) {
	numArgs := len(cmd.args)
	// TODO write a proper flag parser
	if numArgs != 2 && numArgs != 4 {
		ReplySimpleError(con, errWrongNumberOfArgs)
		return
	}
	key := cmd.args[0]
	value := cmd.args[1]
	var expiry int64 = -1

	if numArgs == 4 {
		pxFlag := strings.ToLower(cmd.args[2])
		if pxFlag != "px" {
			ReplySimpleError(con, errSyntax)
			return
		}

		exp, err := strconv.ParseInt(cmd.args[3], 10, 64)
		if err != nil || exp < 0 {
			ReplySimpleError(con, "invalid expiry time")
			return
		}
		expiry = exp
	}

	store.Set(key, value, expiry)
	ReplySuccess(con)
}

func (cmd *GetCommand) Execute(con net.Conn) {
	if len(cmd.args) != 1 {
		ReplySimpleError(con, errWrongNumberOfArgs)
		return
	}
	value, exist := store.Get(cmd.args[0])
	if !exist {
		ReplyNullBulkString(con)
	} else {
		ReplyBulkString(con, value)
	}
}

func (cmd *InfoCommand) Execute(con net.Conn) {
	if len(cmd.args) == 0 {
		ReplySimpleError(con, "unsupported option. Currently only info replication is supported")
		return
	}
	if len(cmd.args) > 1 {
		ReplySimpleError(con, errWrongNumberOfArgs)
		return
	}

	if cmd.args[0] != "replication" {
		ReplySimpleError(con, errSyntax)
		return
	}

	ReplyBulkString(con, store.Info.String())
}

func (cmd *ReplConfCommand) Execute(con net.Conn) {
	if len(cmd.args) < 2 {
		ReplySimpleError(con, errWrongNumberOfArgs)
		return
	}
	// TODO
	switch cmd.args[0] {
	case "listening-port":
		break
	case "capa":
		break
	default:
		ReplySimpleError(con, errSyntax)
		return
	}

	ReplySuccess(con)
}

func (cmd *PSYNCCommand) Execute(con net.Conn) {
	if len(cmd.args) != 3 {
		ReplySimpleError(con, errWrongNumberOfArgs)
		return
	}

	response := fmt.Sprintf("FULLRESYNC %s 0", store.Info.MasterReplId)
	ReplySimpleString(con, response)
}
