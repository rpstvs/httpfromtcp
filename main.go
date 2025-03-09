package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

const inputFilePath = "messages.txt"

func main() {
	file, err := os.Open(inputFilePath)

	if err != nil {
		log.Fatalf("couldnt open file %s: %s\n", inputFilePath, err)
	}

	fmt.Printf("Reading data from %s\n", inputFilePath)
	fmt.Println("=====================================")
	currentLine := ""
	for {
		buffer := make([]byte, 8, 8)
		n, err := file.Read(buffer)

		if err != nil {
			if errors.Is(err, io.EOF) {
				fmt.Printf("read: %s \n", currentLine)
				break
			}
			fmt.Printf("error: %s\n", err.Error())
			break
		}
		str := string(buffer[:n])

		parts := strings.Split(str, "\n")

		if len(parts) > 1 {
			currentLine += parts[0]
			fmt.Printf("read: %s \n", currentLine)
			currentLine = parts[1]
		} else {
			currentLine += parts[0]
		}

	}
}
