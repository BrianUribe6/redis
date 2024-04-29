package command

import (
	"github.com/codecrafters-io/redis-starter-go/app/resp"
	"github.com/codecrafters-io/redis-starter-go/app/resp/client"
	"github.com/codecrafters-io/redis-starter-go/app/store"
)

type GetCommand Command

func (cmd *GetCommand) Execute(con client.Client) RESPValue {
	if len(cmd.args) != 1 {
		return resp.EncodeSimpleError(errWrongNumberOfArgs)
	}
	value, exist := store.Get(cmd.args[0])
	if !exist {
		return resp.EncodeNullBulkString()
	}
	return resp.EncodeBulkString(value)
}
