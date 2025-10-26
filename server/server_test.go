package server

import (
	"io"
	"log/slog"
	"net"
	"testing"
)

func TestServerBootstrap(t *testing.T) {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		t.Fatalf("cannot open TCP listener: %v", err)
	}
	noopLogger := slog.New(slog.NewTextHandler(io.Discard, nil))
	server := NewServer(listener, noopLogger)

	n := 10
	startErrorChan := make(chan error, n)
	for i := 0; i < n; i++ {
		go func() {
			err := server.Start()
			startErrorChan <- err
		}()
	}

	for i := 0; i < n-1; i++ {
		err := <-startErrorChan
		if err == nil {
			t.Error("there should be N-1 invocation of `Start` that returns error")
			return
		}
	}

	client, err := net.Dial("tcp", ":8080")
	if err != nil {
		t.Fatalf("cannot open TCP listener: %v", err)
	}
	if _, err := client.Write([]byte("*3\r\n$3\r\nset\r\n$4\r\nkey1\r\n$6\r\nvalue1\r\n")); err != nil {
		t.Errorf("cannot send 'SET key1 value1' to server: %s", err.Error())
		return
	}

	buff := make([]byte, 4096)
	nRead, err := io.ReadFull(client, buff[:5])
	if err != nil {
		t.Errorf("cannot read server response: %s", err.Error())
		return
	}

	if string(buff[:nRead]) != "+OK\r\n" {
		t.Errorf("calling SET should return 'OK' response, but istead it returns '%s'", string(buff[:nRead]))
		return
	}
}
