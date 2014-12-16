package dispatch

import (
	"fmt"
	"io"

	"github.com/golang/protobuf/proto"
)

// ContentType indicates a Sink's data type.
type ContentType int

const (
	Bytes    ContentType = iota + 1 // binary bytes
	Text                            // printable string
	ProtoBuf                        // protobuf message
)

func (c ContentType) String() string {
	switch c {
	case Bytes:
		return "Bytes"
	case Text:
		return "Text"
	case ProtoBuf:
		return "ProtoBuf"
	default:
		return fmt.Sprintf("ContextType(%d)", c)
	}
}

// A Sink carries binary bytes, string, or a protobuf message. Internally datas
// are converted to []byte.
type Sink struct {
	ContentType
	b []byte
}

// NewTextSink creates a Sink by string. The Sink's ContentType will be Text.
func NewTextSink(v string) *Sink {
	return &Sink{
		ContentType: Text,
		b:           []byte(v),
	}
}

// NewBytesSink creates a Sink by []byte. The Sink's ContentType will be Bytes.
func NewBytesSink(b []byte) *Sink {
	return &Sink{
		ContentType: Bytes,
		b:           b,
	}
}

// NewProtoSink creates a Sink by a protobuf message. The Sink's ContentType
// will be ProtoBuf.
func NewProtoSink(m proto.Message) *Sink {
	if b, err := proto.Marshal(m); err == nil {
		return &Sink{
			ContentType: ProtoBuf,
			b:           b,
		}
	}
	return &Sink{}
}

// Bytes returns data as []byte.
func (s *Sink) Bytes() []byte {
	return s.b
}

// String returns data as a string.
func (s *Sink) String() string {
	return string(s.b)
}

// UnmarshalProtoMessage unmarshals data to a protobuf message.
func (s *Sink) UnmarshalProtoMessage(m proto.Message) error {
	return proto.Unmarshal(s.b, m)
}

// Write writes data to a Writer.
func (s *Sink) Write(w io.Writer) {
	w.Write(s.b)
}

// Request represents a request message.
// Protocol() indicates the Request's type and how it will be processed.
// Address() indicates where the Request will be sent.
// Body() carries Request's ContentType and raw data.
type Request interface {
	Protocol() string
	Address() string
	Body() *Sink
}

// SimpleRequest is a simple implemention of Request.
type SimpleRequest struct {
	protocol string
	address  string
	body     *Sink
}

func (s *SimpleRequest) Protocol() string {
	return s.protocol
}

func (s *SimpleRequest) Address() string {
	return s.address
}

func (s *SimpleRequest) Body() *Sink {
	return s.body
}

func NewSimpleRequest(protocol, address string, body *Sink) *SimpleRequest {
	return &SimpleRequest{
		protocol: protocol,
		address:  address,
		body:     body,
	}
}

// Response represents a response message corrsponding to a Request.
// Error() returns an none nil error if any error occurs.
// Body() carries Response's ContentType and raw data.
type Response interface {
	Error() error
	Body() *Sink
}

// SimpleResponse is a simple implemention of Response.
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
