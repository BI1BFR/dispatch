package dispatch

import "fmt"

type ProtocolNotImplementError string

func (e ProtocolNotImplementError) Error() string {
	return fmt.Sprintf("Protocol %s not implemented.", e)
}

type DestNotFoundError string

func (e DestNotFoundError) Error() string {
	return fmt.Sprintf("Dest %s not found.", e)
}

type ContextCanceledError struct{}

func (e ContextCanceledError) Error() string {
	return "Context canceled."
}

type PanicError struct {
	err   interface{}
	stack []byte
}

func (e PanicError) Error() string {
	return fmt.Sprintf("PANIC: \nInfomation: %v \nStack: \n%s", e.err, string(e.stack))
}
