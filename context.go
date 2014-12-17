package dispatch

import (
	"runtime/debug"
	"time"
)

// Mutex is implemented as a buffered channel in order to make it selectable.
// Always create Mutex by calling NewMutex() to ensure channel's buffer size
// is 1. Usage:
//
//    m := NewMutex() // create
//    m <- struct{}{} // lock
//    <-m             // unlock
//
type Mutex chan struct{}

func (m Mutex) Lock() {
	m <- struct{}{}
}

func (m Mutex) Unlock() {
	<-m
}

// NewMutex creates an unlocked Mutex.
func NewMutex() Mutex {
	return make(chan struct{}, 1)
}

// Context is a goroutine scope object which carries Mutexes locked by current
// goroutine and a cancelation signal.
//
// A Context MUST only be created when:
// 1. Spawning a new goroutine.
// 2. In a new goroutine which has not been associated with a Context. (e.g. in
// http.HandlerFunc)
// In other cases, always use associated Context and pass it around.
type Context struct {
	m      []Mutex
	cancel chan struct{}
}

// NewContext creates a Context with no cancelation.
func NewContext() *Context {
	ctx, _ := NewContextWithCancel()
	return ctx
}

// NewContextWithTimeOut creates a Context that will be canceled after indicated
// time duration.
func NewContextWithTimeOut(t time.Duration) *Context {
	ctx, cancel := NewContextWithCancel()
	time.AfterFunc(t, cancel)
	return ctx
}

// NewContextWithCancel creates a Context that will be canceled when the
// returned cancel function is called.
// Calling cancel more than once will cause panic.
func NewContextWithCancel() (ctx *Context, cancel func()) {
	c := make(chan struct{})
	return &Context{cancel: c}, func() { close(c) }
}

func (ctx *Context) Canceled() <-chan struct{} {
	return ctx.cancel
}

// AcquireOrCancel should be called before accessing shared resources. It returns
// false if cancel is signaled before success lock. Usage pattern:
//
//    if ctx.AcquireOrCancel(m) {
//            defer ctx.Release()
//            // ... access shared resources ...
//    } else {
//            // ... canceled ...
//    }
//
// A Context can lock mutiple Mutexes. Note that a race condition may occur in
// this case (e.g. one locks A then waits to lock B, other locks B then waits
// to lock A). To simplify the problem, it's recommended always locking only one
// Mutex at a time.
func (ctx *Context) AcquireOrCancel(m Mutex) bool {
	select {
	case m <- struct{}{}:
		ctx.m = append(ctx.m, m)
		return true
	case <-ctx.cancel:
		return false
	}
}

// Release unlocks last acquired Mutex.
func (ctx *Context) Release() {
	if n := len(ctx.m); n > 0 {
		<-ctx.m[n-1]
		ctx.m = ctx.m[:n-1]
	}
}

// ReleaseAll unlocks all acquired Mutexes, in First-In-Last-Out order.
func (ctx *Context) ReleaseAll() {
	for i := len(ctx.m) - 1; i >= 0; i-- {
		<-ctx.m[i]
	}
	ctx.m = []Mutex{}
}

// Handler is an interface wraps Serve() which processes a Request and returns a
// Response.
// Typically, Handler will be registed to a Dest and serve a particular protocol.
type Handler interface {
	Serve(ctx *Context, m Mutex, r Request) Response
}

// HandlerFunc is an adapter from an ordinary function to a Handler.
type HandlerFunc func(ctx *Context, m Mutex, r Request) Response

// Serve calls f(ctx, m, r). If panic occurs, it returns a Response with a
// PanicError.
func (f HandlerFunc) Serve(ctx *Context, m Mutex, r Request) (rsp Response) {
	defer func() {
		if err := recover(); err != nil {
			rsp = NewSimpleResponse(nil, PanicError{err, debug.Stack()})
		}
	}()

	return f(ctx, m, r)
}

// LockedHandlerFunc is an adatper from an ordinary function to a auto-locked
// Handler.
type LockedHandlerFunc func(r Request) Response

// Serve tries to lock Mutex then calls f(r). If ctx is canceled before success
// lock, it returns a Response with a ContextCanceledError; If panic occurs, it
// returns a Response with a PanicError.
func (f LockedHandlerFunc) Serve(ctx *Context, m Mutex, r Request) (rsp Response) {
	if ctx.AcquireOrCancel(m) {
		defer ctx.Release()

		defer func() {
			if err := recover(); err != nil {
				rsp = NewSimpleResponse(nil, PanicError{err, debug.Stack()})
			}
		}()

		return f(r)

	}
	return NewSimpleResponse(nil, ContextCanceledError{})
}
