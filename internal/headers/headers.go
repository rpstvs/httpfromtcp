package headers

import (
	"bytes"
	"fmt"
	"strings"
)

type Headers map[string]string

const crlf = "\r\n"

func NewHeaders() Headers {
	return map[string]string{}
}

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	idx := bytes.Index(data, []byte(crlf))

	if idx == -1 {
		return 0, false, nil
	}

	if idx == 0 {
		return 2, true, nil
	}

	parts := bytes.SplitN(data[:idx], []byte(":"), 2)
	key := strings.ToLower(string(parts[0]))

	if key != strings.TrimRight(key, " ") {
		return 0, false, fmt.Errorf("invalid header name: %s", key)
	}

	value := bytes.TrimSpace(parts[1])
	key = strings.TrimSpace(key)

	if !validTokens([]byte(key)) {
		return 0, false, fmt.Errorf("invalid header token found: %s", key)
	}

	h.Set(key, string(value))
	return idx + 2, false, nil

}

func (h Headers) Set(key, value string) {

	key = strings.ToLower(key)

	if val, ok := h[key]; ok {
		value = strings.Join([]string{val, value}, ",")
	}
	h[key] = value
}

func (h Headers) Override(key, value string) {
	key = strings.ToLower(key)
	h[key] = value
}

func (h Headers) Remove(key string) {
	key = strings.ToLower(key)
	delete(h, key)
}

func (h Headers) Get(key string) (string, bool) {

	key = strings.ToLower(key)

	if _, ok := h[key]; !ok {
		return "", false
	}
	return h[key], true
}

var tokenChars = []byte{'!', '#', '$', '%', '&', '\'', '*', '+', '-', '.', '^', '_', '`', '|', '~'}

func validTokens(data []byte) bool {
	for _, v := range data {
		if !(v >= 'A' && v <= 'Z' || v >= 'a' && v <= 'z' || v >= '0' && v <= '9' || v == '-') {
			return false
		}
	}
	return true
}
