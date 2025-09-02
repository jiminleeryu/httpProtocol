package request

import (
	"fmt"
	"io"
	"bytes"
	"jiminryu.httpProtocol/internal/headers"
)

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

type parserState string
const (
	StateInit parserState = "init"
	StateHeaders parserState = "headers"
	StateDone parserState = "done"
	StateError parserState = "error"
)

type Request struct {
	RequestLine RequestLine
	Headers 	headers.Headers
	state 		parserState
}

func newRequest() *Request {
	return &Request{
		state: StateInit,
		Headers: *headers.NewHeaders(),
	}
}

var ErrorMalformedRequestLine = fmt.Errorf("malformed request line")
var ErrorUnsupportedHttpVersion = fmt.Errorf("unsupported http version")
var ErrorRequestInErrorState = fmt.Errorf("Error request in error state")
var SEPARATOR = []byte("\r\n")

const bufferSize int = 1024

/**
 * Parses the request line of an HTTP request. 
*/
func parseRequestLine(b []byte) (*RequestLine, int, error) {
	idx := bytes.Index(b, SEPARATOR)
	if idx == -1 {
		return nil, 0, nil
	}
	startLine := b[:idx]
	read := idx + len(SEPARATOR)

	parts := bytes.Split(startLine, []byte(" "))
	if len(parts) != 3 {
		return nil, 0, ErrorMalformedRequestLine
	}

	httpParts := bytes.Split(parts[2], []byte("/"))
	if len(httpParts) != 2 || string(httpParts[0]) != "HTTP" || string(httpParts[1]) != "1.1" { 
		return nil, 0, ErrorMalformedRequestLine
	}

	rl := &RequestLine{
		Method : string(parts[0]),
		RequestTarget : string(parts[1]),
		HttpVersion : string(httpParts[1]),
	}
	return rl, read, nil
}

func (r *Request) parse(data []byte) (int, error) {
	read := 0
outer:
	for {
		currentData := data[read:]
		switch r.state {
			case StateError:
				return 0, ErrorRequestInErrorState
			case StateInit:
				rl, n, err := parseRequestLine(currentData)
				if err != nil {
					r.state = StateError
					return 0, err
				}

				if n == 0 {
					break outer // unable to move forward
				}

				r.RequestLine = *rl
				read += n

				r.state = StateDone
				r.state = StateHeaders
			case StateHeaders:

				n, done, err := r.Headers.Parse(currentData)
				if err != nil {
					r.state = StateError
					return 0, err
				}

				if n == 0 { // needs to return already read data
					break outer
				}

				read += n

				if done {
					r.state = StateDone
				}

			case StateDone:
				break outer
			default:
				panic("unknown state")
		}
	}
	return read, nil
}

func (r *Request) done() bool {
	return r.state == StateDone || r.state == StateError
}

func (r *Request) error() bool {
	return r.state == StateError
}

func RequestFromReader(reader io.Reader) (*Request, error){
	request := newRequest()

	// buffer could get overrun, header or body can exceed 1k
	buf := make([]byte, bufferSize)
	bufLen := 0
	for !request.done() {
		n, err := reader.Read(buf[bufLen:])
		// TODO: what to do here, handle EOF
		if err != nil {
			return nil, err
		}

		bufLen += n

		readN, err := request.parse(buf[:bufLen + n])
		if err != nil {
			return nil, err
		}

		copy(buf, buf[readN:bufLen])
		bufLen -= readN
	}
	return request, nil
}