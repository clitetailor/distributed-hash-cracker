package main

import (
	"github.com/clitetailor/distributed-hash-decrypter/master/manager"
	"fmt"
	"net"
	"bufio"
	"log"
)

func main() {
	ln, err := net.Listen("tcp", ":8080")
	ln2, err2 := net.Listen("tcp", ":25000")

	if err != nil {
		log.Fatal(err)
	}
	
	if err2 != nil {
		log.Fatal(err)
	}

	m := manager.NewManager(ln2)

	go func() {
		m.Run()
	}()
	
	for {
		conn, err := ln.Accept()
		
		if err != nil {
			log.Println(err)
			conn.Close()
			continue
		}
		
		go HandleConnection(conn, m.In, m.Out)
	}
}

// HandleConnection handles connection input and output.
func HandleConnection(conn net.Conn, in chan string, out chan string) {
	reader := bufio.NewReader(conn)

	for {
		request, err := reader.ReadString('\n')

		if err != nil {
			log.Println(err)
			conn.Close()
			return
		}

		go func() {
			in <- request
			response := <- out
			
			_, err2 := fmt.Fprintf(conn, response + "\n")
			if err2 != nil {
				log.Println(err2)
				conn.Close()
				return
			}
		}()
	}
}
