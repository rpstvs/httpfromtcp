package server

import (
	"io"

	"github.com/rpstvs/httpfromtcp/internal/request"
)

type HandlerError struct {
	statusCode int
	message    string
}

type Handler func(w io.Writer, req *request.Request) *HandlerError
