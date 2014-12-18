package dispatch

import "github.com/huangml/mux"

// AddressBook is used to get a Dest by a Request(typically by address).
//
// Note that AddressBook is not necessarily to be a Dest container, it may also
// be a factory or a loader, etc.
type AddressBook interface {
	Lookup(r Request) Dest
}

// MuxAddressBook is a AddressBook uses a Mux to look up Dest.
type MuxAddressBook struct {
	m *mux.Mux
}

// Register binds a Dest to an address.
func (b *MuxAddressBook) Register(address string, dest Dest) {
	b.m.Map(address, dest)
}

// Lookup matches Request's address to registed Dests.
func (b *MuxAddressBook) Lookup(r Request) Dest {
	if d := b.m.Match(r.Address()); d != nil {
		return d.(Dest)
	}
	return nil
}

// NewMuxAddressBook creates a MuxAddressBook with provided Mux.
//
// Note that passed Mux MUST be either empty or only has values of type Dest.
func NewMuxAddressBook(m *mux.Mux) *MuxAddressBook {
	return &MuxAddressBook{m: m}
}
