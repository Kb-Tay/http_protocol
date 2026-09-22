package request

import (
	"errors"
	"fmt"
	"http_protocol/internal/buffer"
	"http_protocol/internal/headers"
	"http_protocol/internal/utils"
	"io"
	"log"
	"strconv"
	"strings"
)

const HTTP_VERSION = "1.1"

type State int

const (
	Initialised State = iota
	RequestStateParsingHeaders
	RequestStateParsingBody
	Completed
)

type Request struct {
	RequestLine RequestLine
	Headers headers.Headers 
	Body []byte
	// headers field
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
		Headers: headers.NewHeaders(),
	} // it will initialised with default values

	// read chunks within a loop
	buf := buffer.New()

	for request.State != Completed {
		bytes := make([]byte, 8)
		n, _ := reader.Read(bytes)
		//if err != nil {
		//	if err == io.EOF {
		//		request.State = Completed
		//		continue
		//	}
		//}
		// **Key part: Need to only read the amount you take in
		// when init a slice of fixed size, the bytes slice will contain
		// an array of x00 bytes
		buf.Read(bytes[:n])

		switch (request.State) {
		case Initialised, RequestStateParsingBody:
			n, err := request.parse(buf.GetBuffer())
			if n > 0 {
				buf.Parsed(n)	
			}

			if err != nil {
				return &request, err
			}

		case RequestStateParsingHeaders:
			for request.State != RequestStateParsingBody {
				n, err := request.parseSingle(buf.GetBuffer())
			
				if err != nil {
					return &request, err
				}

				if n == 0 {
					break
				}

				buf.Parsed(n)
			}
		
		case Completed:
		}
	}

	return &request, nil
}

func (r *Request) parseRequestLine(data []byte) (int, error) {
	// still need ot split because we might receive more data than just the headers
	parts := strings.SplitN(string(data), "\r\n", 2)

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

	r.State = RequestStateParsingHeaders

	return len(requestLine) + 2, nil
}

func (r *Request) parseRequestBody(data[] byte) (int, error) {
	v := r.Headers.Get("content-length")

	if v == "" {
		r.State = Completed
		return 0, nil
	}

	contentLength, err := strconv.Atoi(v)
	fmt.Printf("Content-Length: %d, Body: %s\n", contentLength, data)

	if err != nil {
		return 0, errors.New("Invalid Content-Length")
	}

	if contentLength == 0 {
		r.State = Completed
		return 0, nil
	}

	// is there a case where bytes length == 0 but then its due to the buffer not reading yet
	if contentLength > 0 && len(data) == 0 {
		return 0, errors.New("Length of Body smaller than Content-Length")
	}

	bytesToParse := min(len(data), contentLength - len(r.Body))
	r.Body = append(r.Body, data[:bytesToParse]...)

	if len(r.Body) == contentLength {
		r.State = Completed
	}

	return bytesToParse, nil
}

func (r *Request) parseSingle(data []byte) (int, error) {
	n, done, err := r.Headers.Parse(data)

	if done {
		r.State = RequestStateParsingBody 
	}

	return n, err
}

func (r *Request) parse(data []byte) (int, error) {
	// split by \r\n here then loop each section

	// reads the byte until the first \r\n
	// discards the rest of the Request for now
	var n int
  var err error

	switch r.State {
		case Initialised:
		n, err = r.parseRequestLine(data)
	case RequestStateParsingBody:
		n, err = r.parseRequestBody(data)
	}

	return n, err 
}

func isValidMethod(method string) bool {
	return utils.IsAlphabetic(method)
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
