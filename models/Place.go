package models

import (
	"github.com/oakmound/oak/v4/alg/floatgeom"
)

type Instruction struct {
	Direction string
	Point     float64
}

func newInstruction(direction string, point float64) *Instruction {
	return &Instruction{
		Direction: direction,
		Point:     point,
	}
}

type Place struct {
	area                 *floatgeom.Rect2
	directionsForParking *[]Instruction
	directionsForLeaving *[]Instruction
	number               int
	isAvailable          bool
}

func NewPlace(x, y, x2, y2 float64, row, number int) *Place {
	directionsForParking := getInstructionsForParking(x, y, row)
	directionsForLeaving := getInstructionsForLeaving()
	area := floatgeom.NewRect2(x, y, x2, y2)

	return &Place{
		area:                 &area,
		directionsForParking: directionsForParking,
		directionsForLeaving: directionsForLeaving,
		number:               number,
		isAvailable:          true,
	}
}

func getInstructionsForParking(x, y float64, row int) *[]Instruction {
	var directions []Instruction

	if row == 1 {
		directions = append(directions, *newInstruction("down", 45))
	} else if row == 2 {
		directions = append(directions, *newInstruction("down", 135))
	} else if row == 3 {
		directions = append(directions, *newInstruction("down", 225))
	} else if row == 4 {
		directions = append(directions, *newInstruction("down", 315))
	}

	directions = append(directions, *newInstruction("right", x+5))
	directions = append(directions, *newInstruction("down", y+5))

	return &directions
}

func getInstructionsForLeaving() *[]Instruction {
	var directions []Instruction

	directions = append(directions, *newInstruction("right", 600))
	directions = append(directions, *newInstruction("up", 15))
	directions = append(directions, *newInstruction("left", 355))

	return &directions
}

func (p *Place) GetArea() *floatgeom.Rect2 {
	return p.area
}

func (p *Place) GetNumber() int {
	return p.number
}

func (p *Place) GetInstructionsForParking() *[]Instruction {
	return p.directionsForParking
}

func (p *Place) GetInstructionsForLeaving() *[]Instruction {
	return p.directionsForLeaving
}

func (p *Place) GetIsAvailable() bool {

	return p.isAvailable
}

func (p *Place) SetIsAvailable(isAvailable bool) {
	p.isAvailable = isAvailable
}
