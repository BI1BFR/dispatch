package msg

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
