package dispatch

import (
	"net/http"

	"github.com/huangml/dispatch/msg"
)

type Dispatcher struct {
	b AddressBook
}

func (d *Dispatcher) Call(addr string, l *msg.Locker, req *msg.Request) *msg.Response {
	if p := d.b.Find(addr); p != nil {
		return p.Call(l, req)
	} else {
		return msg.ErrWithText(http.StatusNotFound, "dest not found")
	}
}

func (d *Dispatcher) Send(addr string, req *msg.Request) error {
	if p := d.b.Find(addr); p != nil {
		p.Send(req)
		return nil
	} else {
		return msg.ErrWithText(http.StatusNotFound, "dest not found").Err
	}
}
