package event

import (
	"encoding/json"
)

type Type string

type Event struct {
	Type `json:"event"`
}

func Parse[E any](data []byte) (e E, ok bool) {
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

func Bytes(e any) []byte {
	data, _ := json.Marshal(e)
	return data
}
