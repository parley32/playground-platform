package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
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

func main() {
	ln, err := net.Listen("tcp", "127.0.0.1:8080")
	if err != nil {
		fmt.Println("failed to listen:", err)
		return
	}
	defer ln.Close()
	fmt.Println("listening on:", ln.Addr())

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("connection failed:", err)
			return
		}
		go handleConn(prodValues{}, conn)
	}
}

func handleConn(serverValues ServerValues, conn net.Conn) {
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

		_, err = conn.Write([]byte(handleCommand(line) + "\n"))
		if err != nil {
			fmt.Println("write error:", err)
			return
		}
	}
}

// case-sensitive, for simplicity
func handleCommand(input string) string {
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
	default:
		return "ERR unknown command"
	}
}
