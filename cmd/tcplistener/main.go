package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
)

const port = ":42069"

func main() {

	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Println("Couldnt establish a listener ")
		return

	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("couldnt complete connection")
			break
		}
		log.Println("A connection has been accepted")

		lines := getLinesChannel(conn)

		for line := range lines {
			fmt.Println(line)
		}

		fmt.Println("Connection to ", conn.RemoteAddr(), "closed")

	}
}

func getLinesChannel(f io.ReadCloser) <-chan string {
	currentLine := ""
	c := make(chan string)
	go func() {
		defer f.Close()
		defer close(c)
		for {
			buffer := make([]byte, 8, 8)
			n, err := f.Read(buffer)

			if err != nil {
				if currentLine != "" {
					c <- currentLine
					currentLine = ""
				}
				if errors.Is(err, io.EOF) {
					break
				}
				fmt.Printf("error: %s\n", err.Error())
				break
			}

			str := string(buffer[:n])

			parts := strings.Split(str, "\n")

			for i := 0; i < len(parts)-1; i++ {
				c <- fmt.Sprintf("%s%s", currentLine, parts[i])
				currentLine = ""
			}
			currentLine += parts[len(parts)-1]

		}
	}()

	return c
}
