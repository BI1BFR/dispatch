package local

import (
	"sync"

	"github.com/huangml/dispatch/msg"
)

type LockedDest struct {
	h   Handler
	mtx sync.Mutex
}

func NewLockedDest(h Handler) *LockedDest {
	return &LockedDest{h: h}
}

func (d *LockedDest) Call(ctx *msg.Context, r msg.Request) msg.Response {
	return d.h.Handle(ctx, &d.mtx, r)
}

func (d *LockedDest) Send(r msg.Request) error {
	ctx, _ := msg.NewContext()
	go d.h.Handle(ctx, &d.mtx, r)
	return nil
}
