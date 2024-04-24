package command

import (
	"github.com/codecrafters-io/redis-starter-go/app/resp/client"
)

type PingCommand Command

func (cmd *PingCommand) Execute(con client.Client) {
	if len(cmd.args) == 0 {
		con.SendSimpleString("PONG")
	} else if len(cmd.args) == 1 {
		con.SendBulkString(cmd.args[0])
	} else {
		con.SendSimpleError(errWrongNumberOfArgs)
	}
}
