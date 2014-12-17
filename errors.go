package dispatch

import "fmt"

// ProtocolNotImplementError records an error when no Handler found by Request's
// Protocol().
type ProtocolNotImplementError string

// Error returns error infomation with the protocol.
func (e ProtocolNotImplementError) Error() string {
	return fmt.Sprintf("Protocol %s not implemented.", string(e))
}

// DestNotFoundError records an error when no Dest found by Request's Address().
type DestNotFoundError string

// Error returns error information with the address.
func (e DestNotFoundError) Error() string {
	return fmt.Sprintf("Dest %s not found.", string(e))
}

// ContextCanceledError records an error when Context canceled before processing
// a Request.
type ContextCanceledError struct{}

// Error returns "Context canceled."
func (e ContextCanceledError) Error() string {
	return "Context canceled."
}

// PanicError records an error when panic occur when calling a Handler.
type PanicError struct {
	err   interface{}
	stack []byte
}

// Error returns panic information and debug stack.
func (e PanicError) Error() string {
	return fmt.Sprintf("PANIC: \nInfomation: %v \nStack: \n%s", e.err, string(e.stack))
}
