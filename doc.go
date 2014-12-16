// Package dispatch implements a message processing mechanism.
//
// The basic idea is goroutine-per-request, which means each request runs in a
// separate goroutine.
package dispatch
