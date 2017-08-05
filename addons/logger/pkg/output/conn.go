package output

import (
	"io"
	"net"
	"sync"
)

type conn struct {
	net.Conn
	sync.Mutex

	addr    string
	network string
}

func (c *conn) dial() error {
	c.Lock()
	defer c.Unlock()

	if c.Conn != nil {
		c.Close()
	}

	var err error
	c.Conn, err = net.Dial(c.network, c.addr)
	return err
}

func (c *conn) Write(data []byte) (int, error) {
	n, err := c.Conn.Write(data)
	if err == io.EOF {
		c.dial()
	}
	return n, err
}

func dial(network, addr string) (net.Conn, error) {
	c := &conn{network: network, addr: addr}
	err := c.dial()
	if err != nil {
		return nil, err
	}
	return c, err
}
