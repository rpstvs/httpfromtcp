package server

import (
	"fmt"
	"log"
	"net"
	"sync/atomic"
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

	response := "HTTP/1.1 200 OK\r\n" +
		"Content-Type: text/plain\r\n" +
		"\r\n" +
		"Hello World!"

	_, err := conn.Write([]byte(response))
	if err != nil {
		log.Println("couldnt write response")
	}
}
