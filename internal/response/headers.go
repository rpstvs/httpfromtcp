package response

import (
	"fmt"

	"github.com/rpstvs/httpfromtcp/internal/headers"
)

func GetDefaultHeaders(contentLen int) headers.Headers {
	headersResp := headers.NewHeaders()
	headersResp.Set("Content-Length", fmt.Sprintf("%d", contentLen))
	headersResp.Set("Connection", "close")
	headersResp.Set("Content-Type", "text/plain")

	return headersResp
}
