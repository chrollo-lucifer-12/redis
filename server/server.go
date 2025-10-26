package server

import (
	"fmt"
	"io"
	"log/slog"
	"net"
	"sync/atomic"
)

type server struct {
	listener net.Listener
	logger   *slog.Logger
	started  atomic.Bool
}

func NewServer(listener net.Listener, logger *slog.Logger) *server {
	return &server{
		listener: listener,
		logger:   logger,
		started:  atomic.Bool{},
	}
}

func (s *server) Start() error {
	if !s.started.CompareAndSwap(false, true) {
		return fmt.Errorf("server already started")
	}
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			return err
		}
		go s.handleConn(conn)
	}
	return nil
}

func (s *server) handleConn(conn net.Conn) {
	defer conn.Close()
	for {
		buff := make([]byte, 4096)
		n, err := conn.Read(buff)
		if err != nil {
			if err == io.EOF {
				s.logger.Info("client disconnected", conn.RemoteAddr().String())
			} else if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				s.logger.Error("read connection timeout", "err", err)
			} else {
				s.logger.Error("read error from", "addr", conn.RemoteAddr(), "err", err)
			}
			break
		}
		if n == 0 {
			break
		}
		if _, err := conn.Write(buff[:n]); err != nil {
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				s.logger.Error("write connection timeout", "err", err)
			} else {
				s.logger.Error("read error from", "addr", conn.RemoteAddr(), "err", err)
			}
		}
	}
}

func (s *server) Stop() error {
	if err := s.listener.Close(); err != nil {
		s.logger.Error("cannot stop listener", "err", err)
	}
	return nil
}
