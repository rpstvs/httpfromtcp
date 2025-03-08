package main

import (
	"fmt"
	"io"
	"log"
	"os"
)

func main() {
	file, err := os.Open("messages.txt")

	if err != nil {
		log.Fatal("couldnt open file")
	}

	defer file.Close()

	buffer := make([]byte, 8)

	for {
		_, err := file.Read(buffer)
		fmt.Printf("read: %s \n", buffer)
		if err == io.EOF {
			break
		}
	}
}
