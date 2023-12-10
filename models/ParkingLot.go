package models

import (
	"sync"
)

type ParkingLot struct {
	spots         []*Place
	queueCars     *Queue
	mu            sync.Mutex
	availableCond *sync.Cond
}

func NewParkingLot(spots []*Place) *ParkingLot {
	queue := NewQueue()

	p := &ParkingLot{
		spots:     spots,
		queueCars: queue,
	}
	p.availableCond = sync.NewCond(&p.mu)

	return p
}

func (p *ParkingLot) GetPlaces() []*Place {
	return p.spots
}

func (p *ParkingLot) GetPlaceAvailable() *Place {
	p.mu.Lock()
	defer p.mu.Unlock()

	for {
		for _, spot := range p.spots {
			if spot.GetIsAvailable() {
				spot.SetIsAvailable(false)
				return spot
			}
		}
		p.availableCond.Wait()
	}
}

func (p *ParkingLot) SetFreePlace(spot *Place) {
	p.mu.Lock()
	defer p.mu.Unlock()

	spot.SetIsAvailable(true)
	p.availableCond.Signal()
}

func (p *ParkingLot) GetQueue() *Queue {
	return p.queueCars
}
