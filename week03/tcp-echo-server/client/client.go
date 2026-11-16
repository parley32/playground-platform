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

	responseIsExpected := make(chan bool)
	response := make(chan string)

	go func() {
		fmt.Printf("> ")
		for {
			text, err := stdinReader.ReadString('\n')
			if err != nil {
				fmt.Println("input reading failed: ", err)
				return
			}
			responseIsExpected <- true
			_, err = conn.Write([]byte(text))
			if err != nil {
				fmt.Println("failed to send text to server: ", err)
				return
			}
		}
	}()

	go func() {
		defer close(response)
		for {
			reply, err := connReader.ReadString('\n')
			if err != nil {
				fmt.Println("reply error: ", err)
				return
			}
			response <- reply
		}
	}()

	for {
		select {
		case <-responseIsExpected:
			reply, ok := <-response
			if !ok {
				return
			}
			fmt.Print("\r\033[K" + reply)
		case reply, ok := <-response:
			fmt.Print("\r\033[K" + reply)
			if !ok {
				return
			}
		}
		fmt.Printf("> ")
	}

}
