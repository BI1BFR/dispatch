package remote

import (
	"net/http"
	"time"

	"github.com/huangml/dispatch/msg"
)

type Auth struct {
	UserName string
	Password string
}

type Request struct {
	Proto string
	*msg.Sink
	*Auth
	TimeOut time.Duration
}

func (r *Request) Protocol() string {
	return r.Proto
}

func (r *Request) Body() *msg.Sink {
	return r.Sink
}

type Response struct {
	*statusError
	*msg.Sink
}

func (r *Response) Error() error {
	return r.statusError
}

func (r *Response) Body() *msg.Sink {
	return r.Sink
}

type statusError struct {
	statusCode int
	text       string
}

func (e statusError) StatusCode() int {
	return e.statusCode
}

func (e statusError) Error() string {
	return http.StatusText(e.statusCode) + " [" + e.text + "]"
}

func Err(statusCode int) statusError {
	return statusError{
		statusCode: statusCode,
	}
}

func ErrWithText(statusCode int, text string) statusError {
	return statusError{
		statusCode: statusCode,
		text:       text,
	}
}
