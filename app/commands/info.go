package command

import (
	"github.com/codecrafters-io/redis-starter-go/app/resp/client"
	"github.com/codecrafters-io/redis-starter-go/app/store"
)

type InfoCommand Command

func (cmd *InfoCommand) Execute(con client.Client) {
	if len(cmd.args) == 0 {
		con.SendSimpleError("unsupported option. Currently only info replication is supported")
		return
	}
	if len(cmd.args) > 1 {
		con.SendSimpleError(errWrongNumberOfArgs)
		return
	}

	if cmd.args[0] != "replication" {
		con.SendSimpleError(errSyntax)
		return
	}

	con.SendBulkString(store.Info.String())
}
