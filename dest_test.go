package dispatch

import (
	"sync"
	"testing"
	"time"
)

func TestConcurrentDest(t *testing.T) {
	h := func(r Request) Response {
		time.Sleep(time.Millisecond * 3)
		return SimpleResponse(nil, nil)
	}

	m := NewConcurrentDest(LockedHandlerFunc(h))

	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		go func() {
			wg.Add(1)
			defer wg.Done()

			m.Call(NewContext(), nil)
		}()
	}

	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-time.After(time.Millisecond * 5):
		t.Error()
	case <-done:
	}
}

func TestLockedDest(t *testing.T) {
	h := func(r Request) Response {
		time.Sleep(time.Millisecond * 3)
		return SimpleResponse(nil, nil)
	}

	m := NewLockedDest(LockedHandlerFunc(h))

	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		go func() {
			wg.Add(1)
			defer wg.Done()

			m.Call(NewContext(), nil)
		}()
	}

	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-time.After(time.Millisecond * 5):
	case <-done:
		t.Error()
	}
}
