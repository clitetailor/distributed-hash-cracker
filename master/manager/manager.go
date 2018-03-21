package manager

import (
	"net"
	"../worker"
)

// Manager manages worker clusters.
type Manager struct {
	ln net.Listener
	workers map[int]Worker

	in chan string
	out chan string
	exit chan bool
}

// New initializes and returns a new Manager.
func New(ln net.Listener) Manager {
	return {
		workers: make(map[int]Worker),
		in: make(chan string)	}
}

func (manager Manager) Run() {
	for {
		manager.conn, err := manager.ln.Accept()
	
		if err != nil {
			log.Output(1, err)
			manager.conn.Close()
			return
		}

		manager.Add(conn)
	}
}

func (manager Manager) Add(conn net.Conn) {
	id := len(manager.workers)

	worker := worker.New(id, conn)
	manager.workers[id] = worker

	go func() {
		worker.Run()
	}
}
