package command

import (
	"net"
	"strings"

	"strconv"

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
	switch strings.ToLower(label) {
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
	ReplySimpleError(con, "unknown command, may not be implemented yet")
}

func (cmd *SetCommand) Execute(con net.Conn) {
	numArgs := len(cmd.args)
	if numArgs < 2 || numArgs > 4 {
		ReplySimpleError(con, "wrong number of arguments for 'set' command.")
		return
	}
	key := cmd.args[0]
	value := cmd.args[1]
	var expiry int64 = -1

	if numArgs > 3 && strings.ToLower(cmd.args[2]) == "px" {
		ms, err := strconv.ParseInt(cmd.args[3], 10, 64)
		if err != nil || ms < 0 {
			ReplySimpleError(con, "invalid expiry time")
			return
		}
		expiry = ms
	} else {
		ReplySimpleError(con, "syntax error")
		return
	}

	store.Set(key, value, expiry)
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
