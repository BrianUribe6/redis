package client

import (
	"bufio"
	"net"
)

type Client struct {
	conn net.Conn
	*bufio.Reader
	*bufio.Writer
	BytesRead int
}

func New(conn net.Conn) Client {
	return Client{
		conn,
		bufio.NewReader(conn),
		bufio.NewWriter(conn),
		0,
	}
}

func (c *Client) Connection() net.Conn {
	return c.conn
}

func (c *Client) Close() error {
	return c.conn.Close()
}
