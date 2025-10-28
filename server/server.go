package server

import (
	"errors"
	"fmt"
	"hash/fnv"
	"io"
	"log/slog"
	"net"
	"strings"
	"sync"
	"sync/atomic"
	"unsafe"
)

const nShard = 1000

func calculateShard(s string) int {
	hasher := fnv.New64()
	_, _ = hasher.Write([]byte(s))
	hash := hasher.Sum64()
	return int(hash % uint64(nShard))
}

type server struct {
	listener     net.Listener
	logger       *slog.Logger
	started      atomic.Bool
	clients      map[int64]net.Conn
	lastClientId int64
	clientsLock  sync.Mutex
	shuttingDown bool
	dbLock       [nShard]sync.RWMutex
	db           [nShard]map[string]string
}

func NewServer(listener net.Listener, logger *slog.Logger) *server {
	s := &server{
		listener:     listener,
		logger:       logger,
		started:      atomic.Bool{},
		clients:      make(map[int64]net.Conn),
		lastClientId: 0,
		clientsLock:  sync.Mutex{},
		shuttingDown: false,
		db:           [nShard]map[string]string{},
		dbLock:       [nShard]sync.RWMutex{},
	}

	for i := 0; i < nShard; i++ {
		s.db[i] = make(map[string]string)
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
		s.clients[clientId] = conn
		s.clientsLock.Unlock()
		go s.handleConn(clientId, conn)
	}
	return nil
}

func (s *server) handleConn(clientId int64, conn net.Conn) {
	//reader := bufio.NewReader(conn)
	reader := newMessageReader(conn)
	for {
		//request, err := readArray(reader, true)
		length, err := reader.ReadArrayLen()
		if err != nil {
			if !errors.Is(err, io.EOF) {
				s.logger.Error("error reading from client", "id", clientId, "err", err)
			}
			break
		}
		if length == 0 {
			//	s.logger.Error("missing command in the request", slog.Int64("clientId", clientId))
			break
		}

		//commandName, ok := request[0].(string)
		commandName, err := reader.ReadString()
		if err != nil {
			break
		}
		unsafeToUpper(commandName)

		switch strings.ToUpper(commandName) {
		case "GET":
			err = s.handleGetCommand(reader, conn)
		case "SET":
			err = s.handleSetCommand(reader, conn)
		default:
		}

		if _, err := conn.Write([]byte("+OK\r\n")); err != nil {
			// s.logger.Error(
			// 	"error writing to client",
			// 	slog.Int64("clientId", clientId),
			// 	slog.String("err", err.Error()),
			// )
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

	if err := conn.Close(); err != nil {
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

	for clientId, conn := range s.clients {
		if err := conn.Close(); err != nil {
			s.logger.Error("cannot close client", "clientId", clientId, "err", err)
		}
	}
	clear(s.clients)
	if err := s.listener.Close(); err != nil {
		s.logger.Error("cannot stop listener", "err", err)
	}
	return nil
}

func (s *server) handleGetCommand(reader *messageReader, conn io.Writer) error {
	key, err := reader.ReadString()
	if err != nil {
		_, err = conn.Write([]byte("-ERR key is not a string\r\n"))
		return err
	}
	shard := calculateShard(key)
	s.dbLock[shard].RLock()
	value, ok := s.db[shard][key]
	s.dbLock[shard].RUnlock()
	if ok {
		resp := fmt.Sprintf("$%d\r\n%s\r\n", len(value), value)
		_, err = conn.Write([]byte(resp))
	} else {
		_, err = conn.Write([]byte("_\r\n"))
	}

	return err
}
func (s *server) handleSetCommand(reader *messageReader, conn io.Writer) error {
	key, err := reader.ReadString()
	if err != nil {
		_, err := conn.Write([]byte("-ERR key is not a string\r\n"))
		return err
	}
	value, err := reader.ReadString()
	if err != nil {
		_, err := conn.Write([]byte("-ERR value is not a string\r\n"))
		return err
	}
	shard := calculateShard(key)
	s.dbLock[shard].Lock()
	s.db[shard][key] = value
	s.dbLock[shard].Unlock()
	_, err = conn.Write([]byte("+OK\r\n"))
	return err
}

func unsafeToUpper(s string) {
	bytes := unsafe.Slice(unsafe.StringData(s), len(s))
	for i := 0; i < len(bytes); i++ {
		b := bytes[i]
		if b >= 'a' && b <= 'z' {
			b = b + 'A' - 'a'
			bytes[i] = b
		}
	}
}
