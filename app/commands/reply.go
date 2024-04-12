package command

import (
	"fmt"
	"net"
	"strings"
)

func ReplyBulkString(conn net.Conn, msg string) {
	conn.Write([]byte(createBulkString(msg)))
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
	conn.Write([]byte("$-1\r\n"))
}

func ReplyArrayBulk(conn net.Conn, arr []string) {
	respArray := fmt.Sprintf("*%d\r\n", len(arr))
	var sb strings.Builder

	sb.WriteString(respArray)
	for _, val := range arr {
		sb.WriteString(createBulkString(val))
	}

	conn.Write([]byte(sb.String()))
}

func ReplySuccess(conn net.Conn) {
	ReplySimpleString(conn, "OK")
}

func createBulkString(msg string) string {
	lenght := len(msg)
	s := fmt.Sprintf("$%d\r\n%s\r\n", lenght, msg)
	return s
}
