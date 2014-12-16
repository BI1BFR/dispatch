package dispatch

import "github.com/huangml/mux"

// Dest is an in-memory entity that accepts Requests and returns Responses.
// It's typically implemented with a Mutex to make Requests processed
// sequentially and one or more Handlers to do the logic.
//
// Call() is used for synchronous communication, it returns after Request has
// been processed; Send() is used for asynchronous communication, it spawns a
// new goroutine to process Request and returns immediately.
//
// See ConcurrentDest, LockedDest and MuxDest for more details.
type Dest interface {
	Call(ctx *Context, r Request) Response
	Send(r Request) error
}

// ConcurrentDest is a Dest that Requests are processed concurrently.
type ConcurrentDest struct {
	h Handler
}

// NewUnlockedDest creates an ConcurrentDest with provided Handler.
func NewUnlockedDest(h Handler) *ConcurrentDest {
	return &ConcurrentDest{h: h}
}

// Call is for synchronous communication.
func (d *ConcurrentDest) Call(ctx *Context, r Request) Response {
	return d.h.Serve(ctx, NewMutex(), r)
}

// Send is for asynchronous communication.
func (d *ConcurrentDest) Send(r Request) error {
	go d.h.Serve(NewContext(), NewMutex(), r)
	return nil
}

// LockedDest is a Dest that Requests are processed sequentially.
type LockedDest struct {
	h Handler
	m Mutex
}

// NewLockedDest creates a LockedDest with provided Handler.
func NewLockedDest(h Handler) *LockedDest {
	return &LockedDest{h: h, m: NewMutex()}
}

// Call is for synchronous communication.
func (d *LockedDest) Call(ctx *Context, r Request) Response {
	return d.h.Serve(ctx, d.m, r)
}

// Send is for asynchronous communication.
func (d *LockedDest) Send(r Request) error {
	go d.h.Serve(NewContext(), d.m, r)
	return nil
}

// MuxDest is a Dest carries mutiple Handlers.
type MuxDest struct {
	mu  *mux.Mux
	mtx Mutex
}

// NewMuxDest creates a MuxDest with provided Mux.
// Note that passed Mux MUST either be empty or only has values of type Handler.
func NewMuxDest(m *mux.Mux) *MuxDest {
	return &MuxDest{mu: m, mtx: NewMutex()}
}

// Handle binds a Handler to a pattern.
func (d *MuxDest) Handle(pattern string, h Handler) {
	d.mu.Map(pattern, h)
}

// Call is for synchronous communication. It returns a Response with
// ProtocolNotImplementError if no registered Handler matches the Request's
// protocol.
func (d *MuxDest) Call(ctx *Context, r Request) Response {
	if h := d.mu.Match(r.Protocol()); h != nil {
		return h.(Handler).Serve(ctx, d.mtx, r)
	}
	return NewSimpleResponse(nil, ProtocolNotImplementError(r.Protocol()))
}

// Send is for asynchronous communication. It returns a Response with
// ProtocolNotImplementError if no registered Handler matches the Request's
// protocol.
func (d *MuxDest) Send(r Request) error {
	if h := d.mu.Match(r.Protocol()); h != nil {
		go h.(Handler).Serve(NewContext(), d.mtx, r)
		return nil
	}
	return ProtocolNotImplementError(r.Protocol())
}
