package command

import (
	"fmt"
	"net"
)

func ReplyBulkString(conn net.Conn, msg string) {
	lenght := len(msg)
	s := fmt.Sprintf("$%d\r\n%s\r\n", lenght, msg)
	conn.Write([]byte(s))
}

func ReplySimpleString(conn net.Conn, msg string) {
	s := fmt.Sprintf("+%s\r\n", msg)
	conn.Write([]byte(s))
}

func ReplySimpleError(conn net.Conn, errMsg string) {
	s := fmt.Sprintf("-%s\r\n", errMsg)
	conn.Write([]byte(s))
}

func ReplyNullBulkString(conn net.Conn) {
	s := fmt.Sprintf("$-1\r\n")
	conn.Write([]byte(s))
}
