package command

import (
	"net"
	"strconv"
	"strings"

	"github.com/codecrafters-io/redis-starter-go/app/resp"
	"github.com/codecrafters-io/redis-starter-go/app/store"
)

type SetCommand Command

func (cmd *SetCommand) Execute(con net.Conn) {
	numArgs := len(cmd.args)
	// TODO write a proper flag parser
	if numArgs != 2 && numArgs != 4 {
		resp.ReplySimpleError(con, errWrongNumberOfArgs)
		return
	}
	key := cmd.args[0]
	value := cmd.args[1]
	var expiry int64 = -1

	if numArgs == 4 {
		pxFlag := strings.ToLower(cmd.args[2])
		if pxFlag != "px" {
			resp.ReplySimpleError(con, errSyntax)
			return
		}

		exp, err := strconv.ParseInt(cmd.args[3], 10, 64)
		if err != nil || exp < 0 {
			resp.ReplySimpleError(con, "invalid expiry time")
			return
		}
		expiry = exp
	}

	store.Set(key, value, expiry)
	resp.ReplySuccess(con)
}
