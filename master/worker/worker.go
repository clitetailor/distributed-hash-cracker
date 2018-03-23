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
}

// NewWorker initializes and returns a new Worker.
func NewWorker(conn net.Conn) *Worker {
	return &Worker {
		conn: conn,
		Done: false,
		In: make(chan lib.DataTransfer),
		Out: make(chan lib.DataTransfer) }
}

// Run runs and manager the connection to worker.
func (worker *Worker) Run() error {
	writer := json.NewEncoder(worker.conn)
	reader := json.NewDecoder(worker.conn)

	kill := make(chan error)

	go func() {
		for {
			var response lib.DataTransfer

			err := reader.Decode(&response)
			if err != nil {
				log.Println(err)
				kill <- err
				return
			}

			worker.Out <- response
		}
	}()

	for {
		select {
		case response := <- worker.In:
			err := writer.Encode(response)
			if err != nil {
				return err
			}

		case err := <- kill:
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
}
