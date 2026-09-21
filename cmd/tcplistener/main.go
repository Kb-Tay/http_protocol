package main

import (
	"fmt"
	"http_protocol/internal/headers"
	"http_protocol/internal/request"
	"log"
	"net"
)

const BYTES_READ = 8
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
		request, err := request.RequestFromReader(conn)
		fmt.Printf("Request line:\n- Method: %s\n- Target: %s\n- Version: %s\n%s", request.RequestLine.Method, 
			request.RequestLine.RequestTarget, request.RequestLine.HttpVersion, headerToString(request.Headers))

		conn.Close()
	}
}

func headerToString(headers headers.Headers) string {
	s := "Headers:\n"

	for k, v := range(headers) {
		s += fmt.Sprintf("- %s: %s\n", k, v)
	}
	
	return s
}
