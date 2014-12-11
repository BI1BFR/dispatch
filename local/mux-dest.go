package local

import (
	"errors"
	"sync"

	"github.com/huangml/dispatch/msg"
	"github.com/huangml/dispatch/mux"
)

type MuxDest struct {
	m   *mux.Mux
	mtx sync.Mutex
}

func NewMuxDest() *MuxDest {
	return &MuxDest{m: mux.New()}
}

func (d *MuxDest) HandleFunc(pattern string, h HandlerFunc) {
	d.m.Bind(pattern, h)
}

func (d *MuxDest) Handle(pattern string, h Handler) {
	d.m.Bind(pattern, h)
}

func (d *MuxDest) Call(ctx *msg.Context, r msg.Request) msg.Response {
	if h, _ := d.m.Match(r.Protocol()); h != nil {
		return h.(Handler).Handle(ctx, &d.mtx, r)
	}
	return msg.Err("protocol " + r.Protocol() + " not implemented.")
}

func (d *MuxDest) Send(r msg.Request) error {
	if h, _ := d.m.Match(r.Protocol()); h != nil {
		ctx, _ := msg.NewContext()
		go h.(Handler).Handle(ctx, &d.mtx, r)
		return nil
	}
	return errors.New("protocol " + r.Protocol() + " not implemented.")
}
