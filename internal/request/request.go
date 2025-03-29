package request

import (
	"bytes"
	"fmt"
	"io"
	"strings"
)

type Request struct {
	RequestLine RequestLine
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

const crlf = "\r\n"

func RequestFromReader(reader io.Reader) (*Request, error) {

	dat, err := io.ReadAll(reader)

	if err != nil {
		return nil, err
	}

	requestLine, err := parseRequestLine(dat)

	if err != nil {
		return nil, err
	}

	return &Request{
		RequestLine: *requestLine,
	}, nil

}

func parseRequestLine(data []byte) (*RequestLine, error) {
	idx := bytes.Index(data, []byte(crlf))

	if idx == -1 {
		return nil, fmt.Errorf("could not find CRLF in request-line")
	}

	requestLineText := string(data[:idx])
	requestLine, err := requestLineFromString(requestLineText)

	if err != nil {
		return nil, err
	}
	return requestLine, nil
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
