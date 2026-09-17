package request

import (
	"errors"
	"http_protocol/internal/buffer"
	"io"
	"log"
	"strings"
	"unicode"
)

const HTTP_VERSION = "1.1"

type State int 

const (
	Initialised State = iota
	Completed
)

type Request struct {
	RequestLine RequestLine
	State State // 0 init, 1 done
}

type RequestLine struct {
    HttpVersion   string
    RequestTarget string
    Method        string
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	request := Request{
		State: Initialised,	
	} // it will initialised with default values
	
	// read chunks within a loop 
	buf := buffer.New()

	for request.State != Completed {
		bytes := make([]byte, 8)
		n, err := reader.Read(bytes)
		if err != nil {
			if err == io.EOF {
				request.State = Completed
				continue;
			}
		}
		// **Key part: Need to only read the amount you take in
		// when init a slice of fixed size, the bytes slice will contain 
		// an array of x00 bytes
		buf.Read(bytes[:n])		

		n, err = request.parse(buf.GetBuffer())

		if err != nil {
			return &request, err 
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

	if n > 0 {
		r.State = Completed
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
	
	if header != "HTTP" {
		return "", errors.New("Invalid protocol")
	}

	if version != HTTP_VERSION {
		return "", errors.New("Incompatible http version")
	}

	return version, nil
}

