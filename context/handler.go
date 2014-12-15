package context

import (
	"errors"
	"fmt"
	"runtime/debug"

	"github.com/huangml/dispatch/msg"
)

type HandlerFunc func(ctx *Context, m Mutex, r msg.Request) msg.Response

func (f HandlerFunc) Handle(ctx *Context, m Mutex, r msg.Request) msg.Response {
	var rsp msg.Response
	defer func() {
		if err := recover(); err != nil {
			rsp = msg.ErrResponse(fmt.Errorf("PANIC: %s\n%s", err, debug.Stack()))
		}
	}()

	rsp = f(ctx, m, r)
	return rsp
}

type LockedHandlerFunc func(r msg.Request) msg.Response

func (f LockedHandlerFunc) ToHandlerFunc() HandlerFunc {
	return func(ctx *Context, m Mutex, r msg.Request) msg.Response {
		if ctx.AcquireOrCancel(m) {
			defer ctx.Release()

			return f(r)
		} else {
			return msg.ErrResponse(errors.New("Request Canceled"))
		}
	}
}
