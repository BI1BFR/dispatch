package dispatch

import (
	"net/http"
	"time"
)

type Auth struct {
	UserName string
	Password string
}

type RemoteRequest struct {
	Proto string
	*Sink
	*Auth
	TimeOut time.Duration
}

func (r *RemoteRequest) Protocol() string {
	return r.Proto
}

func (r *RemoteRequest) Body() *Sink {
	return r.Sink
}

type RemoteResponse struct {
	*statusError
	*Sink
}

func (r *RemoteResponse) Error() error {
	return r.statusError
}

func (r *RemoteResponse) Body() *Sink {
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
