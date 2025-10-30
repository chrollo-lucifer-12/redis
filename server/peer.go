package server

import "net"

type peer struct {
	conn  net.Conn
	msgCh chan []byte
}

func newPeer(conn net.Conn, msgCh chan []byte) *peer {
	return &peer{
		conn:  conn,
		msgCh: msgCh,
	}
}

func (p *peer) readLoop() {
	buf := make([]byte, 1024)
	for {
		n, err := p.conn.Read(buf)
		if err != nil {
			close(p.msgCh)
			return
		}
		msgBuf := make([]byte, n)
		copy(msgBuf, buf[:n])
		p.msgCh <- msgBuf
	}
}
