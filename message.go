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
	Proto, Dst string
	*Sink
}

func (s *SimpleRequest) Protocol() string {
	return s.Proto
}

func (s *SimpleRequest) Dest() string {
	return s.Dst
}

func (s *SimpleRequest) Body() *Sink {
	return s.Sink
}

type Response interface {
	Error() error
	Body() *Sink
}

type SimpleResponse struct {
	Err  error
	Body *Sink
}

type errResponse struct {
	e error
}

func ErrResponse(err error) Response {
	return errResponse{e: err}
}

func (e errResponse) Error() error {
	return e.e
}

func (e errResponse) Body() *Sink {
	return nil
}
