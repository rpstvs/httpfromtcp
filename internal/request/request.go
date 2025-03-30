package request

import (
	"bytes"
	"fmt"
	"io"
	"strings"
)

type Request struct {
	RequestLine RequestLine
	state       requestState
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

type requestState int

const crlf = "\r\n"

const (
	initialized requestState = iota
	done
)

const bufferSize = 8

func RequestFromReader(reader io.Reader) (*Request, error) {

	buf := make([]byte, bufferSize)
	readToIndex = 0

	req := Request{
		state: initialized,
	}

	for req.state != done {
		if len(buf) == bufferSize {
			buf2 := make([]byte, len(buf)*2)
			copy(buf2, buf)
		}
	}

}

func parseRequestLine(data []byte) (*RequestLine, int, error) {
	idx := bytes.Index(data, []byte(crlf))

	if idx == -1 {
		return nil, 0, fmt.Errorf("could not find CRLF in request-line")
	}

	requestLineText := string(data[:idx])
	requestLine, err := requestLineFromString(requestLineText)

	if err != nil {
		return nil, 0, err
	}
	return requestLine, idx + 2, nil
}

func requestLineFromString(str string) (*RequestLine, error) {
	reqLine := strings.Split(str, " ")

	if len(reqLine) != 3 {
		return nil, fmt.Errorf("bad request line: %s", str)
	}

	method := reqLine[0]

	for _, c := range method {
		if c < 'A' || c > 'Z' {
			return nil, fmt.Errorf("invalid method: %s", method)
		}
	}

	requestTarget := reqLine[1]

	versionParts := strings.Split(reqLine[2], "/")

	if len(versionParts) != 2 {
		return nil, fmt.Errorf("malformed start-line: %s", str)
	}

	httpPart := versionParts[0]

	if httpPart != "HTTP" {
		return nil, fmt.Errorf("unrecognized HTTP-version: %s", httpPart)
	}

	version := versionParts[1]

	if version != "1.1" {
		return nil, fmt.Errorf("unrecognized HTTP-version: %s", version)
	}

	return &RequestLine{
		Method:        method,
		RequestTarget: requestTarget,
		HttpVersion:   version,
	}, nil

}

func (r *Request) parse(data []byte) (int, error) {

	switch r.state {
	case initialized:
		requestLine, n, err := parseRequestLine(data)
		if n == 0 {
			return 0, nil
		}
		r.RequestLine = *requestLine
		r.state = done
	case done:
		return 0, fmt.Errorf("error: trying to read data in a done state")
	default:
		return 0, fmt.Errorf("error: unknown state")
	}
	return int(r.state), nil
}
