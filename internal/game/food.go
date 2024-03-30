package game

type Food struct {
	Point Point `json:"point"`
	game  *Game
}

func (f Food) HandleCollision(s *Snake) {
	if s.Head() == f.Point {
		s.Grow()
		s.game.ReloadFood()
	}
}
