package local

import (
	"fmt"
	"runtime/debug"
	"sync"

	"github.com/huangml/dispatch/msg"
)

type Handler interface {
	Handle(ctx *msg.Context, l sync.Locker, r msg.Request) msg.Response
}

type HandlerFunc func(ctx *msg.Context, l sync.Locker, r msg.Request) msg.Response

func (f HandlerFunc) Handle(ctx *msg.Context, l sync.Locker, r msg.Request) (rsp msg.Response) {
	defer func() {
		if err := recover(); err != nil {
			rsp = msg.Err(fmt.Sprintf("PANIC: %s\n%s", err, debug.Stack()))
		}
	}()

	rsp = f(ctx, l, r)
	return
}

func AutoLockHandlerFunc(f func(r msg.Request) msg.Response) HandlerFunc {
	return func(ctx *msg.Context, l sync.Locker, r msg.Request) msg.Response {
		ctx.Lock(l)
		defer ctx.Unlock()

		return f(r)
	}
}
