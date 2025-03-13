package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	serverAddr := "localhost:42069"
	udpaddr, err := net.ResolveUDPAddr("udp", serverAddr)

	if err != nil {
		log.Fatal("couldnt resolve ipaddress")
	}

	conn, err := net.DialUDP("udp", nil, udpaddr)

	defer conn.Close()

	reader := bufio.NewReader(os.Stdin)

	for {

		fmt.Println(">")
		message, err := reader.ReadString('\n')

		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
			os.Exit(1)
		}
		conn.Write([]byte(message))
	}
}
