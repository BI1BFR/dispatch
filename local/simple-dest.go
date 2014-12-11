package local

import "github.com/huangml/dispatch/msg"

type SimpleDest struct {
	h Handler
}

func NewSimpleDest(h Handler) *SimpleDest {
	return &SimpleDest{h: h}
}

func (d *SimpleDest) Call(ctx *msg.Context, r msg.Request) msg.Response {
	return d.h.Handle(ctx, nil, r)
}

func (d *SimpleDest) Send(r msg.Request) error {
	ctx, _ := msg.NewContext()
	go d.h.Handle(ctx, nil, r)
	return nil
}
