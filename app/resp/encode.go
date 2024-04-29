package resp

import (
	"fmt"
	"strings"
)

func EncodeBulkString(msg string) []byte {
	lenght := len(msg)
	s := fmt.Sprintf("$%d\r\n%s\r\n", lenght, msg)
	return []byte(s)
}

func EncodeSimpleString(msg string) []byte {
	s := fmt.Sprintf("+%s\r\n", msg)
	return []byte(s)
}

func EncodeSimpleError(errMsg string) []byte {
	s := fmt.Sprintf("-%s\r\n", errMsg)
	return []byte(s)
}

func EncodeNullBulkString() []byte {
	return []byte("$-1\r\n")
}

func EncodeArrayBulk(values ...string) []byte {
	respArray := fmt.Sprintf("*%d\r\n", len(values))
	var sb strings.Builder

	sb.WriteString(respArray)
	for _, val := range values {
		sb.WriteString(string(EncodeBulkString(val)))
	}

	return []byte(sb.String())
}

func Success() []byte {
	return EncodeSimpleString("OK")
}
