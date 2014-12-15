package dispatch

import (
	"runtime/debug"
	"time"
)

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

func NewContext() *Context {
	ctx, _ := NewContextWithCancel()
	return ctx
}

func NewContextWithTimeOut(t time.Duration) *Context {
	ctx, cancel := NewContextWithCancel()
	time.AfterFunc(t, cancel)
	return ctx
}

func NewContextWithCancel() (ctx *Context, cancel func()) {
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

type Handler interface {
	Serve(ctx *Context, m Mutex, r Request) Response
}

type HandlerFunc func(ctx *Context, m Mutex, r Request) Response

func (f HandlerFunc) Serve(ctx *Context, m Mutex, r Request) (rsp Response) {
	return safeServe(f, ctx, m, r)
}

type LockedHandlerFunc func(r Request) Response

func (f LockedHandlerFunc) Serve(ctx *Context, m Mutex, r Request) Response {
	if ctx.AcquireOrCancel(m) {
		defer ctx.Release()

		return safeServe(func(ctx *Context, m Mutex, r Request) Response {
			return f(r)
		}, ctx, m, r)
	} else {
		return ErrResponse(ContextCanceledError{})
	}
}

func safeServe(f HandlerFunc, ctx *Context, m Mutex, r Request) (rsp Response) {
	defer func() {
		if err := recover(); err != nil {
			rsp = ErrResponse(PanicError{err, debug.Stack()})
		}
	}()

	return f(ctx, m, r)
}
