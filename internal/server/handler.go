package server

import (
	"io"

	"github.com/rpstvs/httpfromtcp/internal/request"
	"github.com/rpstvs/httpfromtcp/internal/response"
)

type HandlerError struct {
	StatusCode response.StatusCode
	Message    string
}

type Handler func(w io.Writer, req *request.Request) *HandlerError
