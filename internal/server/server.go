package server

import (
	"fmt"
	"log"
	"net"
	"sync/atomic"

	"github.com/rpstvs/httpfromtcp/internal/request"
	"github.com/rpstvs/httpfromtcp/internal/response"
)

type Handler func(w *response.Writer, req *request.Request)
type Server struct {
	listener net.Listener
	closed   atomic.Bool
	handler  Handler
}

func Serve(port int, handler Handler) (*Server, error) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, err
	}

	server := Server{
		listener: listener,
		handler:  handler,
	}

	go server.listen()

	return &server, nil
}

func (server *Server) listen() {
	for {
		conn, err := server.listener.Accept()
		if err != nil {
			if server.closed.Load() {
				return
			}
			log.Println("couldnt establish conn")
			continue

		}
		go server.handle(conn)

	}

}

func (server *Server) Close() error {
	server.closed.Store(true)
	if server.listener != nil {
		return server.listener.Close()
	}
	return nil
}

func (server *Server) handle(conn net.Conn) {
	defer conn.Close()

	w := response.NewWriter(conn)

	req, err := request.RequestFromReader(conn)

	if err != nil {
		w.WriteStatusLine(response.StatusBadRequest)
		body := []byte(fmt.Sprintf("Error parsing request: %v", err))
		w.WriteHeaders(response.GetDefaultHeaders(len(body)))
		w.WriteBody(body)
		return
	}
	server.handler(w, req)

}
