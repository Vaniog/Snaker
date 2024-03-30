package event

import (
	"github.com/Vaniog/Snaker/internal/game"
)

type Type string

const TypeRotate Type = "rotate"

type Event interface{}

type Rotate struct {
	Drc game.Direction `json:"direction"`
}

type GameStart struct {
}
