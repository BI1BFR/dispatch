package dispatch

import "github.com/huangml/mux"

type Dest interface {
	Call(ctx *Context, r Request) Response
	Send(r Request) error
}

type LockedDest struct {
	h Handler
	m Mutex
}

func NewLockedDest(h Handler) *LockedDest {
	return &LockedDest{h: h, m: NewMutex()}
}

func (d *LockedDest) Call(ctx *Context, r Request) Response {
	return d.h.Serve(ctx, d.m, r)
}

func (d *LockedDest) Send(r Request) error {
	go d.h.Serve(NewContext(), d.m, r)
	return nil
}

type MuxDest struct {
	mu  *mux.Mux
	mtx Mutex
}

func NewMuxDest(m *mux.Mux) *MuxDest {
	return &MuxDest{mu: m, mtx: NewMutex()}
}

func (d *MuxDest) Handle(pattern string, h Handler) {
	d.mu.Map(pattern, h)
}

func (d *MuxDest) Call(ctx *Context, r Request) Response {
	if h := d.mu.Match(r.Protocol()); h != nil {
		return h.(Handler).Serve(ctx, d.mtx, r)
	}
	return &SimpleResponse{Err: ProtocolNotImplementError(r.Protocol())}
}

func (d *MuxDest) Send(r Request) error {
	if h := d.mu.Match(r.Protocol()); h != nil {
		go h.(Handler).Serve(NewContext(), d.mtx, r)
		return nil
	}
	return ProtocolNotImplementError(r.Protocol())
}
