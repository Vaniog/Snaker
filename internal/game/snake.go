package game

type Snake struct {
	Body []Point   `json:"body"`
	Drc  Direction `json:"direction"`
	// Len can be unequal to Len(Body), so snake will grow
	Len   int  `json:"len"`
	Alive bool `json:"alive"`

	game *Game

	// for rollback
	lastTail Point
}

func NewSnake(game *Game, head Point, drc Direction, len int) *Snake {
	s := Snake{Len: len, game: game}
	s.Alive = true
	s.Body = append(s.Body, head)
	s.lastTail = s.Head()
	s.Drc = drc
	return &s
}

func (s *Snake) Rotate(drc Direction) {
	if s.Drc.Opposite() != drc {
		s.Drc = drc
	}
}

func (s *Snake) Move() {
	head := s.Body[len(s.Body)-1]
	if len(s.Body) == s.Len {
		s.lastTail = s.Body[0]
		s.Body = s.Body[1:]
	}
	s.Body = append(s.Body, head.Move(s.Drc))
}

func (s *Snake) Rollback() {
	s.Body = s.Body[:len(s.Body)-1]
	s.Body = append([]Point{s.lastTail}, s.Body...)
}

func (s *Snake) Head() Point {
	return s.Body[len(s.Body)-1]
}

func (s *Snake) SetHead(p Point) {
	s.Body[len(s.Body)-1] = p
}

func (s *Snake) Grow() {
	s.Len++
}

func (s *Snake) HandleCollision(s2 *Snake) {
	for _, p := range s.Body {
		if s2.Head() == p {
			s2.Alive = false
			s2.Rollback()
		}
	}
}
