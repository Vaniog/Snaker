package game

import "time"

type Options struct {
	FPS           int
	FrameDuration time.Duration
	SnakeLen      int
	SnakesAmount  int
	Duration      time.Duration
	Field         Field
	BotsAmount    int
}

const defaultFps = 10

var DuoOptions = Options{
	FPS:           defaultFps,
	FrameDuration: time.Second / defaultFps,
	SnakeLen:      5,
	SnakesAmount:  2,
	Duration:      time.Minute * 10,
	Field: Field{
		W: 40,
		H: 40,
	},
	BotsAmount: 0,
}

var SoloOptions = Options{
	FPS:           defaultFps,
	FrameDuration: time.Second / defaultFps,
	SnakeLen:      5,
	SnakesAmount:  1,
	Duration:      time.Minute * 10,
	Field: Field{
		W: 40,
		H: 40,
	},
	BotsAmount: 1,
}
