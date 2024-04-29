package command

import (
	"strconv"
	"strings"

	"github.com/codecrafters-io/redis-starter-go/app/resp"
	"github.com/codecrafters-io/redis-starter-go/app/resp/client"
	"github.com/codecrafters-io/redis-starter-go/app/store"
)

type SetCommand Command

func (cmd *SetCommand) Execute(con client.Client) RESPValue {
	numArgs := len(cmd.args)
	// TODO write a proper flag parser
	if numArgs != 2 && numArgs != 4 {
		return resp.EncodeSimpleError(errWrongNumberOfArgs)
	}
	key := cmd.args[0]
	value := cmd.args[1]
	var expiry int64 = -1

	if numArgs == 4 {
		pxFlag := strings.ToLower(cmd.args[2])
		if pxFlag != "px" {
			return resp.EncodeSimpleError(errSyntax)
		}

		exp, err := strconv.ParseInt(cmd.args[3], 10, 64)
		if err != nil || exp < 0 {
			return resp.EncodeSimpleError("invalid expiry time")
		}
		expiry = exp
	}

	store.Set(key, value, expiry)
	return resp.Success()
}
