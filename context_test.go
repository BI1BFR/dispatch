package dispatch

import (
	"strconv"
	"sync"
	"testing"
	"time"
)

func TestMutex(t *testing.T) {
	m := NewMutex()
	m <- struct{}{}
	<-m

	select {
	case m <- struct{}{}:
	case <-time.After(time.Nanosecond):
		t.Error()
	}

	select {
	case m <- struct{}{}:
		t.Error()
	case <-time.After(time.Nanosecond):
	}
}

func TestAcquireOrCancel(t *testing.T) {
	m := NewMutex()

	var wg sync.WaitGroup

	go func() {
		wg.Add(1)
		defer wg.Done()

		c1 := NewContext()
		if c1.AcquireOrCancel(m) {
			defer c1.Release()

			time.Sleep(time.Millisecond * 10)
		} else {
			t.Error()
		}
	}()

	go func() {
		wg.Add(1)
		defer wg.Done()

		c2 := NewContextWithTimeOut(time.Millisecond)
		if c2.AcquireOrCancel(m) {
			defer c2.Release()

			t.Error()
		}
	}()

	go func() {
		wg.Add(1)
		defer wg.Done()

		c3 := NewContextWithTimeOut(time.Millisecond * 20)
		if c3.AcquireOrCancel(m) {
			defer c3.Release()
		} else {
			t.Error()
		}
	}()

	wg.Wait()
}

func TestRelease(t *testing.T) {
	m1, m2 := NewMutex(), NewMutex()
	c := NewContext()

	released := func(m Mutex) bool {
		select {
		case m <- struct{}{}:
			<-m
			return true
		default:
			return false
		}
	}

	c.AcquireOrCancel(m1)

	if released(m1) || !released(m2) {
		t.Error()
	}

	c.AcquireOrCancel(m2)
	c.Release()

	if released(m1) || !released(m2) {
		t.Error()
	}

	c.Release()

	if !released(m2) || !released(m2) {
		t.Error()
	}

	func() {
		if c.AcquireOrCancel(m1) {
			defer c.Release()

			if c.AcquireOrCancel(m2) {
				defer c.Release()
			} else {
				t.Error()
			}
		} else {
			t.Error()
		}
	}()

	if !released(m2) || !released(m2) {
		t.Error()
	}

	func() {
		if c.AcquireOrCancel(m1) {
			defer c.ReleaseAll()

			if !c.AcquireOrCancel(m2) {
				t.Error()
			}
		} else {
			t.Error()
		}
	}()

	if !released(m2) || !released(m2) {
		t.Error()
	}
}

func TestHandler(t *testing.T) {
	var val int

	h1 := func(r Request) Response {
		val++
		return SimpleResponse(TextSink(strconv.Itoa(val)), nil)
	}

	r := LockedHandlerFunc(h1).Serve(NewContext(), NewMutex(), nil)
	if r.Body().String() != "1" {
		t.Error()
	}

	var wg sync.WaitGroup
	m := NewMutex()
	h2 := func(r Request) Response {
		time.Sleep(time.Millisecond * 10)
		return h1(r)
	}

	for i := 0; i < 5; i++ {
		go func(i int) {
			wg.Add(1)
			defer wg.Done()

			r := LockedHandlerFunc(h2).Serve(NewContext(), m, nil)
			if r.Body().String() != strconv.Itoa(i+1) {
				t.Error()
			}
		}(i)
	}

	go func() {
		wg.Add(1)
		defer wg.Done()

		time.Sleep(time.Millisecond * 5)

		r := LockedHandlerFunc(h2).Serve(NewContext(), m, nil)
		if r.Body().String() != strconv.Itoa(7) {
			t.Error()
		}
	}()

	wg.Wait()
}
