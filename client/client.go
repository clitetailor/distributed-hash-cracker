package main

import (
	"bufio"
	"fmt"
	"os"
	"flag"
	"net"
	"strings"
)

func main() {
	host := flag.String("host", "", "node host")
	port := flag.Int("port", 8080, "node port")
	flag.Parse()

	address := getAddress(*host, *port)
	fmt.Print(address)
	conn, err := net.Dial("tcp", ":8080")

	if err != nil {
		fmt.Print(err.Error())
	} else {
		handleConnection(conn)
	}
}

func getAddress(host string, port int) (string) {
	return fmt.Sprintf("%s:%d\n", host, port)
}

func handleConnection(conn net.Conn) {
	prompt := bufio.NewReader(os.Stdin)
	reader := bufio.NewReader(conn)

	for {
		fmt.Print("> ")
		hash, err := prompt.ReadString('\n')

		if strings.HasPrefix(hash, "exit") {
			conn.Close()
			return
		}

		if err != nil {
			fmt.Println(err.Error())
			conn.Close()
			return
		}

		fmt.Fprintf(conn, hash)

		response, err2 := reader.ReadString('\n')
		if err2 != nil {
			fmt.Println(err2.Error())
			conn.Close()
			return
		}

		fmt.Println(response)
	}
}
