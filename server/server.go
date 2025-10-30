package server

import (
	"fmt"
	"log/slog"
	"net"
	"sync"
	"sync/atomic"
)

type server struct {
	listener     net.Listener
	logger       *slog.Logger
	started      atomic.Bool
	clients      map[int64]*peer
	lastClientId int64
	clientsLock  sync.Mutex
	shuttingDown bool
	db           *database
	listDB       *listDB
}

func NewServer(listener net.Listener, logger *slog.Logger) *server {
	s := &server{
		listener:     listener,
		logger:       logger,
		started:      atomic.Bool{},
		clients:      make(map[int64]*peer),
		lastClientId: 0,
		clientsLock:  sync.Mutex{},
		shuttingDown: false,
		db:           newDB(),
		listDB:       newListDb(),
	}
	return s
}

func (s *server) Start() error {
	if !s.started.CompareAndSwap(false, true) {
		return fmt.Errorf("server already started")
	}
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			s.clientsLock.Lock()
			isShuttingDown := s.shuttingDown
			s.clientsLock.Unlock()
			if !isShuttingDown {
				return err
			}
			return nil
		}
		s.clientsLock.Lock()
		s.lastClientId += 1
		clientId := s.lastClientId
		msgCh := make(chan []byte)
		createPeer := newPeer(conn, msgCh)
		s.clients[clientId] = createPeer
		s.clientsLock.Unlock()
		go s.startCleaner()
		go createPeer.readLoop()
		go s.handleConn(clientId, createPeer)
	}
	return nil
}

func (s *server) handleConn(clientId int64, p *peer) {

	for rawMsg := range p.msgCh {
		msgStr := string(rawMsg)
		request, err := parseResp(msgStr)
		if err != nil {
			break
		}
		err = s.handleCommands(request, p.conn)
		if err != nil {
			break
		}
	}

	s.clientsLock.Lock()
	if _, ok := s.clients[clientId]; !ok {
		s.clientsLock.Unlock()
		return
	}
	delete(s.clients, clientId)
	s.clientsLock.Unlock()

	if err := p.conn.Close(); err != nil {
		s.logger.Error("cannot close client", "clientId", clientId, "err", err)
	}
}

func (s *server) Stop() error {
	s.clientsLock.Lock()
	defer s.clientsLock.Unlock()

	if s.shuttingDown {
		return fmt.Errorf("already shutting down")
	}
	s.shuttingDown = true

	for clientId, peer := range s.clients {
		if err := peer.conn.Close(); err != nil {
			s.logger.Error("cannot close client", "clientId", clientId, "err", err)
		}
	}
	clear(s.clients)
	if err := s.listener.Close(); err != nil {
		s.logger.Error("cannot stop listener", "err", err)
	}
	return nil
}
