package server

import (
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"
)

func (s *server) handleCommands(commands []string, w io.Writer) error {
	cmd := commands[0]
	var err error
	switch cmd {
	case "GET":
		err = s.handleGetCommand(commands[1:], w)
	case "SET":
		err = s.handleSetCommand(commands[1:], w)
	case "LPUSH":
		err = s.handleLpushCommand(commands[1:], w)
	case "RPUSH":
		err = s.handleRpushCommand(commands[1:], w)
	case "LPOP":
		err = s.handleLpopCommand(commands[1:], w)
	case "RPOP":
		err = s.handleRpopCommand(commands[1:], w)
	case "LLEN":
		err = s.handleLlenCommand(commands[1:], w)
	case "LRANGE":
		err = s.handleLRangeCommand(commands[1:], w)
	case "LTRIM":
		err = s.handleTrimCommand(commands[1:], w)
	}
	return err
}

func (s *server) handleTrimCommand(commands []string, conn io.Writer) error {
	key := commands[0]
	startStr := commands[1]
	stopStr := commands[2]

	start, _ := strconv.Atoi(startStr)
	stop, _ := strconv.Atoi(stopStr)
	s.listDB.LTRIM(key, start, stop)

	_, err := conn.Write([]byte("+OK\r\n"))
	return err
}

func (s *server) handleLRangeCommand(commands []string, conn io.Writer) error {
	key := commands[0]
	startStr := commands[1]
	stopStr := commands[2]

	start, _ := strconv.Atoi(startStr)
	stop, _ := strconv.Atoi(stopStr)

	elements := s.listDB.LRANGE(key, start, stop)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("*%d\r\n", len(elements)))

	for _, el := range elements {
		sb.WriteString(fmt.Sprintf("$%d\r\n%s\r\n", len(el), el))
	}

	_, err := conn.Write([]byte(sb.String()))
	return err
}

func (s *server) handleLpopCommand(commands []string, conn io.Writer) error {
	key := commands[0]
	countStr := commands[1]
	count, _ := strconv.Atoi(countStr)

	elements := s.listDB.LPOP(key, count)

	var resp string
	if len(elements) == 0 {
		if count == 1 {
			resp = "_\r\n"
		} else {
			resp = "*0\r\n"
		}
	} else {
		if count == 1 {
			resp = fmt.Sprintf("$%d\r\n%s\r\n", len(elements[0]), elements[0])
		} else {
			var b strings.Builder
			b.WriteString(fmt.Sprintf("*%d\r\n", len(elements)))
			for _, el := range elements {
				b.WriteString(fmt.Sprintf("$%d\r\n%s\r\n", len(el), el))
			}
			resp = b.String()
		}
	}

	_, err := conn.Write([]byte(resp))
	return err
}

func (s *server) handleRpopCommand(commands []string, conn io.Writer) error {
	key := commands[0]
	countStr := commands[1]
	count, _ := strconv.Atoi(countStr)

	elements := s.listDB.RPOP(key, count)

	var resp string
	if len(elements) == 0 {
		if count == 1 {
			resp = "_\r\n"
		} else {
			resp = "*0\r\n"
		}
	} else {
		if count == 1 {
			resp = fmt.Sprintf("$%d\r\n%s\r\n", len(elements[0]), elements[0])
		} else {
			var b strings.Builder
			b.WriteString(fmt.Sprintf("*%d\r\n", len(elements)))
			for _, el := range elements {
				b.WriteString(fmt.Sprintf("$%d\r\n%s\r\n", len(el), el))
			}
			resp = b.String()
		}
	}

	_, err := conn.Write([]byte(resp))
	return err
}

func (s *server) handleRpushCommand(commands []string, conn io.Writer) error {
	key := commands[0]
	values := commands[1:]
	len := s.listDB.RPUSH(key, values)
	resp := fmt.Sprintf(":%d\r\n", len)
	_, err := conn.Write([]byte(resp))
	return err
}

func (s *server) handleLpushCommand(commands []string, conn io.Writer) error {
	key := commands[0]
	values := commands[1:]
	len := s.listDB.LPUSH(key, values)
	resp := fmt.Sprintf(":%d\r\n", len)
	_, err := conn.Write([]byte(resp))
	return err
}

func (s *server) handleLlenCommand(commands []string, conn io.Writer) error {
	key := commands[0]
	len := s.listDB.LLEN(key)
	resp := fmt.Sprintf(":%d\r\n", len)
	_, err := conn.Write([]byte(resp))
	return err
}

func (s *server) handleGetCommand(commands []string, conn io.Writer) error {
	key := commands[0]
	value, ok := s.db.mp.Load(key)
	var err error
	if ok {
		v := value.(data)
		if !v.expiry.IsZero() && time.Now().After(v.expiry) {
			s.db.mp.Delete(key)
			_, err = conn.Write([]byte("$-1\r\n"))
			return err
		}
		resp := fmt.Sprintf("$%d\r\n%s\r\n", len(v.value), v.value)
		_, err = conn.Write([]byte(resp))
	} else {
		_, err = conn.Write([]byte("$-1\r\n"))
	}

	return err
}

func (s *server) handleSetCommand(commands []string, conn io.Writer) error {
	key := commands[0]
	value := commands[1]
	var exp time.Time
	if len(commands) > 2 {
		if len(commands) >= 4 {
			option := strings.ToUpper(commands[2])
			expValue := commands[3]
			duration, err := strconv.Atoi(expValue)
			if err != nil {
				_, _ = conn.Write([]byte("-ERR invalid expire time\r\n"))
				return nil
			}

			switch option {
			case "EX":
				exp = time.Now().Add(time.Duration(duration) * time.Second)
			case "PX":
				exp = time.Now().Add(time.Duration(duration) * time.Millisecond)
			default:
				_, _ = conn.Write([]byte("-ERR syntax error\r\n"))
				return nil
			}
		} else {
			_, _ = conn.Write([]byte("-ERR syntax error\r\n"))
			return nil
		}
	}
	s.db.mp.Store(key, data{value: value, expiry: exp})
	_, err := conn.Write([]byte("+OK\r\n"))
	return err
}

func (s *server) startCleaner() {
	for {
		time.Sleep(1 * time.Second)
		s.db.mp.Range(func(key, value interface{}) bool {
			val := value.(data)
			if time.Now().After(val.expiry) {
				s.db.mp.Delete(key)
			}
			return true
		})
	}
}
