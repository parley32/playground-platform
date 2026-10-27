package main

import (
	"bufio"
	"errors"
	"io"
	"net"
	"testing"
	"time"
)

type mockValues struct {
	timeCap time.Duration
}

func (m *mockValues) getTimeOut() time.Duration {
	return m.timeCap * time.Millisecond
}

func TestTimeout(t *testing.T) {
	serverConn, clientConn := net.Pipe()
	defer serverConn.Close()
	defer clientConn.Close()
	testReader := bufio.NewReader(clientConn)

	err := clientConn.SetReadDeadline(time.Now().Add(1 * time.Second))
	if err != nil {
		t.Fatalf("test broke: %s", err)
	}

	go handleConn(&mockValues{timeCap: 1}, serverConn)
	_, err = testReader.ReadString('\n')
	if !errors.Is(err, io.EOF) {
		t.Fatalf("failed to close connection on timeout: %s", err)
	}
}

func TestClientClosure(t *testing.T) {
	serverConn, clientConn := net.Pipe()
	defer serverConn.Close()
	defer clientConn.Close()
	testReader := bufio.NewReader(clientConn)

	sent := "ECHO test\n"
	want := "test\n"
	done := make(chan struct{})
	go func() {
		defer close(done)
		handleConn(&mockValues{timeCap: 5000}, serverConn)
	}()

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

	clientConn.Close()

	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("handleConn didnt return as soon as the client closed")
	}
}

func TestClient(t *testing.T) {
	serverConn, clientConn := net.Pipe()
	defer serverConn.Close()
	defer clientConn.Close()
	testReader := bufio.NewReader(clientConn)

	sent := "ECHO test\n"
	want := "test\n"
	go handleConn(&mockValues{timeCap: 5000}, serverConn)

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

func TestHandleCommand(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "Echo command success",
			input: "ECHO test\n",
			want:  "test",
		},
		{
			name:  "Echo command alone",
			input: "ECHO\n",
			want:  "ERR usage ECHO <text>",
		},
		{
			name:  "Ping command success",
			input: "PING\r\n",
			want:  "PONG",
		},
		{
			name:  "Wrongly cased ping command",
			input: "ping\n",
			want:  "ERR unknown command",
		},
		{
			name:  "Echo command with spaces",
			input: "   ECHO   test\n",
			want:  "test",
		},
		{
			name:  "Missing command failure",
			input: "test\n",
			want:  "ERR unknown command",
		},
		{
			name:  "Incorrect command",
			input: "ECHOECHO\n",
			want:  "ERR unknown command",
		},
		{
			name:  "Maintaining inner spaces",
			input: "ECHO a   b \n",
			want:  "a   b",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := handleCommand(test.input); got != test.want {
				t.Fatalf("want %q, got %q", test.want, got)
			}
		})
	}
}
