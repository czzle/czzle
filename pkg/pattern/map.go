package pattern

import "image/color"

type Map struct {
	points []bool
	n      int
	size   int
	Color  color.Color
}

func (p *Map) Set(x, y int, v bool) {
	at := y*p.size + x
	if at > p.n || at < 0 {
		return
	}
	p.points[at] = v
}

func (p Map) Get(x, y int) bool {
	at := y*p.size + x
	if at > p.n || at < 0 {
		return false
	}
	return p.points[at]
}

func NewMap(size int, c color.Color) *Map {
	n := size * size
	return &Map{
		points: make([]bool, n),
		n:      n,
		size:   size,
		Color:  c,
	}
}
