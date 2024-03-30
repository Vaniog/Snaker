package game

import "time"

type Options struct {
	FPS           int
	FrameDuration time.Duration
	SnakeLen      int
	Duration      time.Duration
	Field         Field
}

var DefaultOptions = NewOptions(
	10,
	5,
	Field{W: 40, H: 40},
	10*time.Minute,
)

func NewOptions(fps, snakeLen int, field Field, dur time.Duration) Options {
	return Options{
		fps,
		time.Second / time.Duration(fps),
		snakeLen,
		dur,
		field,
	}
}
