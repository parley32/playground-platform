package main

import (
	"bufio"
	"errors"
	"io"
	"net"
	"strconv"
	"sync"
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
	store := newStore()
	serverConn, clientConn := net.Pipe()
	defer serverConn.Close()
	defer clientConn.Close()
	testReader := bufio.NewReader(clientConn)

	err := clientConn.SetReadDeadline(time.Now().Add(1 * time.Second))
	if err != nil {
		t.Fatalf("test broke: %s", err)
	}

	go handleConn(store, &mockValues{timeCap: 1}, serverConn)
	_, err = testReader.ReadString('\n')
	if !errors.Is(err, io.EOF) {
		t.Fatalf("failed to close connection on timeout: %s", err)
	}
}

func TestClientClosure(t *testing.T) {
	store := newStore()
	serverConn, clientConn := net.Pipe()
	defer serverConn.Close()
	defer clientConn.Close()
	testReader := bufio.NewReader(clientConn)

	sent := "ECHO test\n"
	want := "test\n"
	done := make(chan struct{})
	go func() {
		defer close(done)
		handleConn(store, &mockValues{timeCap: 5000}, serverConn)
	}()

	_, err := clientConn.Write([]byte(sent))
	if err != nil {
		t.Fatalf("failed to write to client: %q", err)
	}
	received, err := testReader.ReadString('\n')
	if err != nil {
		t.Fatalf("server failed to respond: %q", err)
	}
	if received != want {
		t.Fatalf("error, got %q, wanted %q", received, want)
	}

	clientConn.Close()

	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("handleConn didnt return as soon as the client closed")
	}
}

func TestSharedStore(t *testing.T) {
	serverConn, clientConn := net.Pipe()
	defer serverConn.Close()
	defer clientConn.Close()
	testReader := bufio.NewReader(clientConn)
	store := newStore()

	sent1 := "SET a 1\n"
	want1 := "OK\n"

	serverConn2, clientConn2 := net.Pipe()
	defer serverConn2.Close()
	defer clientConn2.Close()
	testReader2 := bufio.NewReader(clientConn2)
	sent2 := "GET a\n"
	want2 := "1\n"

	channel1 := make(chan struct{})
	channel2 := make(chan struct{})

	go func() {
		defer close(channel1)
		handleConn(store, &mockValues{timeCap: 5000}, serverConn)
	}()

	go func() {
		defer close(channel2)
		handleConn(store, &mockValues{timeCap: 5000}, serverConn2)
	}()

	_, err := clientConn.Write([]byte(sent1))
	if err != nil {
		t.Fatalf("failed to write to client: %q", err)
	}
	received, err := testReader.ReadString('\n')
	if err != nil {
		t.Fatalf("server failed to respond: %q", err)
	}
	if received != want1 {
		t.Fatalf("error, got %q, wanted %q", received, want1)
	}

	_, err = clientConn2.Write([]byte(sent2))
	if err != nil {
		t.Fatalf("failed to write to second client: %q", err)
	}
	received, err = testReader2.ReadString('\n')
	if err != nil {
		t.Fatalf("server failed to respond: %q", err)
	}
	if received != want2 {
		t.Fatalf("error, got %q, wanted %q", received, want2)
	}

	clientConn.Close()
	clientConn2.Close()
	select {
	case <-channel1:
	case <-time.After(time.Second):
		t.Fatal("connection 1 didnt return as soon as the client closed")
	}
	select {
	case <-channel2:
	case <-time.After(time.Second):
		t.Fatal("connection 2 didnt return as soon as the client closed")
	}
}

func TestHandleCommand(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		want       string
		wantKey    string
		wantStored string
		seedKey    string
		seedValue  string
	}{
		{
			name:  "Echo command success",
			input: "ECHO test\n",
			want:  "test",
		},
		{
			name:  "Empty command",
			input: "\n",
			want:  "ERR unknown command",
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
		{
			name:       "Setting a value with whitespaces",
			input:      "SET a   Hello    World   \n",
			want:       "OK",
			wantKey:    "a",
			wantStored: "Hello    World",
		},
		{
			name:  "Incorrectly setting a value with whitespaces only",
			input: "SET a  \n",
			want:  "ERR usage SET <key> <value>",
		}, {
			name:       "Correctly finding the key despite extra whitespaces between command and key",
			input:      "SET     a b \n",
			want:       "OK",
			wantKey:    "a",
			wantStored: "b",
		},
		{
			name:  "Incorrect set command",
			input: "set a b\n",
			want:  "ERR unknown command",
		},
		{
			name:  "Incorrect SET usage",
			input: "SET a\n",
			want:  "ERR usage SET <key> <value>",
		},
		{
			name:  "Invalid SET alone",
			input: "SET\n",
			want:  "ERR usage SET <key> <value>",
		},
		{
			name:  "invalid GET",
			input: "GET a   b\n",
			want:  "ERR usage GET <key>",
		},
		{
			name:  "invalid GET alone",
			input: "GET\n",
			want:  "ERR usage GET <key>",
		},
		{
			name:      "valid GET with trailing spaces",
			input:     "GET a     \n",
			want:      "1",
			seedKey:   "a",
			seedValue: "1",
		},
		{
			name:      "Correct GET",
			input:     "GET c\n",
			want:      "1",
			seedKey:   "c",
			seedValue: "1",
		},
		{
			name:  "invalid GET key",
			input: "GET b\n",
			want:  "ERR key not found",
		},
		{
			name:       "SET correctly overwrites value",
			input:      "SET specialKey valueB\n",
			want:       "OK",
			wantKey:    "specialKey",
			wantStored: "valueB",
			seedKey:    "specialKey",
			seedValue:  "valueA",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			store := newStore()

			if len(test.seedKey) > 0 {
				store.set(test.seedKey, test.seedValue)
			}
			if got := handleCommand(store, test.input); got != test.want {
				t.Fatalf("want %q, got %q", test.want, got)
			}
			if len(test.wantStored) > 0 {
				stored, _ := store.get(test.wantKey)
				if stored != test.wantStored {
					t.Fatalf("want stored %q, got %q", test.wantStored, store.storedValues[test.wantKey])
				}
			}
		})
	}
}

func TestMutex(t *testing.T) {
	store := newStore()
	var wg sync.WaitGroup

	start := make(chan struct{})
	for index := range 50 {
		wg.Go(func() {
			<-start
			store.set(strconv.Itoa(index), strconv.Itoa(index*2))
			value, ok := store.get(strconv.Itoa(index))
			if !ok || value != strconv.Itoa(index*2) {
				t.Errorf("Mutex ERR: want %q, got %q", strconv.Itoa(index*2), value)
			}
		})
	}
	close(start)
	wg.Wait()

	for index := range 50 {
		value, ok := store.get(strconv.Itoa(index))
		if !ok || value != strconv.Itoa(index*2) {
			t.Fatalf("Mutex ERR: want %q, got %q", strconv.Itoa(index*2), value)
		}
	}
}
