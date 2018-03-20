package main

import (
	"fmt"
	"net"
	"bufio"
	"../lib"
)

func main() {
	connectToMaster()
}

func connectToMaster() {
	conn, err := net.Dial("tcp", ":25000")

	if err != nil {
		fmt.Println(err.Error())
		conn.Close()
		return
	}
	for {
		handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	for {
		text, err := bufio.NewReader(conn).ReadString('\n')

		if err != nil {
			fmt.Println(err.Error())
			conn.Close()
			return
		}

		fmt.Println(text)
	}
}
