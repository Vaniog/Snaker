package game

import (
	"encoding/json"
	"math/rand"
	"time"
)

type Game struct {
	Field  Field    `json:"field"`
	Snakes []*Snake `json:"snakes"`
	Food   Food     `json:"food"`
	Opts   Options  `json:"options"`

	StartTime time.Time `json:"start"`
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
	return time.Now().After(g.StartTime.Add(g.Opts.Duration))
}

func (g *Game) JSON() []byte {
	data, _ := json.Marshal(g)
	return data
}

func (g *Game) ReloadFood() {
	g.Food.Point = Point{
		rand.Intn(g.Field.W),
		rand.Intn(g.Field.H),
	}
}

func (g *Game) Start() {
	g.ShuffleSnakes()
	g.ReloadFood()
	g.StartTime = time.Now()
}
