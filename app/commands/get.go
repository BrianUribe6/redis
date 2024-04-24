package command

import (
	"github.com/codecrafters-io/redis-starter-go/app/resp/client"
	"github.com/codecrafters-io/redis-starter-go/app/store"
)

type GetCommand Command

func (cmd *GetCommand) Execute(con client.Client) {
	if len(cmd.args) != 1 {
		con.SendSimpleError(errWrongNumberOfArgs)
		return
	}
	value, exist := store.Get(cmd.args[0])
	if !exist {
		con.SendNullBulkString()
	} else {
		con.SendBulkString(value)
	}
}
