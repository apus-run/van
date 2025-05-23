package pipeconn

import (
	"context"
	"errors"
	"net"

	"github.com/apus-run/van/pkg/concurrency"
)

const (
	// Network is the network reported by a pipe's address.
	Network = "pipe"
)

var (
	// ErrClosed indicates that a call to Accept() failed because the listener was closed
	ErrClosed = errors.New("listener was closed")

	// ErrAlreadyClosed indicates that a call to Close() failed because the listener had already been closed.
	ErrAlreadyClosed = errors.New("already closed")

	pipeAddr = func() net.Addr {
		c1, c2 := net.Pipe()
		addr := c1.RemoteAddr()
		_ = c1.Close()
		_ = c2.Close()
		return addr
	}()
)

// DialContextFunc is a function for dialing a pipe listener.
type DialContextFunc func(context.Context) (net.Conn, error)

type pipeListener struct {
	closed       concurrency.Signal
	serverConnsC chan net.Conn
}

// NewPipeListener returns a net.Listener that accepts connections which are local pipe connections (i.e., via
// net.Pipe()). It also returns a function that implements a context-aware dial.
func NewPipeListener() (net.Listener, DialContextFunc) {
	lis := &pipeListener{
		closed:       concurrency.NewSignal(),
		serverConnsC: make(chan net.Conn),
	}

	return lis, lis.DialContext
}

func (l *pipeListener) Accept() (net.Conn, error) {
	if l.closed.IsDone() {
		return nil, ErrClosed
	}
	select {
	case conn := <-l.serverConnsC:
		return conn, nil
	case <-l.closed.Done():
		return nil, ErrClosed
	}
}

func (l *pipeListener) DialContext(ctx context.Context) (net.Conn, error) {
	if l.closed.IsDone() {
		return nil, ErrClosed
	}

	serverConn, clientConn := net.Pipe()

	select {
	case l.serverConnsC <- serverConn:
		return clientConn, nil
	case <-l.closed.Done():
		return nil, ErrClosed
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

func (l *pipeListener) Addr() net.Addr {
	return pipeAddr
}

func (l *pipeListener) Close() error {
	if !l.closed.Signal() {
		return ErrAlreadyClosed
	}
	return nil
}
