package manager

import (
	"github.com/clitetailor/gohashgodistributed/lib"
	"github.com/clitetailor/gohashgodistributed/lib/charset"
	"github.com/clitetailor/gohashgodistributed/master/worker"
	"log"
	"net"
)

var nodeID = 0

// Manager manages worker clusters.
type Manager struct {
	ln      net.Listener
	workers map[*worker.Worker]*worker.Worker
	In      chan string
	Out     chan string
}

// NewManager initializes and returns new Manager.
func NewManager(ln net.Listener) *Manager {
	return &Manager{
		ln:      ln,
		workers: make(map[*worker.Worker]*worker.Worker),
		In:      make(chan string),
		Out:     make(chan string)}
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
		request := <-manager.In

		start := []rune("a")
		end := []rune("999999")

		if len(manager.workers) == 0 {
			manager.Out <- "No workers found!"
			continue
		}

		ranges := charset.Range(start, end, len(manager.workers))

		i := 0
		for _, w := range manager.workers {
			w.In <- lib.DataTransfer{
				Type:  "data",
				Start: ranges[i][0],
				End:   ranges[i][1],
				Code:  request}

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

	w := worker.NewWorker(conn)
	workers[w] = w

	kill := make(chan bool)
	go func() {
		nodeID++
		log.Println("Add Node:", nodeID)
		log.Println("Conns:", len(workers))

		err := w.Run()

		log.Println("Remove Node:", nodeID)

		if err != nil {
			log.Println(err)

			w.Destroy()
			delete(workers, w)
		}

		kill <- true
		log.Println("Conns:", len(workers))
	}()

	go func() {
		for {
			select {
			case response := <-w.Out:
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

			case <-kill:
				return
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
