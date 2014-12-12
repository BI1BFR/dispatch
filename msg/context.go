package msg

import (
	"sync"
	"time"
)

type Context struct {
	l      sync.Locker
	cancel chan struct{}
}

func NewContext() (ctx *Context, cancel func()) {
	c := make(chan struct{})
	return &Context{cancel: c}, func() { close(c) }
}

func NewContextWithTimeOut(t time.Duration) *Context {
	ctx, cancel := NewContext()
	time.AfterFunc(t, cancel)
	return ctx
}

func (ctx *Context) Lock(l sync.Locker) {
	if ctx.l != nil {
		panic("context is already locked")
	}
	if l == nil {
		return
	}

	ctx.l = l
	l.Lock()
}

func (ctx *Context) Unlock() {
	if ctx.l == nil {
		return
	}

	ctx.l.Unlock()
	ctx.l = nil
}

func (ctx *Context) Locked() bool {
	return ctx.l != nil
}

func (ctx *Context) Cancel() <-chan struct{} {
	return ctx.cancel
}
