package models

import "sync"

type Orchestrator struct {
	Cars  []*Car
	Mutex sync.Mutex
}

func NewOrchestrator() *Orchestrator {
	return &Orchestrator{
		Cars: make([]*Car, 0),
	}
}

func (cm *Orchestrator) Add(car *Car) {
	cm.Mutex.Lock()
	defer cm.Mutex.Unlock()
	cm.Cars = append(cm.Cars, car)
}

func (cm *Orchestrator) Remove(car *Car) {
	cm.Mutex.Lock()
	defer cm.Mutex.Unlock()
	for i, c := range cm.Cars {
		if c == car {
			cm.Cars = append(cm.Cars[:i], cm.Cars[i+1:]...)
			break
		}
	}
}

func (cm *Orchestrator) GetCars() []*Car {
	cm.Mutex.Lock()
	defer cm.Mutex.Unlock()
	return cm.Cars
}
