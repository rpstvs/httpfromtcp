package request

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRequestLineParser(t *testing.T) {
	r, err := RequestFromReader(strings.Reader("GET / HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n"))
	require.NoError(t, err)
}
