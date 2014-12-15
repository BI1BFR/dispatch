package dest

import (
	"github.com/huangml/dispatch/context"
	"github.com/huangml/dispatch/msg"
)

type LockedDest struct {
	h context.HandlerFunc
	m context.Mutex
}

func NewLockedDest(h context.HandlerFunc) *LockedDest {
	return &LockedDest{h: h, m: context.NewMutex()}
}

func (d *LockedDest) Call(ctx *context.Context, r msg.Request) msg.Response {
	return d.h(ctx, d.m, r)
}

func (d *LockedDest) Send(r msg.Request) error {
	go d.h(context.New(), d.m, r)
	return nil
}
