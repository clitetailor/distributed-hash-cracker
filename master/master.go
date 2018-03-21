package main

import (
	"github.com/clitetailor/distributed-hash-decrypter/lib"
	"github.com/clitetailor/distributed-hash-decrypter/master/manager"
	"github.com/clitetailor/distributed-hash-decrypter/master/worker"
	"fmt"
	"net"
	"bufio"
	"log"
)

func main() {
	ln, err := net.Listen("tcp", ":8080")
	ln2, err2 := net.Listen("tcp", ":25000")

	if err != nil {
		log.Output(1, err.Error())
		return
	}
	fmt.Println("Listening  on port 8080")
	
	if err2 != nil {
		log.Output(1, err2.Error())
		return
	}
	fmt.Println("Listening on ports 25000")

	m := manager.New(ln2)

	go func() {
		m.Run()
	}()
	
	for {
		conn, err := ln.Accept()
		
		if err != nil {
			log.Output(1, err.Error())
			conn.Close()
			continue
		}
		
		go handleConnection(conn, m.In, m.Out)
	}
}

func handleConnection(conn net.Conn, in chan string, out chan string) {
	for {
		reader := bufio.NewReader(conn)

		for {
			request, err := reader.ReadString('\n')
			if err != nil {
				log.Output(1, err)
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
