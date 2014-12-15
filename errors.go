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
