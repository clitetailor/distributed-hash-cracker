package worker

import (
	"bufio"
	"net"
	"fmt"
	"log"
)

// Worker stores information about worker cluster.
type Worker struct {
	id int
	conn net.Conn
	in chan string
	out chan string
	exit chan bool
	dis chan bool
}

// New initializes and returns a new Worker.
func New(id int, conn net.Conn) Worker {
	return Worker {
		id: id,
		conn: conn,
		in: make(chan string),
		out: make(chan string),
		exit: make(chan bool),
		dis: make(chan bool) }
}

// GetDis returns chan that signals when worker is disconnected.
func (worker Worker) GetDis() (chan bool) {
	return worker.dis
}

// Run runs and manager the connection to worker.
func (worker Worker) Run() {
	reader := bufio.NewReader(worker.conn)

	for {
		_, err := fmt.Fprintf(worker.conn, <- worker.in)

		if err != nil {
			log.Output(1, err.Error())
			worker.conn.Close()
			worker.dis <- true
			return
		}

		response, err2 := reader.ReadString('\n')
		worker.out <- response
		
		if err2 != nil {
			log.Output(1, err2.Error())
			worker.CloseConn()
			return
		}
	}
}

// CloseConn closes the connection to worker.
func (worker Worker) CloseConn() {
	worker.conn.Close()
	worker.dis <- true
}
