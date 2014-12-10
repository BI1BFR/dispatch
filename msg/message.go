package msg

import (
	"net/http"

	"github.com/huangml/dispatch/sink"
)

type Request struct {
	Protocol string
	Context  map[string]*sink.Sink
	*sink.Sink
}

func NewRequest(protocol string, s *sink.Sink) *Request {
	return &Request{
		Protocol: protocol,
		Context:  make(map[string]*sink.Sink),
		Sink:     s,
	}
}

type Response struct {
	Err *MsgError
	*sink.Sink
}

func NewResponse(s *sink.Sink) *Response {
	return &Response{
		Sink: s,
	}
}

func ErrWithText(statusCode int, text string) *Response {
	return &Response{
		Err: &MsgError{
			StatusCode: statusCode,
			Text:       text,
		},
	}
}

func Err(statusCode int) *Response {
	return &Response{
		Err: &MsgError{
			StatusCode: statusCode,
		},
	}
}

type MsgError struct {
	StatusCode int
	Text       string
}

func (e *MsgError) Error() string {
	return http.StatusText(e.StatusCode) + " [" + e.Text + "]"
}
