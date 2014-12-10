package dispatch

import "github.com/huangml/dispatch/msg"

type Dest interface {
	Call(*msg.Locker, *msg.Request) *msg.Response
	Send(*msg.Request)
}

type AddressBook interface {
	Find(addr string) Dest
}
