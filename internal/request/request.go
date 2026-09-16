package request

import (
	"errors"
	"fmt"
	"io"
	"log"
	"strings"
	"unicode"
)

const HTTP_VERSION = "1.1"

type Request struct {
	RequestLine RequestLine
	State int // 0 init, 1 done
}

type RequestLine struct {
    HttpVersion   string
    RequestTarget string
    Method        string
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	request := Request{} // it will initialised with default values
	
	// read chunks within a loop 
	data := make([]byte, 0)

	for {
		bytes := make([]byte, 8)
		n, err := reader.Read(bytes)
		if n > 0 {
			// **Key part: Need to only read the amount you take in
			// when init a slice of fixed size, the bytes slice will contain 
			// an array of x00 bytes
			data = append(data, bytes[:n]...) 		
		}

		n, parseErr := request.parse(data)

		if parseErr != nil {
			return &request, parseErr
		}

		if request.State == 1 {
			return &request, nil
		}

		if err == io.EOF {
			break
		}
	}

	return &request, nil
}

func (r *Request) parseRequestLine(data []byte) (int, error) {
	// still need ot split because we might receive more data than just the headers
	parts	:= strings.SplitN(string(data), "\r\n", 2)
	
	if len(parts) < 2 {
		return 0, nil
	}

	requestLine, _ := parts[0], parts[1]
	parts = strings.Split(requestLine, " ")

	if len(parts) < 3 {
		return 0, errors.New("Invalid request format") 
	}

	method, target, protocol := parts[0], parts[1], parts[2]

  if !isValidMethod(method) {
		return 0, errors.New("Invalid method")
	}

	version, err := parseHttpVersion(protocol)

	if err != nil {
		return 0, err
	}
	
	r.RequestLine.RequestTarget = target
	r.RequestLine.Method = method
	r.RequestLine.HttpVersion = version

	return len(requestLine), nil
}

func (r *Request) parse(data []byte) (int, error) {
	// reads the byte until the first \r\n 
	// discards the rest of the Request for now
	n, err := r.parseRequestLine(data)
	fmt.Printf("received: %v\n", n)

	if n > 0 {
		fmt.Printf("Hit: %v", r)
		r.State = 1
		return n, nil
	}

	if err != nil {
		return 0, err 
	}

	return 0, nil 
}

func isValidMethod(method string) bool {
	for _, c := range method {
		if !unicode.IsLetter(c) {
			return false
		}
	}

	return true
}


func parseHttpVersion(protocol string) (string, error) {
	parts := strings.Split(protocol, "/")

	if len(parts) < 2 {
		log.Println("Invalid protocol")
	}

	header, version := parts[0], parts[1]
	fmt.Printf("header: %s %s", header, version)
	
	if header != "HTTP" {
		return "", errors.New("Invalid protocol")
	}

	if version != HTTP_VERSION {
		return "", errors.New("Incompatible http version")
	}

	return version, nil
}

