package manager

import (
	"github.com/clitetailor/distributed-hash-decrypter/master/worker"
	"github.com/clitetailor/distributed-hash-decrypter/lib/charset"
	"github.com/clitetailor/distributed-hash-decrypter/lib"
	"net"
	"log"
	"fmt"
)

// Manager manages worker clusters.
type Manager struct {
	ln net.Listener
	workers map[int]worker.Worker

	In chan string
	Out chan string
}

// New initializes and returns a new Manager.
func New(ln net.Listener) Manager {
	return Manager {
		ln: ln,
		workers: make(map[int]worker.Worker),
		
		In: make(chan string),
		Out: make(chan string) }
}

// Run runs manager tasks.
func (manager *Manager) Run() {
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
func (manager *Manager) Deliver() {
	for {
		request := <- manager.In

		start := []rune("a")
		end := []rune("999")
		
		if len(manager.workers) == 0 {
			manager.Out <- "No workers found!"
			continue
		}

		ranges := charset.Range(start, end, len(manager.workers))

		i := 0
		for _, w := range manager.workers {
			w.In <- lib.DataTransfer {
				Type: "data",
				Start: ranges[i][0],
				End: ranges[i][1],
				Code: request }

			i++
			go func(w worker.Worker) {
				select {
				case <- w.IsStopped:
					return

				case response := <- w.Out:
					switch response.Type {
					case "found":
						manager.BroadcastStop()
						manager.Out <- response.Result
						
					case "notfound":
						if manager.Done() {
							manager.Out <- "Not found!"
						}
					}
				}
			}(w)
		}
	}
}

// Done checks whether all workers have finished.
func (manager *Manager) Done() bool {
	for _, w := range manager.workers {
		if !w.Done {
			return false
		}
	}

	return true
}

// Add adds new worker connection to manager.
func (manager *Manager) Add(conn net.Conn) {
	workers := manager.workers
	id := len(workers)

	worker := worker.New(conn)
	manager.workers[id] = worker

	fmt.Println("Conns: ", len(workers))

	go func() {
		err := worker.Run()
		if err != nil {
			log.Output(2, err.Error())

			worker.Destroy()
			delete(workers, id)

			fmt.Println("Conns: ", len(workers))
		}
	}()
}

// BroadcastStop sends stop signal to all working workers.
func (manager *Manager) BroadcastStop() {
	for _, worker := range manager.workers {
		if !worker.Done {
			worker.Stop()
		}
	}
}
