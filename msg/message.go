package msg

import "github.com/huangml/dispatch/sink"

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
	Err error
	*sink.Sink
}

func ErrResponse(err error) *Response {
	return &Response{
		Err: err,
	}
}

func NewResponse(s *sink.Sink) *Response {
	return &Response{
		Sink: s,
	}
}
