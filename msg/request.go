package msg

type Request interface {
	Protocol() string
	Body() *Sink
}

type SimpleRequest struct {
	Proto string
	*Sink
}

func (s *SimpleRequest) Protocol() string {
	return s.Proto
}

func (s *SimpleRequest) Body() *Sink {
	return s.Sink
}
