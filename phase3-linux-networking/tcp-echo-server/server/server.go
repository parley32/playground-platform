package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"
)

const idleTimeout = 30 * time.Second
const bufferSize = 4096

type ServerValues interface {
	getTimeOut() time.Duration
	getBufferSize() int
}

type prodValues struct{}

func (prodValues) getTimeOut() time.Duration {
	return idleTimeout
}

func (prodValues) getBufferSize() int {
	return bufferSize
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
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	ln, err := net.Listen("tcp", "127.0.0.1:8080")
	if err != nil {
		fmt.Println("failed to listen:", err)
		return
	}
	defer ln.Close()
	fmt.Println("listening on:", ln.Addr())
	store := newStore()
	serve(ctx, store, ln, prodValues{})
	stop()
}

func serve(ctx context.Context, store *store, ln net.Listener, values ServerValues) error {
	var wg sync.WaitGroup

	stop := context.AfterFunc(ctx, func() {
		ln.Close()
	})
	defer stop()

	for {
		conn, err := ln.Accept()
		if errors.Is(err, net.ErrClosed) {
			fmt.Println("server shutting down")
			break
		}
		if err != nil {
			fmt.Println("connection failed:", err)
		} else {
			wg.Go(func() {
				handleConn(ctx, store, values, conn)
			})
		}

	}
	wg.Wait()
	return nil
}

func clientShutdown(conn net.Conn) {
	err := conn.SetReadDeadline(time.Now())
	if err != nil {
		fmt.Println("failed to set read deadline:", err)
		return
	}
}

func handleConn(ctx context.Context, s *store, serverValues ServerValues, conn net.Conn) {
	defer conn.Close()
	fmt.Println("client connected:", conn.RemoteAddr())

	reader := bufio.NewReaderSize(conn, serverValues.getBufferSize())

	stop := context.AfterFunc(ctx, func() {
		clientShutdown(conn)
	})
	defer stop()

	for {
		err := conn.SetReadDeadline(time.Now().Add(serverValues.getTimeOut()))
		if err != nil {
			fmt.Println("failed to set read deadline:", err)
			return
		}
		if ctx.Err() != nil {
			setShutdownError(err, conn)
			return
		}
		line, err := reader.ReadSlice('\n')
		if errors.Is(err, bufio.ErrBufferFull) {
			err = loopRead(reader)
			if err == nil {
				fmt.Println(conn.RemoteAddr(), " max input size: ", serverValues.getBufferSize())
				_, err = conn.Write([]byte("ERR max input size: " + strconv.Itoa(serverValues.getBufferSize()) + "\n"))
				if err != nil {
					fmt.Println("write error:", err)
					return
				}
				continue
			}
		}
		if err != nil {
			if errors.Is(err, os.ErrDeadlineExceeded) {
				if ctx.Err() != nil {
					setShutdownError(err, conn)
				} else {
					fmt.Println("client timed out:", conn.RemoteAddr())
				}
			} else if errors.Is(err, io.EOF) {
				fmt.Println("client disconnected:", conn.RemoteAddr())
			} else {
				fmt.Println("error occurred:", err)
			}
			return
		}

		_, err = conn.Write([]byte(handleCommand(s, string(line)) + "\n"))
		if err != nil {
			fmt.Println("write error:", err)
			return
		}
	}
}

func setShutdownError(err error, conn net.Conn) {
	fmt.Println("server shutting down, closing connection:", conn.RemoteAddr())
	_, err = conn.Write([]byte("server shutting down, closing connection\n"))
	if err != nil {
		fmt.Println("write error: ", err)
	}
}

func loopRead(reader *bufio.Reader) error {
	for {
		_, err := reader.ReadSlice('\n')
		if !errors.Is(err, bufio.ErrBufferFull) {
			return err
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
