package command

import (
	"strconv"
	"strings"

	"github.com/codecrafters-io/redis-starter-go/app/resp/client"
	"github.com/codecrafters-io/redis-starter-go/app/store"
)

type SetCommand Command

func (cmd *SetCommand) Execute(con client.Client) {
	numArgs := len(cmd.args)
	// TODO write a proper flag parser
	if numArgs != 2 && numArgs != 4 {
		con.SendSimpleError(errWrongNumberOfArgs)
		return
	}
	key := cmd.args[0]
	value := cmd.args[1]
	var expiry int64 = -1

	if numArgs == 4 {
		pxFlag := strings.ToLower(cmd.args[2])
		if pxFlag != "px" {
			con.SendSimpleError(errSyntax)
			return
		}

		exp, err := strconv.ParseInt(cmd.args[3], 10, 64)
		if err != nil || exp < 0 {
			con.SendSimpleError("invalid expiry time")
			return
		}
		expiry = exp
	}

	store.Set(key, value, expiry)
	con.SendSuccess()
}
