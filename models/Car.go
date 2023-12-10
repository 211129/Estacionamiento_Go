package models

import (
	"github.com/oakmound/oak/v4/render"
	"github.com/oakmound/oak/v4/render/mod"
	"sync"
	"time"

	"github.com/oakmound/oak/v4/alg/floatgeom"
	"github.com/oakmound/oak/v4/entities"
	"github.com/oakmound/oak/v4/scene"
)

const (
	entranceSpotX = 355.00
	speed         = 10
)

type Car struct {
	area   floatgeom.Rect2
	entity *entities.Entity
	mu     sync.Mutex
}

func NewCar(ctx *scene.Context) *Car {
	area := floatgeom.NewRect2(300, 490, 320, 510)

	sprite, _ := render.LoadSprite("assets/images/car.png")

	newSwitch := render.NewSwitch("up", map[string]render.Modifiable{
		"up":    sprite,
		"down":  sprite.Copy().Modify(mod.FlipY),
		"left":  sprite.Copy().Modify(mod.Rotate(90)),
		"right": sprite.Copy().Modify(mod.Rotate(-90)),
	})

	entity := entities.New(ctx, entities.WithRect(area), entities.WithRenderable(newSwitch), entities.WithDrawLayers([]int{1, 2}))

	return &Car{
		area:   area,
		entity: entity,
	}
}

func (c *Car) move(direction string, point float64, manager *Orchestrator) {
	if direction == "left" {
		for c.X() > point {
			if !c.isACar("left", manager.GetCars()) {
				c.entity.Renderable.(*render.Switch).Set("left")
				c.ShiftX(-1)
				time.Sleep(speed * time.Millisecond)
			}
		}
	} else if direction == "right" {
		for c.X() < point {
			if !c.isACar("right", manager.GetCars()) {
				c.entity.Renderable.(*render.Switch).Set("right")
				c.ShiftX(1)
				time.Sleep(speed * time.Millisecond)
			}
		}
	} else if direction == "up" {
		for c.Y() > point {
			if !c.isACar("up", manager.GetCars()) {
				c.entity.Renderable.(*render.Switch).Set("up")
				c.ShiftY(-1)
				time.Sleep(speed * time.Millisecond)
			}
		}
	} else if direction == "down" {
		for c.Y() < point {
			if !c.isACar("down", manager.GetCars()) {
				c.entity.Renderable.(*render.Switch).Set("down")
				c.ShiftY(1)
				time.Sleep(speed * time.Millisecond)
			}
		}
	}
}

func (c *Car) Enqueue(manager *Orchestrator) {
	c.move("up", 45, manager)
}

func (c *Car) JoinDoor(manager *Orchestrator) {
	c.move("right", entranceSpotX, manager)
}

func (c *Car) ExitDoor(manager *Orchestrator) {
	c.move("left", 300, manager)
}

func (c *Car) Park(spot *Place, manager *Orchestrator) {
	for index := 0; index < len(*spot.GetInstructionsForParking()); index++ {
		directions := *spot.GetInstructionsForParking()
		if directions[index].Direction == "right" {
			c.move("right", directions[index].Point, manager)
		} else if directions[index].Direction == "down" {
			c.move("down", directions[index].Point, manager)
		}
	}
}

func (c *Car) Leave(spot *Place, manager *Orchestrator) {
	for index := 0; index < len(*spot.GetInstructionsForLeaving()); index++ {
		directions := *spot.GetInstructionsForLeaving()
		if directions[index].Direction == "left" {
			c.move("left", directions[index].Point, manager)
		} else if directions[index].Direction == "right" {
			c.move("right", directions[index].Point, manager)
		} else if directions[index].Direction == "up" {
			c.move("up", directions[index].Point, manager)
		} else if directions[index].Direction == "down" {
			c.move("down", directions[index].Point, manager)
		}
	}
}

func (c *Car) LeavePlace(manager *Orchestrator) {
	spotY := c.Y()
	c.move("up", spotY-30, manager)
}

func (c *Car) GoLong(manager *Orchestrator) {
	c.move("left", -20, manager)
}

func (c *Car) ShiftY(dy float64) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.entity.ShiftY(dy)
}

func (c *Car) ShiftX(dx float64) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.entity.ShiftX(dx)
}

func (c *Car) X() float64 {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.entity.X()
}

func (c *Car) Y() float64 {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.entity.Y()
}

func (c *Car) Remove() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.entity.Destroy()
}

func (c *Car) isACar(direction string, cars []*Car) bool {
	minDistance := 30.0
	for _, car := range cars {
		if direction == "left" {
			if c.X() > car.X() && c.X()-car.X() < minDistance && c.Y() == car.Y() {
				return true
			}
		} else if direction == "right" {
			if c.X() < car.X() && car.X()-c.X() < minDistance && c.Y() == car.Y() {
				return true
			}
		} else if direction == "up" {
			if c.Y() > car.Y() && c.Y()-car.Y() < minDistance && c.X() == car.X() {
				return true
			}
		} else if direction == "down" {
			if c.Y() < car.Y() && car.Y()-c.Y() < minDistance && c.X() == car.X() {
				return true
			}
		}
	}
	return false
}
