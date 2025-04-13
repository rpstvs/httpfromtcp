package response

import (
	"fmt"
	"io"
	"strconv"

	"github.com/rpstvs/httpfromtcp/internal/headers"
)

type StatusCode int

const (
	StatusOK                  = 200
	StatusBadRequest          = 400
	StatusInternalServerError = 500
)

var StatusLine = map[StatusCode]string{
	200: "HTTP/1.1 200 OK",
	400: "HTTP/1.1 400 Bad Request",
	500: "HTTP/1.1 500 Internal Server Error",
}

func WriteStatusLine(w io.Writer, statusCode StatusCode) error {
	val, ok := StatusLine[statusCode]

	if !ok {
		return fmt.Errorf("status code not found")
	}

	w.Write([]byte(val))
	return nil
}

func GetDefaultHeaders(contentLen int) headers.Headers {
	headersResp := headers.Headers{}
	headersResp.Set("Content-Length", strconv.Itoa(contentLen))
	headersResp.Set("Connection", "close")
	headersResp.Set("Content-Type", "text/plain")
	headersResp.Get("")
	return headersResp
}

func WriteHeaders(w io.Writer, headers headers.Headers) error {
	w.Write([]byte(headers.Get("Cont").val))
	return nil
}
