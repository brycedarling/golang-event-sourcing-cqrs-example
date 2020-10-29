package web

import (
	"errors"
	"log"
	"net"
	"time"
)

// NewTCPListener ...
func NewTCPListener(address string) (net.Listener, func(), error) {
	netLis, err := net.Listen("tcp", address)
	if err != nil {
		return nil, nil, err
	}

	tcpLis, ok := netLis.(*net.TCPListener)
	if !ok {
		return nil, nil, errors.New("Couldn't wrap listener as TCP")
	}

	l := &tcpListener{
		TCPListener: tcpLis,
		shutdown:    make(chan int),
	}

	return l, l.Shutdown, nil
}

type tcpListener struct {
	*net.TCPListener
	shutdown chan int
}

var _ net.Listener = (*tcpListener)(nil)

func (l *tcpListener) Accept() (net.Conn, error) {
	for {
		// Wait up to one second for a new connection
		l.SetDeadline(time.Now().Add(time.Second))

		conn, err := l.TCPListener.Accept()

		select {
		case <-l.shutdown:
			// Check for the channel being closed
			return nil, ErrShutdown
		default:
			// If the channel is still open, continue as normal
		}

		if err != nil {
			// If this is a timeout, then continue to wait for new connections
			if err, ok := err.(net.Error); ok && err.Timeout() && err.Temporary() {
				continue
			}
		}

		return conn, err
	}
}

func (l *tcpListener) Shutdown() {
	log.Printf("Closing tcp listener on %s", l.Addr())

	close(l.shutdown)

	l.TCPListener.Close()
}

// ErrShutdown ...
var ErrShutdown = errors.New("tcp listener is shutdown")
