package command

import (
	"net"

	"github.com/codecrafters-io/redis-starter-go/app/resp"
	"github.com/codecrafters-io/redis-starter-go/app/store"
)

type InfoCommand Command

func (cmd *InfoCommand) Execute(con net.Conn) {
	if len(cmd.args) == 0 {
		resp.ReplySimpleError(con, "unsupported option. Currently only info replication is supported")
		return
	}
	if len(cmd.args) > 1 {
		resp.ReplySimpleError(con, errWrongNumberOfArgs)
		return
	}

	if cmd.args[0] != "replication" {
		resp.ReplySimpleError(con, errSyntax)
		return
	}

	resp.ReplyBulkString(con, store.Info.String())
}
