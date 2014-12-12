package local

import (
	"errors"
	"fmt"
	"sync"

	"github.com/huangml/dispatch/msg"
	"github.com/huangml/mux"
)

type MuxDest struct {
	m   *mux.Mux
	mtx sync.Mutex
}

func NewMuxDest() *MuxDest {
	return &MuxDest{m: mux.NewPathMux()}
}

func (d *MuxDest) HandleFunc(pattern string, h HandlerFunc) {
	d.m.Map(pattern, h)
}

func (d *MuxDest) Handle(pattern string, h Handler) {
	d.m.Map(pattern, h)
}

func (d *MuxDest) Call(ctx *msg.Context, r msg.Request) msg.Response {
	if h, _ := d.m.Match(r.Protocol()); h != nil {
		return h.(Handler).Handle(ctx, &d.mtx, r)
	}
	return msg.ErrResponse(fmt.Errorf("protocol %v not implemented", r.Protocol()))
}

func (d *MuxDest) Send(r msg.Request) error {
	if h, _ := d.m.Match(r.Protocol()); h != nil {
		ctx, _ := msg.NewContext()
		go h.(Handler).Handle(ctx, &d.mtx, r)
		return nil
	}
	return errors.New("protocol " + r.Protocol() + " not implemented.")
}
