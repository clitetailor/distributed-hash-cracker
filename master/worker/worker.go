package worker

import (
	"bufio"
	"net"
	"fmt"
	"log"
)

// Worker stores information about worker cluster.
type Worker struct {
	conn net.Conn
	in chan string
	out chan string
	exit chan bool
	dis chan bool
}

// New initializes and returns a new Worker.
func New(conn net.Conn) Worker {
	worker := new(Worker)
	worker.conn = conn

	worker.in = make(chan string)
	worker.out = make(chan string)
	worker.exit = make(chan bool)
	worker.dis = make(chan bool)
	
	return *worker
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
			log.Println(err)
			worker.conn.Close()
			worker.dis <- true
			return
		}

		response, err2 := reader.ReadString('\n')
		worker.out <- response
		
		if err2 != nil {
			fmt.Println(err2.Error())
			worker.CloseConn()
			return
		}
	}
}

func (worker Worker) CloseConn() {
	worker.conn.Close()
	worker.dis <- true
}
