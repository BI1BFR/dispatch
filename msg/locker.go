package msg

import "sync"

type Locker struct {
	l      sync.Locker
	locked bool
}

func newLocker(l sync.Locker) *Locker {
	return &Locker{
		l: l,
	}
}

func (l *Locker) Lock() {
	l.l.Lock()
	l.locked = true
}

func (l *Locker) Unlock() {
	l.l.Unlock()
	l.locked = false
}
