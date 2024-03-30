package game

type Field struct {
	W int
	H int
}

func (f Field) InBounds(p Point) bool {
	return p.X >= 0 && p.X < f.W && p.Y >= 0 && p.Y < f.H
}
func (f Field) ToBounds(p Point) Point {
	return Point{(p.X + f.W) % f.W, (p.Y + f.H) % f.H}
}

func (f Field) HandleCollision(s *Snake) {
	if !f.InBounds(s.Head()) {
		s.SetHead(f.ToBounds(s.Head()))
	}
}
