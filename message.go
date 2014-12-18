package dispatch

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/golang/protobuf/proto"
)

// ContentType indicates a Sink's data type.
type ContentType int

const (
	Bytes    ContentType = iota + 1 // binary bytes
	Text                            // printable string
	Json                            // json string
	Protobuf                        // protobuf message
)

func (c ContentType) String() string {
	switch c {
	case Bytes:
		return "Bytes"
	case Text:
		return "Text"
	case Json:
		return "Json"
	case Protobuf:
		return "Protobuf"
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

// TextSink creates a Sink by string. The Sink's ContentType will be Text.
func TextSink(v string) *Sink {
	return &Sink{
		ContentType: Text,
		b:           []byte(v),
	}
}

// BytesSink creates a Sink by []byte. The Sink's ContentType will be Bytes.
func BytesSink(b []byte) *Sink {
	return &Sink{
		ContentType: Bytes,
		b:           b,
	}
}

// JsonSink creates a Sink by a struct. The Sink's ContentType will be Json.
func JsonSink(j interface{}) *Sink {
	if b, err := json.Marshal(j); err != nil {
		return &Sink{
			ContentType: Json,
			b:           b,
		}
	} else {
		return nil
	}
}

// ProtoSink creates a Sink by a protobuf message. The Sink's ContentType
// will be Protobuf.
func ProtoSink(m proto.Message) *Sink {
	if b, err := proto.Marshal(m); err == nil {
		return &Sink{
			ContentType: Protobuf,
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

// UnmarshalJson unmarshals JSON data and store the result in v.
func (s *Sink) UnmarshalJson(v interface{}) error {
	return json.Unmarshal(s.b, v)
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
// Address() indicates where the Request will be sent to.
// Body() returns Request's ContentType and raw data.
type Request interface {
	Protocol() string
	Address() string
	Body() *Sink
}

// simpleRequest is a simple implemention of Request.
type simpleRequest struct {
	protocol string
	address  string
	body     *Sink
}

// Protocol indicates the Request's type and how it will be processed.
func (s *simpleRequest) Protocol() string {
	return s.protocol
}

// Address indicates where the Request will be sent to.
func (s *simpleRequest) Address() string {
	return s.address
}

// Body returns Request's ContentType and raw data.
func (s *simpleRequest) Body() *Sink {
	return s.body
}

// SimpleRequest creates a simple Request.
func SimpleRequest(protocol, address string, body *Sink) Request {
	return &simpleRequest{
		protocol: protocol,
		address:  address,
		body:     body,
	}
}

// Response represents a response message corrsponding to a Request.
// Error() returns an none nil error if any error occurs.
// Body() returns Response's ContentType and raw data.
type Response interface {
	Error() error
	Body() *Sink
}

// simpleResponse is a simple implemention of Response.
type simpleResponse struct {
	e    error
	body *Sink
}

// Error returns an none nil error if any error occurs.
func (s *simpleResponse) Error() error {
	return s.e
}

// Body returns Response's ContentType and raw data.
func (s *simpleResponse) Body() *Sink {
	return s.body
}

// SimpleResponse creates a simple Response
func SimpleResponse(body *Sink, err error) Response {
	return &simpleResponse{
		e:    err,
		body: body,
	}
}
