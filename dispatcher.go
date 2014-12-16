package dispatch

// AddressBook is used to get a Dest by a Request(typically by r.Dest()).
// It is left to the package user to implement a concrete AddressBook.
//
// Note that AddressBook is not necessarily to be a Dest container, it may also
// be a factory or a loader, etc.
type AddressBook interface {
	Lookup(r Request) Dest
}

// Dispatcher dispatches Requests to corresponding Dest by looking up an
// AddressBook.
//
// Note that Dispatcher is also a Dest. You can build a multi-level Dispatcher
// easily.
type Dispatcher struct {
	AddressBook
}

// Call dispatches a call to a Dest. It returns a Response with DestNotFoundError
// if no Dest found in the AddressBook.
func (d *Dispatcher) Call(ctx *Context, r Request) Response {
	if dst := d.Lookup(r); dst != nil {
		return dst.Call(ctx, r)
	}
	return NewSimpleResponse(nil, DestNotFoundError(r.Address()))
}

// Send dispatches a send to a Dest. It returns a Response with DestNotFoundError
// if no Dest found in the AddressBook.
func (d *Dispatcher) Send(r Request) error {
	if dst := d.Lookup(r); dst != nil {
		return dst.Send(r)
	}
	return DestNotFoundError(r.Address())
}
