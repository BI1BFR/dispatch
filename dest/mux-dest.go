package dest

import (
	"fmt"

	"github.com/huangml/dispatch/context"
	"github.com/huangml/dispatch/msg"
	"github.com/huangml/mux"
)

type MuxDest struct {
	mu  *mux.Mux
	mtx context.Mutex
}

func NewMuxDest() *MuxDest {
	return &MuxDest{mu: mux.NewPathMux(), mtx: context.NewMutex()}
}

func (d *MuxDest) HandleFunc(pattern string, h context.HandlerFunc) {
	d.mu.Map(pattern, h)
}

func (d *MuxDest) Call(ctx *context.Context, r msg.Request) msg.Response {
	if h := d.mu.Match(r.Protocol()); h != nil {
		return h.(context.HandlerFunc)(ctx, d.mtx, r)
	}
	return msg.ErrResponse(protocolNotFoundErr(r.Protocol()))
}

func (d *MuxDest) Send(r msg.Request) error {
	if h := d.mu.Match(r.Protocol()); h != nil {
		go h.(context.HandlerFunc)(context.New(), d.mtx, r)
		return nil
	}
	return protocolNotFoundErr(r.Protocol())
}

func protocolNotFoundErr(protocol string) error {
	return fmt.Errorf("protocol %v not implemented", protocol)
}
