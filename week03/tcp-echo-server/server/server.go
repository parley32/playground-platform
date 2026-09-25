package main

import (
	"bufio"
	"fmt"
	"net"
)

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
		go handleConn(conn)
	}
}

func handleConn(conn net.Conn) {
	defer conn.Close()
	fmt.Println("client connected:", conn.RemoteAddr())

	reader := bufio.NewReader(conn)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("client disconnected:", conn.RemoteAddr())
			return
		}

		_, err = conn.Write([]byte("echo: " + line))
		if err != nil {
			fmt.Println("write error:", err)
			return
		}
	}
}
