package command

import (
	"github.com/codecrafters-io/redis-starter-go/app/resp/client"
)

type EchoCommand Command

func (cmd *EchoCommand) Execute(con client.Client) {
	if len(cmd.args) != 1 {
		con.SendSimpleError(errWrongNumberOfArgs)
		return
	}
	con.SendBulkString(cmd.args[0])
}
