package client

import (
	"fmt"
	"strings"
)

func (c *Client) SendBulkString(msg string) {
	c.Write([]byte(createBulkString(msg)))
}

func (c *Client) SendSimpleString(msg string) {
	s := fmt.Sprintf("+%s\r\n", msg)
	c.Write([]byte(s))
}

func (c *Client) SendSimpleError(errMsg string) {
	s := fmt.Sprintf("-%s\r\n", errMsg)
	c.Write([]byte(s))
}

func (c *Client) SendNullBulkString() {
	c.Write([]byte("$-1\r\n"))
}

func (c *Client) SendArrayBulk(values ...string) {
	respArray := fmt.Sprintf("*%d\r\n", len(values))
	var sb strings.Builder

	sb.WriteString(respArray)
	for _, val := range values {
		sb.WriteString(createBulkString(val))
	}

	c.Write([]byte(sb.String()))
}

func (c *Client) SendSuccess() {
	c.SendSimpleString("OK")
}

func createBulkString(msg string) string {
	lenght := len(msg)
	s := fmt.Sprintf("$%d\r\n%s\r\n", lenght, msg)
	return s
}
