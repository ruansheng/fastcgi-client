package fastcgi

import (
	"strings"
	"errors"
)

var(
	responseParseError = errors.New("response parse error")
)

type Response struct {
	contentType string
	xPoweredBy string
	content string
}

func (r *Response)Init(data string) (bool, error) {
	rows := strings.Split(data, "\r\n\r\n")
	if len(rows) != 2 {
		return false, responseParseError
	}

	headers := strings.Split(rows[0], "\r\n")

	for _, line := range headers {
		chunks := strings.Split(line, ":")
		if len(chunks) != 2 {
			continue
		}
		if chunks[0] == "Content-type" {
			r.contentType = chunks[1]
		}
		if chunks[0] == "X-Powered-By" {
			r.xPoweredBy = chunks[1]
		}
	}

	r.content = rows[1]
	return true, nil
}

func (r *Response)GetContentType() string {
	return r.contentType
}

func (r *Response)GetXPoweredBy() string {
	return r.xPoweredBy
}

func (r *Response)GetContent() string {
	return r.content
}

