package resp

import (
	"fmt"
	"strings"

	"github.com/codecrafters-io/redis-starter-go/app/resp/client"
)

func ReplyBulkString(c client.Client, msg string) {
	c.Write([]byte(createBulkString(msg)))
}

func ReplySimpleString(c client.Client, msg string) {
	s := fmt.Sprintf("+%s\r\n", msg)
	c.Write([]byte(s))
}

func ReplySimpleError(c client.Client, errMsg string) {
	s := fmt.Sprintf("-%s\r\n", errMsg)
	c.Write([]byte(s))
}

func ReplyNullBulkString(c client.Client) {
	c.Write([]byte("$-1\r\n"))
}

func ReplyArrayBulk(c client.Client, values ...string) {
	respArray := fmt.Sprintf("*%d\r\n", len(values))
	var sb strings.Builder

	sb.WriteString(respArray)
	for _, val := range values {
		sb.WriteString(createBulkString(val))
	}

	c.Write([]byte(sb.String()))
}

func ReplySuccess(c client.Client) {
	ReplySimpleString(c, "OK")
}

func createBulkString(msg string) string {
	lenght := len(msg)
	s := fmt.Sprintf("$%d\r\n%s\r\n", lenght, msg)
	return s
}
