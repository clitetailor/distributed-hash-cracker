package manager

import (
	"github.com/clitetailor/distributed-hash-decrypter/master/worker"
	"github.com/clitetailor/distributed-hash-decrypter/lib/charset"
	"github.com/clitetailor/distributed-hash-decrypter/lib"
	"net"
	"log"
)

// Manager manages worker clusters.
type Manager struct {
	ln net.Listener
	workers map[int]worker.Worker

	In chan string
	Out chan string
	Done chan bool
}

// New initializes and returns a new Manager.
func New(ln net.Listener) Manager {
	return Manager {
		ln: ln,
		workers: make(map[int]worker.Worker),
		
		In: make(chan string),
		Out: make(chan string),
		Done: make(chan bool) }
}

// Run runs manager tasks.
func (manager Manager) Run() {
	go func() {
		manager.Deliver()
	}()
	
	for {
		conn, err := manager.ln.Accept()

		if err != nil {
			log.Output(1, err.Error())
			conn.Close()
			return
		}

		manager.Add(conn)
	}
}

// Deliver receives data from client and delivers to workers.
func (manager Manager) Deliver() {
	for {
		request := <- manager.In
		nWorker := len(manager.workers)

		start := []rune("a")
		end := []rune("999")
		
		if nWorker == 0 {
			ranges := charset.Range(start, end, nWorker)

			i := 0
			for _, worker := range manager.workers {
				worker.In <- lib.DataTransfer {
					Type: "data",
					Start: ranges[i][0],
					End: ranges[i][1],
					Code: request }

				i++

				go func() {
					select {
					case response := <- worker.Out:
						manager.Out <- response
						manager.Stop()
					case <- worker.StopSignal:
					}
				}()
			}
		} else {
			manager.Out <- "No workers found!"
		}
	}
}

// Add adds new worker connection to manager.
func (manager Manager) Add(conn net.Conn) {
	id := len(manager.workers)

	worker := worker.New(conn)
	manager.workers[id] = worker

	go func() {
		worker.Run()
		delete(manager.workers, id)
	}()
}

// Stop sends stop signal to all workers.
func (manager Manager) Stop() {
	for _, worker := range manager.workers {
		worker.Stop()
	}
}
