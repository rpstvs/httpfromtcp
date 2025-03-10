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
	lines := getLinesChannel(file)

	for cenas := range lines {
		fmt.Printf("read: %s\n", cenas)
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
