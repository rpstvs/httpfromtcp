package response

import (
	"fmt"
	"io"

	"github.com/rpstvs/httpfromtcp/internal/headers"
)

func GetDefaultHeaders(contentLen int) headers.Headers {
	headersResp := headers.NewHeaders()
	headersResp.Set("Content-Length", fmt.Sprintf("%d", contentLen))
	headersResp.Set("Connection", "close")
	headersResp.Set("Content-Type", "text/plain")

	return headersResp
}

func WriteHeaders(w io.Writer, headers headers.Headers) error {
	for k, v := range headers {
		_, err := w.Write([]byte(fmt.Sprintf("%s: %s\r\n", k, v)))
		if err != nil {
			return err
		}
	}
	_, err := w.Write([]byte("\r\n"))
	return err
}
