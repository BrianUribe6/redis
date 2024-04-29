package command

import (
	"github.com/codecrafters-io/redis-starter-go/app/resp"
	"github.com/codecrafters-io/redis-starter-go/app/resp/client"
	"github.com/codecrafters-io/redis-starter-go/app/store"
)

type InfoCommand Command

func (cmd *InfoCommand) Execute(con client.Client) RESPValue {
	if len(cmd.args) == 0 {
		return resp.EncodeSimpleError("unsupported option. Currently only info replication is supported")
	}
	if len(cmd.args) > 1 {
		return resp.EncodeSimpleError(errWrongNumberOfArgs)
	}

	if cmd.args[0] != "replication" {
		return resp.EncodeSimpleError(errSyntax)
	}

	return resp.EncodeBulkString(store.Info.String())
}
