package server

import (
	"bytes"
	"io"
	"log/slog"
	"math/rand"
	"net"
	"sync/atomic"
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

func BenchmarkRedisSet(b *testing.B) {
	// Use TCP instead of Unix socket for Windows
	listener, err := net.Listen("tcp", "127.0.0.1:0") // 0 => random available port
	if err != nil {
		b.Fatalf("cannot open TCP listener: %v", err)
	}
	defer listener.Close()

	address := listener.Addr().String()
	noopLogger := slog.New(slog.NewTextHandler(io.Discard, nil))
	server := NewServer(listener, noopLogger)

	// Start the server in the background
	go func() {
		if err := server.Start(); err != nil {
			b.Errorf("cannot start server: %v", err)
		}
	}()

	b.ResetTimer()

	var id atomic.Int64

	b.RunParallel(func(pb *testing.PB) {
		client, err := net.Dial("tcp", address)
		if err != nil {
			b.Fatalf("cannot connect to server: %v", err)
		}
		defer client.Close()

		randomizer := rand.New(rand.NewSource(id.Add(1)))
		pipelineSize := 100

		buff := make([]byte, 4096)
		writeBuffer := bytes.Buffer{}
		count := 0

		for pb.Next() {
			writeBuffer.WriteString("*3\r\n$3\r\nset\r\n$12\r\n")
			for i := 0; i < 12; i++ {
				writeBuffer.WriteByte(byte(randomizer.Int31()%96 + 32))
			}
			writeBuffer.WriteString("\r\n$12\r\n")
			for i := 0; i < 12; i++ {
				writeBuffer.WriteByte(byte(randomizer.Int31()%96 + 32))
			}
			writeBuffer.WriteString("\r\n")
			count++

			if count >= pipelineSize {
				if _, err := writeBuffer.WriteTo(client); err != nil {
					b.Fatalf("cannot write to server: %v", err)
				}
				if _, err := io.ReadFull(client, buff[:5*count]); err != nil {
					b.Fatalf("cannot read from server: %v", err)
				}
				count = 0
			}
		}

		// Flush remaining commands
		if count > 0 {
			if _, err := writeBuffer.WriteTo(client); err != nil {
				b.Fatalf("cannot write to server: %v", err)
			}
			if _, err := io.ReadFull(client, buff[:5*count]); err != nil {
				b.Fatalf("cannot read from server: %v", err)
			}
		}
	})

	b.StopTimer()

	if err := server.Stop(); err != nil {
		b.Errorf("cannot stop server: %v", err)
		return
	}

	throughput := float64(b.N) / b.Elapsed().Seconds()
	b.ReportMetric(throughput, "ops/sec")
}
