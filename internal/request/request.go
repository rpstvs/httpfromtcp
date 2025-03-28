package request

import (
	"fmt"
	"io"
	"strings"
	"unicode"
)

type Request struct {
	RequestLine RequestLine
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func RequestFromReader(reader io.Reader) (*Request, error) {

	dat, err := io.ReadAll(reader)

	if err != nil {
		fmt.Println(err)
		return &Request{}, nil
	}

	reqParts := strings.Split(string(dat), "\r\n")

	reqLine := strings.Split(reqParts[0], " ")

	if len(reqLine) != 3 {
		return &Request{}, fmt.Errorf("bad request line")
	}

	if !isUpper(reqLine[0]) {
		return &Request{}, fmt.Errorf()
	}

	return &Request{}, nil
}

func isUpper(s string) bool {
	for _, r := range s {
		if !unicode.IsUpper(r) && unicode.IsLetter(r) {
			return false
		}
	}
	return true
}
