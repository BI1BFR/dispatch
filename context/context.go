package context

import "time"

type Mutex chan struct{}

func NewMutex() Mutex {
	l := make(chan struct{}, 1)
	l <- struct{}{}
	return l
}

type Context struct {
	m      chan Mutex
	cancel chan struct{}
}

func New() *Context {
	ctx, _ := NewWithCancel()
	return ctx
}

func NewWithTimeOut(t time.Duration) *Context {
	ctx, cancel := NewWithCancel()
	time.AfterFunc(t, cancel)
	return ctx
}

func NewWithCancel() (ctx *Context, cancel func()) {
	m := make(chan Mutex, 1)
	c := make(chan struct{})
	return &Context{m: m, cancel: c}, func() { close(c) }
}

func (ctx *Context) AcquireOrCancel(m Mutex) bool {
	if len(ctx.m) > 0 {
		panic("Context has already locked a Mutex")
	}

	select {
	case <-m:
		ctx.m <- m
		return true
	case <-ctx.cancel:
		return false
	}
}

func (ctx *Context) Release() {
	if len(ctx.m) == 0 {
		panic("Context is not obtaining any Mutex")
	}

	<-ctx.m <- struct{}{}
}
