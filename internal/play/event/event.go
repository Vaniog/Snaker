package event

import (
	"encoding/json"
)

type Type string

type Event interface{}

func Parse[E Event](data []byte) (e E, ok bool) {
	err := json.Unmarshal(data, &e)
	ok = true
	if err != nil {
		ok = false
	}
	return
}

func ParseType(data []byte) Type {
	var e struct {
		Type Type `json:"event"`
	}
	err := json.Unmarshal(data, &e)
	if err != nil {
		return ""
	}
	return e.Type
}
