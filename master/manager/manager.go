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
	workers map[int]*worker.Worker
	In chan string
	Out chan string
}

// NewManager initializes and returns new Manager.
func NewManager(ln net.Listener) *Manager {
	return &Manager {
		ln: ln,
		workers: make(map[int]*worker.Worker),
		In: make(chan string),
		Out: make(chan string) }
}

// Run runs manager tasks.
func (manager *Manager) Run() {
	go manager.Deliver()
	
	for {
		conn, err := manager.ln.Accept()
		if err != nil {
			log.Println(err)
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
		end := []rune("999999")

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

	w := worker.New(conn)
	workers[id] = &w

	log.Println("Conns:", len(workers))

	go func() {
		err := w.Run()

		if err != nil {
			log.Println(err)

			w.Destroy()
			delete(workers, id)

			log.Println("Conns:", len(workers))
		}
	}()


	go func() {
		for {
			select {
			case <- w.IsStopped:
				
			case response := <- w.Out:
				w.Done = true
				
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
		}
	}()
}

// BroadcastStop sends stop signal to all working workers.
func (manager *Manager) BroadcastStop() {
	for _, worker := range manager.workers {
		if worker.Done != true {
			worker.SendStop()
		}
	}
}
