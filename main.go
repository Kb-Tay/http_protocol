package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

const BYTES_READ = 8;

func main() {
	f, err := os.Open("messages.txt")	
	if err != nil {
		log.Fatal("Failed to read file")	
	}

	ch := getLinesChannel(f)
	
	for  {
		line, ok := <- ch 

		if !ok {
			break;
		}

		fmt.Printf("read: %s\n", line)
	}
}

func isEmptyStr(buf string) bool {
	for _, ch := range buf {
		if ch != '\n' {
			return true 
		}
	}

	return false 
}

func getLinesChannel(f io.ReadCloser) <-chan string {
	ch := make(chan string)
	go func() {
		defer f.Close()
		var input = strings.Builder{}

		for {
			buf := make([]byte, BYTES_READ)
			_, err := f.Read(buf)

			var byteToWrite = buf

			for i, b := range buf {
				if b == '\n' {
					curr, next := buf[:i], buf[i+1:]
					input.Write(curr)
					ch <- input.String()
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
		
		if input.Len() > 0 && !isEmptyStr(input.String()) {
			ch <- input.String()
		}

		close(ch)
	}()

	return ch
}

