package client

import (
	"bufio"
	"net"
)

type Client struct {
	conn net.Conn
	*bufio.Reader
	*bufio.Writer
}

func New(conn net.Conn) Client {
	return Client{
		conn,
		bufio.NewReader(conn),
		bufio.NewWriter(conn),
	}
}

func (c *Client) Connection() net.Conn {
	return c.conn
}

func (c *Client) Close() error {
	return c.conn.Close()
}
