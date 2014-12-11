package msg

import "errors"

type Response interface {
	Error() error
	Body() *Sink
}

type errResponse string

func Err(v string) Response {
	return errResponse(v)
}

func (e errResponse) Error() error {
	return errors.New(string(e))
}

func (e errResponse) Body() *Sink {
	return nil
}

// type errResponse statusError

// func (r errResponse) Body() *Sink {
//         return nil
// }

// func (r errResponse) Error() error {
//         return statusError(r)
// }

// func newErrResponse(statusCode int) Response {
//         return errResponse(newStatusErr(statusCode))
// }

// func newErrResponseWithText(statusCode int, text string) Response {
//         return errResponse(newStatusErrWithText(statusCode, text))
// }

// type StatusError interface {
//         StatusCode() int
//         Error() string
// }

// type statusError struct {
//         statusCode int
//         text       string
// }

// func (e statusError) StatusCode() int {
//         return e.statusCode
// }

// func (e statusError) Error() string {
//         return http.StatusText(e.statusCode) + " [" + e.text + "]"
// }

// func newStatusErr(statusCode int) statusError {
//         return statusError{
//                 statusCode: statusCode,
//         }
// }

// func newStatusErrWithText(statusCode int, text string) statusError {
//         return statusError{
//                 statusCode: statusCode,
//                 text:       text,
//         }
// }
