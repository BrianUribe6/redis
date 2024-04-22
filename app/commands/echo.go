package command

import (
	"net"

	"github.com/codecrafters-io/redis-starter-go/app/resp"
)

type EchoCommand Command

func (cmd *EchoCommand) Execute(con net.Conn) {
	if len(cmd.args) != 1 {
		resp.ReplySimpleError(con, errWrongNumberOfArgs)
		return
	}
	resp.ReplyBulkString(con, cmd.args[0])
}
