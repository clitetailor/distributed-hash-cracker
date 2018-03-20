package main

import (
	"fmt"
	"net"
	"bufio"
	"../lib"
	"./manager"
	"./worker"
)

func main() {
	createServer()
}

func createServer() {
	ln, err := net.Listen("tcp", ":8080")
	ln2, err2 := net.Listen("tcp", ":25000")

	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println("Listening  on port 8080")
	
	if err2 != nil {
		fmt.Println(err2.Error())
		return
	}
	fmt.Println("Listening on ports 25000")

	in := make(chan string)
	out := make(chan string)
	exit := make(chan bool)

	go func() {
		for {
			conn, err := ln.Accept()

			fmt.Println("Connected to client")
			
			if err != nil {
				fmt.Println(err.Error())
				conn.Close()
				continue
			}

			go handleClientConnection(conn, in, out)
		}
	}()

	listenWorkers(ln2, in, out)
}

func listenWorkers(in chan string, out chan string) {
	workerConns := []net.Conn{}

	for {
		conn, err := ln2.Accept()
	
		if err != nil {
			fmt.Println(err)
			conn.Close()
			return
		}

		go handleWorkerConnection(conn, in, out, exit)
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

			in <- request
			response := <- out

			fmt.Print(response)
			fmt.Fprintf(conn, request)
		}
	}
}

func handleWorkerConnection(conn net.Conn, in chan string, out chan string, exit chan bool) {
	reader := bufio.NewReader(conn)

	for {
		_, err := fmt.Fprintf(conn, <- in)

		if err != nil {
			fmt.Println(err)
			conn.Close()

			return
		}

		response, err2 := reader.ReadString('\n')
		out <- response
		
		if err2 != nil {
			fmt.Println(err2.Error())
			conn.Close()
			return
		}
	}
}
