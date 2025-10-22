package main

import (
	"log"
	"log/slog"
	"net"
)

const defaultListenAddress = ":6379"

type Config struct {
	ListenAddress string
}

type Server struct {
	ln net.Listener
	Config
	peers     map[*Peer]bool
	addPeerch chan *Peer
	quitCh    chan struct{}
}

func NewServer(cfg Config) *Server {
	if len(cfg.ListenAddress) == 0 {
		cfg.ListenAddress = defaultListenAddress
	}
	return &Server{
		Config:    cfg,
		peers:     make(map[*Peer]bool),
		addPeerch: make(chan *Peer),
		quitCh:    make(chan struct{}),
	}
}

func (s *Server) Start() error {
	ln, err := net.Listen("tcp", s.ListenAddress)
	if err != nil {
		return err
	}

	s.ln = ln
	go s.loop()

	slog.Info("server running", "listenAddress", s.ListenAddress)

	return s.acceptLoop()
}

func (s *Server) loop() {
	for {
		select {
		case <-s.quitCh:
			return
		case peer := <-s.addPeerch:
			s.peers[peer] = true
		}
	}
}

func (s *Server) acceptLoop() error {
	for {
		conn, err := s.ln.Accept()
		if err != nil {
			slog.Error("Accept error", "err", err)
			return err
		}
		go s.handleConn(conn)
	}
}

func (s *Server) handleConn(conn net.Conn) {
	peer := NewPeer(conn)
	s.addPeerch <- peer
	slog.Info("new peer connected", "remoteAddress", peer.conn.RemoteAddr())
	if err := peer.readLoop(); err != nil {
		slog.Error("read error", "err", err, "remoteAddress", conn.RemoteAddr())
	}
}

func main() {
	server := NewServer(Config{})
	log.Fatal(server.Start())
}
