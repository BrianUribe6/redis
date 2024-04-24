package command

import (
	"github.com/codecrafters-io/redis-starter-go/app/resp/client"
)

type ReplConfCommand Command

func (cmd *ReplConfCommand) Execute(con client.Client) {
	if len(cmd.args) < 2 {
		con.SendSimpleError(errWrongNumberOfArgs)
		return
	}
	// TODO
	switch cmd.args[0] {
	case "listening-port":
		break
	case "capa":
		break
	default:
		con.SendSimpleError(errSyntax)
		return
	}

	con.SendSuccess()
}
