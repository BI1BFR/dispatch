package dispatch

// Dispatcher dispatches Requests to corresponding Dest by looking up an
// AddressBook.
//
// Note that Dispatcher is also a Dest. You can build a multi-level Dispatcher
// easily.
type Dispatcher struct {
	AddressBook
}

// Call dispatches a call request to a Dest. It returns a Response with
// DestNotFoundError if no Dest found in the AddressBook.
func (d *Dispatcher) Call(ctx *Context, r Request) Response {
	if dst := d.Lookup(r); dst != nil {
		return dst.Call(ctx, r)
	}
	return SimpleResponse(nil, DestNotFoundError(r.Address()))
}

// Send dispatches a send request to a Dest. It returns a Response with
// DestNotFoundError if no Dest found in the AddressBook.
func (d *Dispatcher) Send(r Request) error {
	if dst := d.Lookup(r); dst != nil {
		return dst.Send(r)
	}
	return DestNotFoundError(r.Address())
}
