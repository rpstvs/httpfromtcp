package main

import (
	"fmt"
	"log"
	"net"

	"github.com/rpstvs/httpfromtcp/internal/request"
)

const port = ":42069"

func main() {

	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Println("Couldnt establish a listener ")
		return

	}
	defer listener.Close()
	log.Println("Listening for traffic on", port)
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("couldnt complete connection")
			break
		}
		log.Println("A connection has been accepted")

		req, err := request.RequestFromReader(conn)

		if err != nil {
			log.Println(err)
		}

		fmt.Println("Request line:")
		fmt.Printf("- Method: %s \n", req.RequestLine.Method)
		fmt.Printf("- Target: %s \n", req.RequestLine.RequestTarget)
		fmt.Printf("- Version: %s \n", req.RequestLine.HttpVersion)
		fmt.Println("Headers:")
		for k, v := range req.Headers {
			fmt.Printf("- %s: %s \n", k, v)
		}
		fmt.Println("Connection to ", conn.RemoteAddr(), "closed")

	}
}
