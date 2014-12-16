package dispatch

//go:generate stringer -type=ContentType -output=enums-gen.go

import (
	"io"

	"github.com/golang/protobuf/proto"
)

type ContentType int

const (
	Bytes ContentType = iota + 1
	Text
	ProtoBuf
)

type Sink struct {
	ContentType
	b []byte
}

func NewTextSink(v string) *Sink {
	return &Sink{
		ContentType: Text,
		b:           []byte(v),
	}
}

func NewBytesSink(b []byte) *Sink {
	return &Sink{
		ContentType: Bytes,
		b:           b,
	}
}

func NewProtoSink(m proto.Message) *Sink {
	if b, err := proto.Marshal(m); err == nil {
		return &Sink{
			ContentType: ProtoBuf,
			b:           b,
		}
	}
	return &Sink{}
}

func (s *Sink) Bytes() []byte {
	return s.b
}

func (s *Sink) String() string {
	return string(s.b)
}

func (s *Sink) UnmarshalProtoMessage(m proto.Message) error {
	return proto.Unmarshal(s.b, m)
}

func (s *Sink) Write(w io.Writer) {
	w.Write(s.b)
}

type Request interface {
	Protocol() string
	Dest() string
	Body() *Sink
}

type SimpleRequest struct {
	protocol string
	dest     string
	body     *Sink
}

func (s *SimpleRequest) Protocol() string {
	return s.protocol
}

func (s *SimpleRequest) Dest() string {
	return s.dest
}

func (s *SimpleRequest) Body() *Sink {
	return s.body
}

func NewSimpleRequest(protocol, dest string, body *Sink) *SimpleRequest {
	return &SimpleRequest{
		protocol: protocol,
		dest:     dest,
		body:     body,
	}
}

type Response interface {
	Error() error
	Body() *Sink
}

type SimpleResponse struct {
	e    error
	body *Sink
}

func (s *SimpleResponse) Error() error {
	return s.e
}

func (s *SimpleResponse) Body() *Sink {
	return s.body
}

func NewSimpleResponse(body *Sink, err error) *SimpleResponse {
	return &SimpleResponse{
		e:    err,
		body: body,
	}
}
