package main

import (
	"fmt"
	"net"
	"bufio"
	"os"
)

func main() {
	createServer()
}

func createServer() {
	ln, err := net.Listen("tcp", ":25000")

	in := make(chan string)

	if err != nil {
		handleError(err)
	}

	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				handleError(err)
			}
			for {
				handleConnection(conn, in)
			}
		}
	}()

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("> ")
		
		text, err := reader.ReadString('\n')
		if err != nil {
			handleError(err)
		} else {
			text = text + "\n"
			fmt.Print(text)
			in <- text
		}
	}
}

func handleError(err error) {
	fmt.Print(err.Error())
}

func handleConnection(conn net.Conn, in chan string) {
	text := <- in
	conn.Write([]byte(text))
}
