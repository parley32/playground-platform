package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
	"sync"
	"time"
)

const idleTimeout = 30 * time.Second

type ServerValues interface {
	getTimeOut() time.Duration
}

type prodValues struct{}

func (prodValues) getTimeOut() time.Duration {
	return idleTimeout
}

type store struct {
	mu           sync.Mutex
	storedValues map[string]string
}

func newStore() *store {
	return &store{storedValues: make(map[string]string)}
}

func (s *store) get(key string) (string, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	value, ok := s.storedValues[key]
	return value, ok
}

func (s *store) set(key string, value string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.storedValues[key] = value
}

func main() {
	ln, err := net.Listen("tcp", "127.0.0.1:8080")
	if err != nil {
		fmt.Println("failed to listen:", err)
		return
	}
	defer ln.Close()
	fmt.Println("listening on:", ln.Addr())
	store := newStore()

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("connection failed:", err)
			return
		}
		go handleConn(store, prodValues{}, conn)
	}
}

func handleConn(s *store, serverValues ServerValues, conn net.Conn) {
	defer conn.Close()
	fmt.Println("client connected:", conn.RemoteAddr())

	reader := bufio.NewReader(conn)

	for {
		err := conn.SetReadDeadline(time.Now().Add(serverValues.getTimeOut()))
		if err != nil {
			fmt.Println("failed to set read deadline:", err)
			return
		}
		line, err := reader.ReadString('\n')
		if err != nil {
			if errors.Is(err, os.ErrDeadlineExceeded) {
				fmt.Println("client timed out:", conn.RemoteAddr())
			} else if errors.Is(err, io.EOF) {
				fmt.Println("client disconnected:", conn.RemoteAddr())
			} else {
				fmt.Println("error occurred:", err)
			}
			return
		}

		_, err = conn.Write([]byte(handleCommand(s, line) + "\n"))
		if err != nil {
			fmt.Println("write error:", err)
			return
		}
	}
}

// case-sensitive, for simplicity
func handleCommand(s *store, input string) string {
	input = strings.TrimSpace(input)
	if input == "PING" {
		return "PONG"
	}
	command, after, found := strings.Cut(input, " ")

	switch {
	case command == "ECHO" && found:
		return strings.TrimSpace(after)
	case command == "ECHO":
		return "ERR usage ECHO <text>"
	case command == "GET":
		arg, ok := parseGetArgs(after)
		if !ok {
			return "ERR usage GET <key>"
		}
		value, ok := s.get(arg)
		if !ok {
			return "ERR key not found"
		}
		return value
	case command == "SET":
		key, value, ok := parseSetArgs(after)
		if !ok {
			return "ERR usage SET <key> <value>"
		}
		s.set(key, value)
		return "OK"
	default:
		return "ERR unknown command"
	}
}

func parseGetArgs(input string) (string, bool) {
	arg, after, ok := strings.Cut(strings.TrimSpace(input), " ")
	if (ok && len(after) > 0) || len(arg) == 0 {
		return "", false
	}
	return arg, true
}

func parseSetArgs(input string) (string, string, bool) {
	key, value, ok := strings.Cut(strings.TrimSpace(input), " ")
	if !ok || len(value) == 0 || len(key) == 0 {
		return "", "", false
	}
	return key, strings.TrimSpace(value), true

}
