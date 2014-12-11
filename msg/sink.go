package msg

import (
	"io"

	"code.google.com/p/goprotobuf/proto"
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
