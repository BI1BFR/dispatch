package dispatch

import (
	"runtime/debug"
	"time"
)

type Mutex chan struct{}

func NewMutex() Mutex {
	return make(chan struct{}, 1)
}

type Context struct {
	m      []Mutex
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
	c := make(chan struct{})
	return &Context{cancel: c}, func() { close(c) }
}

func (ctx *Context) AcquireOrCancel(m Mutex) bool {
	select {
	case m <- struct{}{}:
		ctx.m = append(ctx.m, m)
		return true
	case <-ctx.cancel:
		return false
	}
}

func (ctx *Context) Release() {
	if n := len(ctx.m); n > 0 {
		<-ctx.m[n-1]
		ctx.m = ctx.m[:n-1]
	} else {
		panic("Context is not obtaining any Mutex")
	}
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
		return &SimpleResponse{Err: ContextCanceledError{}}
	}
}

func safeServe(f HandlerFunc, ctx *Context, m Mutex, r Request) (rsp Response) {
	defer func() {
		if err := recover(); err != nil {
			rsp = &SimpleResponse{Err: PanicError{err, debug.Stack()}}
		}
	}()

	return f(ctx, m, r)
}
