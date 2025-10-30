package server

import (
	"fmt"
	"io"
)

func (s *server) handleCommands(commands []string, w io.Writer) {
	cmd := commands[0]
	switch cmd {
	case "GET":
		s.handleGetCommand(commands[1:], w)
	case "SET":
		s.handleSetCommand(commands[1:], w)
	}
}

func (s *server) handleGetCommand(commands []string, conn io.Writer) error {
	key := commands[0]
	value, ok := s.db.Load(key)
	var err error
	if ok {
		resp := fmt.Sprintf("$%d\r\n%s\r\n", len(value.(string)), value.(string))
		_, err = conn.Write([]byte(resp))
	} else {
		_, err = conn.Write([]byte("_\r\n"))
	}

	return err
}
func (s *server) handleSetCommand(commands []string, conn io.Writer) error {
	key := commands[0]
	value := commands[1]
	s.db.Store(key, value)
	_, err := conn.Write([]byte("+OK\r\n"))
	return err
}
