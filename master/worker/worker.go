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
	Done bool
	In chan lib.DataTransfer
	Out chan lib.DataTransfer
	IsStopped chan bool
}

// New initializes and returns a new Worker.
func New(conn net.Conn) Worker {
	return Worker {
		conn: conn,
		Done: false,
		In: make(chan lib.DataTransfer),
		Out: make(chan lib.DataTransfer),
		IsStopped: make(chan bool) }
}

// Run runs and manager the connection to worker.
func (worker *Worker) Run() error {
	writer := json.NewEncoder(worker.conn)
	reader := json.NewDecoder(worker.conn)

	go func() {
		for {
			var response lib.DataTransfer

			err2 := reader.Decode(&response)
			if err2 != nil {
				log.Println(err2)
				return
			}

			worker.Out <- response
		}
	}()

	for {
		err := writer.Encode(<- worker.In)
		if err != nil {
			return err
		}
	}
}

// SendStop stops worker running tasks.
func (worker *Worker) SendStop() error {
	data := lib.DataTransfer {
		Type: "stop" }

	err := json.NewEncoder(worker.conn).Encode(data)
	if err != nil {
		return err
	}

	return nil
}

// Destroy closes connection to worker and kills all channels.
func (worker *Worker) Destroy() {
	worker.conn.Close()
	close(worker.In)
	close(worker.Out)
	worker.IsStopped <- true
}
