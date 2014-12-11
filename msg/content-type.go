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

// func (c ContentType) ToHTTP() string {
// 	switch c {
// 	case Bytes:
// 		return "application/octet-stream"
// 	case Text:
// 		return "text/plain"
// 	case ProtoBuf:
// 		return "application/x-protobuf"
// 	default:
// 		return "application/octet-stream"
// 	}
// }
