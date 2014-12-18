package dispatch

import "github.com/huangml/mux"

// Dest is an in-memory entity that accepts Requests and returns Responses.
// It's typically implemented with a Mutex to make Requests processed
// sequentially and one or more Handlers to do the logic.
//
// Call() is used for synchronous communication, it returns after Request has
// been processed; Send() is used for asynchronous communication, it returns
// immediately and spawns a new goroutine to process Request.
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

// NewConcurrentDest creates a ConcurrentDest with provided Handler.
func NewConcurrentDest(h Handler) *ConcurrentDest {
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
	Mutex
}

// NewLockedDest creates a LockedDest with provided Handler.
func NewLockedDest(h Handler) *LockedDest {
	return &LockedDest{h: h, Mutex: NewMutex()}
}

// Call is for synchronous communication.
func (d *LockedDest) Call(ctx *Context, r Request) Response {
	return d.h.Serve(ctx, d.Mutex, r)
}

// Send is for asynchronous communication.
func (d *LockedDest) Send(r Request) error {
	go d.h.Serve(NewContext(), d.Mutex, r)
	return nil
}

// MuxDest is a Dest carries mutiple Handlers.
type MuxDest struct {
	mu *mux.Mux
	Mutex
}

// NewMuxDest creates a MuxDest with provided Mux.
// Note that passed Mux MUST be either empty or only has values of type Handler.
func NewMuxDest(m *mux.Mux) *MuxDest {
	return &MuxDest{mu: m, Mutex: NewMutex()}
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
		return h.(Handler).Serve(ctx, d.Mutex, r)
	}
	return SimpleResponse(nil, ProtocolNotImplementError(r.Protocol()))
}

// Send is for asynchronous communication. It returns a Response with
// ProtocolNotImplementError if no registered Handler matches the Request's
// protocol.
func (d *MuxDest) Send(r Request) error {
	if h := d.mu.Match(r.Protocol()); h != nil {
		go h.(Handler).Serve(NewContext(), d.Mutex, r)
		return nil
	}
	return ProtocolNotImplementError(r.Protocol())
}
