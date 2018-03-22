package main

import (
	"bufio"
	"fmt"
	"flag"
	"net"
	"log"
)

func main() {
	host := flag.String("host", "", "node host")
	port := flag.Int("port", 8080, "node port")
	code := flag.String("code", "", "md5sum")
	flag.Parse()

	address := GetAddress(*host, *port)
	conn, err := net.Dial("tcp", address)

	if err != nil {
		HandleError(err)
		return
	}
	SendMessage(conn, *code)
}

// GetAddress returns address base on host and port.
func GetAddress(host string, port int) (string) {
	return fmt.Sprintf("%s:%d", host, port)
}

// SendMessage sends message to connection.
func SendMessage(conn net.Conn, code string) {
	_, err := fmt.Fprintf(conn, code + "\n")

	if err != nil {
		HandleError(err)
	}

	reader := bufio.NewReader(conn)
	response, err2 := reader.ReadString('\n')
	if err2 != nil {
		HandleError(err2)
	}

	fmt.Println(response)
}

// HandleError logs the error and exits.
func HandleError(err error) {
	log.Fatal(err)
}
