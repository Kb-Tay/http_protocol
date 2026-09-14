package request

import (
	"errors"
	"fmt"
	"io"
	"log"
	"strings"
)

type Request struct {
    RequestLine RequestLine
}

type RequestLine struct {
    HttpVersion   string
    RequestTarget string
    Method        string
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	req := Request{
		RequestLine: RequestLine{
			HttpVersion: "",
			RequestTarget: "",
			Method: "",
		},
	}

	bytes, err := io.ReadAll(reader)

	if err != nil {
		log.Fatal("Failed to read request")
	}

	content := string(bytes)
	arr := strings.Split(content, "\r\n")

	reqLine := strings.Split(arr[0], " ")

	if len(reqLine) < 3{
		return &req, errors.New("Invalid request format") 
	}

	method, target, ver := reqLine[0], reqLine[1], reqLine[2]
	
	req.RequestLine = RequestLine{
		parseHttpVersion(ver), target, method,
	}


	return &req, nil
}


func parseHttpVersion(protocol string) string {
	s := strings.Split(protocol, "/")

	if len(s) < 2 {
		log.Println("Invalid protocol")
	}

	header, ver := s[0], s[1]
	
	if header != "HTTP" {
		log.Println("Invalid protocol")
	}

	return ver
}
