package dispatch

import (
	"errors"
	"fmt"

	"github.com/huangml/dispatch/msg"
)

type Dest interface {
	Call(*msg.Context, msg.Request) msg.Response
	Send(msg.Request) error
}

type AddressBook interface {
	Lookup(r msg.Request) Dest
}

type Dispatcher struct {
	b AddressBook
}

func (d *Dispatcher) Call(ctx *msg.Context, r msg.Request) msg.Response {
	if dst := d.b.Lookup(r); dst != nil {
		return dst.Call(ctx, r)
	}
	return msg.ErrResponse(fmt.Errorf("dest not found, protocol: %v", r.Protocol()))
}

func (d *Dispatcher) Send(r msg.Request) error {
	if dst := d.b.Lookup(r); dst != nil {
		return dst.Send(r)
	}
	return errors.New("dest not found, protocol: " + r.Protocol())
}
