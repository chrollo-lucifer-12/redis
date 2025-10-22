package main

import (
	"fmt"
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
}

func NewServer(cfg Config) *Server {
	if len(cfg.ListenAddress) == 0 {
		cfg.ListenAddress = defaultListenAddress
	}
	return &Server{
		Config:    cfg,
		peers:     make(map[*Peer]bool),
		addPeerch: make(chan *Peer),
	}
}

func (s *Server) Start() error {
	ln, err := net.Listen("tcp", s.ListenAddress)
	if err != nil {
		return err
	}

	s.ln = ln
	go s.loop()
	return s.acceptLoop()
}

func (s *Server) loop() {
	for {
		select {
		case peer := <-s.addPeerch:
			s.peers[peer] = true
		default:
			fmt.Println("defualt")
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

	peer.readLoop()
}

func main() {
	server := NewServer(Config{})
	err := server.Start()
	if err != nil {
		slog.Error("error starting  server", "err", err)
	}
}
