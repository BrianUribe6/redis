package command

import (
	"github.com/codecrafters-io/redis-starter-go/app/resp"
	"github.com/codecrafters-io/redis-starter-go/app/resp/client"
)

type PingCommand Command

func (cmd *PingCommand) Execute(con client.Client) RESPValue {
	if len(cmd.args) == 0 {
		return resp.EncodeSimpleString("PONG")
	}
	if len(cmd.args) == 1 {
		return resp.EncodeBulkString(cmd.args[0])
	}
	return resp.EncodeSimpleError(errWrongNumberOfArgs)
}
