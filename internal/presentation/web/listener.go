package web

import (
	"errors"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/brycedarling/go-practical-microservices/internal/infrastructure/config"
)

// NewListener ...
func NewListener(conf *config.Config) (net.Listener, func(), error) {
	netListener, err := net.Listen("tcp", fmt.Sprintf(":%s", conf.Env.Port))
	if err != nil {
		return nil, nil, err
	}

	tcpListener, ok := netListener.(*net.TCPListener)
	if !ok {
		return nil, nil, errors.New("Couldn't wrap listener")
	}

	l := &listener{
		TCPListener: tcpListener,
		shutdown:    make(chan int),
	}

	return l, l.Shutdown, nil
}

type listener struct {
	*net.TCPListener
	shutdown chan int
}

var _ net.Listener = (*listener)(nil)

func (l *listener) Accept() (net.Conn, error) {
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

func (l *listener) Shutdown() {
	log.Printf("Closing listener on %s", l.Addr())

	close(l.shutdown)

	l.TCPListener.Close()
}

// ErrShutdown ...
var ErrShutdown = errors.New("listener is shutdown")
