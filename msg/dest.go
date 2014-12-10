package msg

import (
	"errors"
	"fmt"
	"runtime/debug"
	"sync"
)

type Handler func(l *Locker, r *Request) *Response
type LockedHandler func(r *Request) *Response

type Dest struct {
	name string

	handlerMtx sync.RWMutex
	handlers   map[string]Handler

	mtx sync.Mutex
}

func NewDest(name string) *Dest {
	return &Dest{
		name:     name,
		handlers: make(map[string]Handler),
	}
}

func (d *Dest) SetHandler(protocol string, h Handler) {
	d.handlerMtx.Lock()
	defer d.handlerMtx.Unlock()

	d.handlers[protocol] = h
}

func (d *Dest) SetLockedHandler(protocol string, h LockedHandler) {
	d.SetHandler(protocol, func(l *Locker, r *Request) *Response {
		l.Lock()
		defer l.Unlock()

		return h(r)
	})
}

func (d *Dest) getHandler(protocol string) Handler {
	d.handlerMtx.RLock()
	defer d.handlerMtx.RUnlock()

	h, ok := d.handlers[protocol]
	if !ok {
		h, _ = d.handlers["*"]
	}

	return h
}

func (d *Dest) handle(req *Request) (rsp *Response) {

	if h := d.getHandler(req.Protocol); h != nil {
		defer func() {
			if err := recover(); err != nil {
				rsp = ErrResponse(errors.New(fmt.Sprintf("PANIC: %s\n%s", err, debug.Stack())))
			}
		}()

		rsp = h(newLocker(&d.mtx), req)

	} else {
		rsp = ErrResponse(errors.New("protocol not implemented"))
	}

	return
}

func (d *Dest) Call(l *Locker, req *Request) *Response {
	if l.locked {
		panic("dest.Call() with resource locked")
	}

	return d.handle(req)
}

func (d *Dest) Send(req *Request) {
	go d.handle(req)
}
