package worker

import (
	"net"
	"log"
	"encoding/json"
	"github.com/clitetailor/distributed-hash-decrypter/lib"
)

// Worker stores information about worker cluster.
type Worker struct {
	conn net.Conn
	In chan lib.DataTransfer
	Out chan string
	Done chan bool
	StopSignal chan bool
}

// New initializes and returns a new Worker.
func New(conn net.Conn) Worker {
	return Worker {
		conn: conn,
		In: make(chan lib.DataTransfer),
		Out: make(chan string),
		StopSignal: make(chan bool) }
}

// Run runs and manager the connection to worker.
func (worker Worker) Run() {
	for {
		writer := json.NewEncoder(worker.conn)
	
		err := writer.Encode(<- worker.In)
		if err != nil {
			log.Output(1, err.Error())
			worker.conn.Close()
			return
		}
		
		reader := json.NewDecoder(worker.conn)
		
		var response lib.DataTransfer

		err2 := reader.Decode(&response)
		if err2 != nil {
			log.Output(1, err2.Error())
			worker.conn.Close()
			return
		}

		switch response.Type {
		case "found":
			worker.Out <- response.Result
			
		case "notfound":
			worker.Done <- true
		}
	}
}

// Stop stops worker running tasks.
func (worker Worker) Stop() {
	data := lib.DataTransfer {
		Type: "stop" }

	worker.StopSignal <- true

	err := json.NewEncoder(worker.conn).Encode(data)
	if err != nil {
		log.Output(1, err.Error())
		worker.conn.Close()
	}
}
