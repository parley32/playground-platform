package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("failed to connect to server:", err)
		return
	}
	defer conn.Close()

	fmt.Println("Connected to server, type string and Enter")

	stdinReader := bufio.NewReader(os.Stdin)
	connReader := bufio.NewReader(conn)

	for {
		fmt.Printf("> ")
		text, err := stdinReader.ReadString('\n')
		if err != nil {
			fmt.Println("input reading failed: ", err)
			return
		}
		_, err = conn.Write([]byte(text))
		if err != nil {
			fmt.Println("failed to send text to server: ", err)
			return
		}

		reply, err := connReader.ReadString('\n')
		if err != nil {
			fmt.Println("reply error: ", err)
			return
		}
		fmt.Print(reply)
	}
}
