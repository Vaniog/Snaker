package game

import "time"

type Options struct {
	FPS           int
	FrameDuration time.Duration
	SnakeLen      int
	SnakesAmount  int
	Duration      time.Duration
	Field         Field
}

const defaultFps = 10

var DefaultOptions = Options{
	FPS:           defaultFps,
	FrameDuration: time.Second / defaultFps,
	SnakeLen:      5,
	SnakesAmount:  2,
	Duration:      time.Minute * 10,
	Field: Field{
		W: 40,
		H: 40,
	},
}
