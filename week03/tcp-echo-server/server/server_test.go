package main

import (
	"bufio"
	"net"
	"testing"
)

func TestClient(t *testing.T) {
	serverConn, clientConn := net.Pipe()
	defer serverConn.Close()
	defer clientConn.Close()
	testReader := bufio.NewReader(clientConn)

	sent := "test\n"
	want := "echo: test\n"
	go handleConn(serverConn)

	_, err := clientConn.Write([]byte(sent))
	if err != nil {
		t.Fatalf("failed to write to client: %s", err)
	}

	received, err := testReader.ReadString('\n')
	if err != nil {
		t.Fatalf("server failed to respond: %s", err)
	}
	if received != want {
		t.Fatalf("error, got %s, wanted %s", received, want)
	}
}
