package main

import (
	"fmt"
	"net"
	"bufio"
)

func main() {
	connectToMaster()
}

func connectToMaster() {
	conn, err := net.Dial("tcp", ":25000")
	if err != nil {
		handleError(err)
	}
	for {
		handleConnection(conn)
	}
}

func handleError(err error) {
	fmt.Print(err.Error())
}

func handleConnection(conn net.Conn) {
	text, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		handleError(err)
	}
	fmt.Print(text)
}