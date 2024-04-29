package command

import (
	"fmt"

	"github.com/codecrafters-io/redis-starter-go/app/resp"
	"github.com/codecrafters-io/redis-starter-go/app/resp/client"
)

type ReplConfCommand Command

func (cmd *ReplConfCommand) Execute(con client.Client) RESPValue {
	if len(cmd.args) < 2 {
		return resp.EncodeSimpleError(errWrongNumberOfArgs)
	}
	// TODO
	switch cmd.args[0] {
	case "listening-port":
		break
	case "capa":
		break
	case "getack":
		totalBytes := fmt.Sprint(con.BytesRead)
		return resp.EncodeArrayBulk("replconf", "ACK", totalBytes)
	default:
		return resp.EncodeSimpleError(errSyntax)
	}

	return resp.Success()
}
