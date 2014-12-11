package msg

type ContentType int

const (
	Bytes ContentType = iota + 1
	Text
	ProtoBuf
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
		return "Unknown"
	}
}
