package gofastcgi

import (
	"bufio"
	"io"
)

type streamWriter struct {
	c         *Client
	t         uint8
	requestId uint16
}

func (sw *streamWriter) Write(p []byte) (int, error) {
	nn := 0
	for len(p) > 0 {
		n := len(p)
		if n > maxWrite {
			n = maxWrite
		}
		if err := sw.c.writeRecord(sw.t, sw.requestId, p[:n]); err != nil {
			return nn, err
		}
		nn += n
		p = p[n:]
	}
	return nn, nil
}

func (sw *streamWriter) Close() error {
	sw.c.writeRecord(sw.t, sw.requestId, nil)
	return nil
}

type Buffer struct {
	closer io.Closer
	*bufio.Writer
}

func (b *Buffer) Close() error {
	if err := b.Writer.Flush(); err != nil {
		b.closer.Close()
	}
	return b.closer.Close()
}

func newBuffer(c *Client, t uint8, requestId uint16) *Buffer {
	sw := &streamWriter{c: c, t: t, requestId: requestId}
	bw := bufio.NewWriterSize(sw, maxWrite)
	return &Buffer{sw, bw}
}
