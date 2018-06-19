package fastcgi

import (
	"sync"
	"io"
	"bytes"
	"net"
	"strconv"
	"errors"
	"encoding/binary"
)

var(
	hostEmptyError = errors.New("host is empty")
	portEmptyError = errors.New("port is error")
)

var pad [maxPad]byte

type Client struct {
	lock      sync.Mutex
	conn      io.ReadWriteCloser
	h         header
	buf       bytes.Buffer
	keepAlive bool
}

func New(host string, port int) (*Client, error) {
	if host == "" {
		return nil, hostEmptyError
	}
	if port <= 80 {
		return nil, portEmptyError
	}

	var conn net.Conn
	addr := host + ":" + strconv.FormatInt(int64(port), 10)
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}

	client := &Client{
		conn:      conn,
		keepAlive: false,
	}
	return client, nil
}

func (c *Client)writeBeginRequest(requestId uint16, t uint8, role uint16, flags uint8) error {
	b := [8]byte{0, byte(role), flags}
	return c.writeRecord(t, requestId, b[:])
}

func (c *Client)writePairs(t uint8, requestId uint16, pairs map[string]string) error {
	buffer := newBuffer(c, t, requestId)
	b := make([]byte, 8)
	for k, v := range pairs {
		n := encodeSize(b, uint32(len(k)))
		n += encodeSize(b[n:], uint32(len(v)))
		if _,err := buffer.Write(b[:n]); err != nil {
			return err
		}
		if _,err := buffer.WriteString(k); err != nil {
			return err
		}
		if _,err := buffer.WriteString(v); err != nil {
			return err
		}
	}
	buffer.Close()
	return nil
}

func (c *Client)writeRecord(t uint8, requestId uint16, content []byte) (error){
	c.lock.Lock()
	defer c.lock.Unlock()

	c.buf.Reset()

	contentLength := len(content)
	c.h.init(t, requestId, uint16(contentLength))

	// add header
	if err := binary.Write(&c.buf, binary.BigEndian, c.h); err != nil {
		return err
	}

	// add content
	if _, err := c.buf.Write(content); err != nil {
		return err
	}

	// add padding content
	if _, err := c.buf.Write(pad[:c.h.PaddingLength]); err != nil {
		return err
	}

	_, err := c.conn.Write(c.buf.Bytes())
	return err
}

func (c *Client) Request(env map[string]string, reqParams string) (response *Response, err error){
	var requestId uint16 = 1
	var role uint16 = 1
	var flags uint8 = 0

	err = c.writeBeginRequest(requestId, FcgiBeginRequest, role, flags)
	if err != nil {
		return
	}

	err = c.writePairs(FcgiParams, requestId, env)
	if err != nil {
		return
	}

	// write post request data
	if len(reqParams) > 0 {
		err = c.writeRecord(FcgiStdin, requestId, []byte(reqParams))
		if err != nil {
			return
		}
	}

	rec := &record{}
	err = rec.read(c.conn)
	if err != nil {
		return
	}

	ret := rec.content()

	// parse php-fmt response ret
	response, err = c.parseContent(string(ret))
	return response, err
}

func (c *Client)parseContent(ret string) (*Response, error) {
	response := &Response{}
	_, err := response.init(string(ret))
	if err != nil {
		return nil, err
	}
	return response, nil
}