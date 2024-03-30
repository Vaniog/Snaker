package game

import (
	"log"
)

type Direction string

const (
	Up    Direction = "Up"
	Right Direction = "Right"
	Down  Direction = "Down"
	Left  Direction = "Left"
)

type Point struct {
	X int
	Y int
}

func (p Point) Move(d Direction) Point {
	switch d {
	case Up:
		return Point{p.X, p.Y - 1}
	case Down:
		return Point{p.X, p.Y + 1}
	case Left:
		return Point{p.X - 1, p.Y}
	case Right:
		return Point{p.X + 1, p.Y}
	}

	log.Printf("tried to move in unknown direction: %s", d)
	return p
}

func (d Direction) Opposite() Direction {
	switch d {
	case Up:
		return Down
	case Down:
		return Up
	case Left:
		return Right
	case Right:
		return Left
	}
	log.Printf("tried to move in unknown direction: %s", d)
	return d
}
