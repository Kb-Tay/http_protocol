package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"strings"
)

const BYTES_READ = 8;
const port = ":42069"

func main() {
	listener, err := net.Listen("tcp", ":42069")
	if err != nil {
		log.Fatal("Failed to create listener")
	}

	defer listener.Close()

	for {
		conn, err := listener.Accept()
		
		if err != nil {
			break
		}

		fmt.Println("Connection Accepted")
		ch := getLinesChannel(conn)
		for  {
			line, ok := <- ch 

			if !ok {
				break;
			}

			fmt.Printf("%s\n", line)
		}
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

func getLinesChannel(conn net.Conn) <-chan string {
	ch := make(chan string)
	go func() {
		readFromConn(conn, ch)
		close(ch)
		conn.Close()
		fmt.Println("Connection Closed")
	}()

	return ch
}

func readFromConn(conn net.Conn, ch chan string) {
	var input = strings.Builder{}
	
	for {
		buf := make([]byte, BYTES_READ)
		_, err := conn.Read(buf)

		if err != nil && err != io.EOF {
			break
		}

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
}

