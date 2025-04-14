package server

import (
	"fmt"
	"log"
	"net"
	"sync/atomic"

	"github.com/rpstvs/httpfromtcp/internal/response"
)

type Server struct {
	listener net.Listener
	closed   atomic.Bool
}

func Serve(port int) (*Server, error) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, err
	}

	server := Server{
		listener: listener,
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
	return server.listener.Close()
}

func (server *Server) handle(conn net.Conn) {
	defer conn.Close()

	response.WriteStatusLine(conn, response.StatusOK)
	headers := response.GetDefaultHeaders(0)
	if err := response.WriteHeaders(conn, headers); err != nil {
		fmt.Printf("error: %v", err)
	}
}
