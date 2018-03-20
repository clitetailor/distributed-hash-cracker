package manager

import (
	"../worker"
)

// Manager manages worker clusters.
type Manager struct {
	workers []Worker
}

// New initializes and returns a new Manager.
func New() Manager {
	manager := make(Manager)
	manager.workers = []Worker{}
}
