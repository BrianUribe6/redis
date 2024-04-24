package command

import (
	"github.com/codecrafters-io/redis-starter-go/app/resp"
	"github.com/codecrafters-io/redis-starter-go/app/resp/client"
)

type PingCommand Command

func (cmd *PingCommand) Execute(con client.Client) {
	if len(cmd.args) == 0 {
		resp.ReplySimpleString(con, "PONG")
	} else if len(cmd.args) == 1 {
		resp.ReplyBulkString(con, cmd.args[0])
	} else {
		resp.ReplySimpleError(con, errWrongNumberOfArgs)
	}
}
