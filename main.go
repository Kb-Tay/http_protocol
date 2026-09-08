package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

func main() {
	const BYTES_READ = 8;
	f, err := os.Open("messages.txt")	
	if err != nil {
		log.Fatal("Failed to read file")	
	}

	defer f.Close()
	
	var input = strings.Builder{}

	for {
		buf := make([]byte, BYTES_READ)
		_, err := f.Read(buf)
	
		var byteToWrite = buf

		for i, ch := range buf {
			if ch == '\n' {
				curr, next := buf[:i], buf[i+1:]
				input.Write(curr)
				fmt.Printf("read: %s\n", input.String())
				input.Reset()
				byteToWrite = next
				break
			}
		}

		input.Write(byteToWrite)	

		if err == io.EOF {
			break;
		}
	}

	if input.Len() > 0 && isEmptyStr(input.String()) {
		fmt.Printf("read: %s\n", input.String())
	}
}

func isEmptyStr(buf string) bool {
	for _, ch := range buf {
		if ch != '\n' {
			return false
		}
	}

	return true
}

