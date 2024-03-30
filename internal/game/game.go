package game

import (
	"context"
	"encoding/json"
	"math/rand"
	"time"
)

type Game struct {
	Field  Field    `json:"field"`
	Snakes []*Snake `json:"snakes"`
	Food   Food     `json:"food"`
	Opts   Options  `json:"options"`

	Start time.Time `json:"start"`
}

func NewGame(opts Options) *Game {
	g := Game{Opts: opts}
	g.Field = opts.Field
	return &g
}

func (g *Game) RegisterSnake() *Snake {
	s := NewSnake(g, Point{0, 0}, Up, g.Opts.SnakeLen)

	g.Snakes = append(g.Snakes, s)
	return s
}

func (g *Game) MoveSnakes() {
	for _, s := range g.Snakes {
		if !s.Alive {
			continue
		}
		s.Move()
		g.Field.HandleCollision(s)
		g.Food.HandleCollision(s)
		for _, enemy := range g.Snakes {
			if enemy == s {
				continue
			}
			enemy.HandleCollision(s)
		}
	}
}

func (g *Game) Update() {
	g.MoveSnakes()
}

func (g *Game) ShuffleSnakes() {
	//TODO make better
	for _, s := range g.Snakes {
		s.SetHead(Point{
			rand.Intn(g.Field.W),
			rand.Intn(g.Field.H),
		})
	}
}

func (g *Game) IsEnd() bool {
	return time.Now().After(g.Start.Add(g.Opts.Duration))
}

func (g *Game) Run(ctx context.Context, callback func()) {
	g.ShuffleSnakes()
	g.ReloadFood()
	g.Start = time.Now()

	ticker := time.NewTicker(g.Opts.FrameDuration)
	defer ticker.Stop()
	for !g.IsEnd() {
		select {
		case <-ticker.C:
			g.Update()
			callback()
		case <-ctx.Done():
			return
		}
	}
}

func (g *Game) JSON() ([]byte, error) {
	return json.Marshal(g)
}

func (g *Game) ReloadFood() {
	g.Food.Point = Point{
		rand.Intn(g.Field.W),
		rand.Intn(g.Field.H),
	}
}