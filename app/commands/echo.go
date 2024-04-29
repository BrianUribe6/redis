package command

import (
	"github.com/codecrafters-io/redis-starter-go/app/resp"
	"github.com/codecrafters-io/redis-starter-go/app/resp/client"
)

type EchoCommand Command

func (cmd *EchoCommand) Execute(con client.Client) RESPValue {
	if len(cmd.args) != 1 {
		return resp.EncodeSimpleError(errWrongNumberOfArgs)
	}
	return resp.EncodeBulkString(cmd.args[0])
}
