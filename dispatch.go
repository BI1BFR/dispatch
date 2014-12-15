package dispatch

import "fmt"

type AddressBook interface {
	Lookup(r Request) Dest
}

type Dispatcher struct {
	AddressBook
}

func (d *Dispatcher) Call(ctx *Context, r Request) Response {
	if dst := d.Lookup(r); dst != nil {
		return dst.Call(ctx, r)
	} else {
		return ErrResponse(destNotFoundError(r.Protocol()))
	}
}

func (d *Dispatcher) Send(r Request) error {
	if dst := d.Lookup(r); dst != nil {
		return dst.Send(r)
	} else {
		return destNotFoundError(r.Protocol())
	}
}

func destNotFoundError(protocol string) error {
	return fmt.Errorf("Dest not found, protocol: %v", protocol)
}
