package worker

import (
	"net"
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
	for {
		writer := json.NewEncoder(worker.conn)
	
		err := writer.Encode(<- worker.In)
		if err != nil {
			return err
		}
		
		reader := json.NewDecoder(worker.conn)
		
		var response lib.DataTransfer

		err = reader.Decode(&response)
		if err != nil {
			return err
		}

		worker.Out <- response
	}
}

// Stop stops worker running tasks.
func (worker *Worker) Stop() error {
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
