package main

import (
	"fmt"
	"net"
	"bufio"
	"encoding/json"
	"log"
	"../lib"
	"../lib/charset"
)

func main() {
	conn, err := net.Dial("tcp", ":25000")

	if err != nil {
		fmt.Println(err.Error())
		conn.Close()
		return
	}

	for {
		handleConnection(conn)
	}

	worker := New()
	worker.Init()
}

// Worker stores informations about worker.
type Worker struct {
	conn net.Conn
	nRoutines int
	stop bool
}

// New initializes and returns a new Worker.
func New(conn net.Conn) Worker {
	return Worker {
		conn: conn,
		nRoutines: 3
	}
}

// Init handles connection IO and run tasks.
func (worker Worker) Init() {
	reader := json.NewDecoder(worker.conn)

	for {
		data := new(DataTransfer)
		err := reader.Decode(&data)

		if err != nil {
			log.Output(1, err.Error())
			worker.conn.Close()
			return
		}

		worker.Run()
	}
}

// Run runs worker task.
func (worker Worker) Run(data DataTransfer) {
	writer := json.NewEncoder(worker.conn)

	switch data.Type {
		case "data": {
			worker.stop = false		
			
			worker.RunHash(data)
		}

		case "exit": {
			worker.Stop()
		}
	}
}

func (worker Worker) StartGoroutines(data DataTransfer) {
	for i := 0; i < worker.nRoutines; i++ {
		go worker.RunHash(data)
	}
}

// RunHash runs hash task.
func (worker Worker) RunHash(data DataTransfer) {
	for i := data.Start; charset.Sign(i, data.End) < 0; charset.IncRuneArr(i) {
		if worker.stop {
			return
		}

		if charset.IsValid(i) {
			str := string(i)
			code := charset.HashString(i)

			if code == data.Code {
				data := DataTransfer {
					type: "found",
					result: str
				}
				writer.Encode(data)
				worker.stop = true
				return
			}
		}
	}

	data := DataTransfer {
		type: "notfound"
	}
	writer.Encode(data)
	worker.stop = true
}

// Stop signal worker to stop other tasks.
func (worker Worker) Stop(data DataTransfer) {
	worker.stop = true
}
