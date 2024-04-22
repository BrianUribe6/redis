package command

import (
	"net"

	"github.com/codecrafters-io/redis-starter-go/app/resp"
)

type ReplConfCommand Command

func (cmd *ReplConfCommand) Execute(con net.Conn) {
	if len(cmd.args) < 2 {
		resp.ReplySimpleError(con, errWrongNumberOfArgs)
		return
	}
	// TODO
	switch cmd.args[0] {
	case "listening-port":
		break
	case "capa":
		break
	default:
		resp.ReplySimpleError(con, errSyntax)
		return
	}

	resp.ReplySuccess(con)
}
