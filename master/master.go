package main

import (
	"fmt"
	"net"
	"bufio"
)

func main() {
	createServer()
}

func createServer() {
	ln, err := net.Listen("tcp", ":25000")
	ln2, err2 := net.Listen("tcp", ":8080")

	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println("Listening on ports 25000")

	if err2 != nil {
		fmt.Println(err2.Error())
		return
	}
	fmt.Println("Listening  on port 8080")

	in := make(chan string)
	out := make(chan string)

	go func() {
		for {
			conn, err := ln2.Accept()

			fmt.Println("Connected to client")
			
			if err != nil {
				fmt.Println(err.Error())
				conn.Close()
				continue
			}

			go handleClientConnection(conn, in, out)
		}
	}()

	for {
		conn, err := ln.Accept()
	
		if err != nil {
			fmt.Println(err)
			conn.Close()
			return
		}

		go handleWorkerConnection(conn, in, out)
	}
}

func handleClientConnection(conn net.Conn, in chan string, out chan string) {
	for {
		reader := bufio.NewReader(conn)

		for {
			request, err := reader.ReadString('\n')
			if err != nil {
				fmt.Println(err)
				conn.Close()
				return
			}

			fmt.Print(request)
			fmt.Fprintf(conn, request)
		}
	}
}

func handleWorkerConnection(conn net.Conn, in chan string, out chan string) {
	reader := bufio.NewReader(conn)

	for {
		hash, err := reader.ReadString('\n')

		if err != nil {
			fmt.Println(err.Error())
			conn.Close()
			return
		}

		in <- hash
		response := <- out

		fmt.Printf(response, conn)
	}
}