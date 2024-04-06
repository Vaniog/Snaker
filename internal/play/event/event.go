package event

import (
	"encoding/json"
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

func Parse[E Event](data []byte) (E, error) {
	var e E
	err := json.Unmarshal(data, &e)
	return e, err
}
