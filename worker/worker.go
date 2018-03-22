package main

import (
	"fmt"
	"net"
	"encoding/json"
	"log"
	"sync"
	"strings"
	"github.com/clitetailor/distributed-hash-decrypter/lib"
	"github.com/clitetailor/distributed-hash-decrypter/lib/charset"
)

func main() {
	conn, err := net.Dial("tcp", ":25000")

	if err != nil {
		fmt.Println(err.Error())
		conn.Close()
		return
	}

	worker := New(conn)
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
		nRoutines: 3 }
}

// Init handles connection IO and run tasks.
func (worker *Worker) Init() {
	reader := json.NewDecoder(worker.conn)

	for {
		data := lib.DataTransfer{}
		err := reader.Decode(&data)

		if err != nil {
			log.Output(1, err.Error())
			worker.conn.Close()
			return
		}

		worker.Run(data)
	}
}

// Run runs worker task.
func (worker *Worker) Run(data lib.DataTransfer) {
	fmt.Println(data)
	
	switch data.Type {
	case "data":
		worker.stop = false		
		go func() {
			worker.StartGoroutines(data)
		}()

	case "stop":
		worker.Stop()
	}
}

// StartGoroutines starts worker goroutines.
func (worker *Worker) StartGoroutines(data lib.DataTransfer) {
	ranges := charset.Range(data.Start, data.End, worker.nRoutines)
	notfound := make(chan bool, 3)

	var wg sync.WaitGroup

	for i := 0; i < worker.nRoutines; i++ {
		wg.Add(1)

		go func(i int) {
			defer wg.Done()

			worker.RunHash(lib.DataTransfer {
				Start: ranges[i][0],
				End: ranges[i][1],
				
				Code: data.Code }, notfound)
		}(i)
	}

	wg.Wait()

	if len(notfound) == worker.nRoutines {
		fmt.Println("Not found!")
	
		response := lib.DataTransfer {
			Type: "notfound" }
		
		writer := json.NewEncoder(worker.conn)
		err := writer.Encode(&response)
		if err != nil {
			close(notfound)

			log.Output(1, err.Error())
			worker.conn.Close()
			return
		}
	}

	close(notfound)
}

// RunHash runs hash task.
func (worker *Worker) RunHash(data lib.DataTransfer, notfound chan bool) {
	writer := json.NewEncoder(worker.conn)

	for i := data.Start; charset.Sign(i, data.End) < 0; i = charset.IncRuneArr(i) {
		if worker.stop {
			return
		}

		if charset.IsValid(i) {
			str := string(i)
			code := charset.HashString(str)

			if strings.HasPrefix(data.Code, code) {
				fmt.Println("Found: ", str)

				data := lib.DataTransfer {
					Type: "found",
					Result: str }

				writer.Encode(data)
				worker.stop = true
				return
			}
		}
	}

	notfound <- true
}

// Stop signals workers to stop other tasks.
func (worker *Worker) Stop() {
	worker.stop = true

	fmt.Println("Stopped!")
}
