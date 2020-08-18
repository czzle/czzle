package pattern

import (
	"image/color"
)

type ColorMap struct {
	points [90000]bool
	color  color.Color
}

func (m *ColorMap) Set(x, y int, v bool) {
	m.points[y*300+x] = v
}

func (m *ColorMap) Get(x, y int) bool {
	return m.points[y*300+x]
}

func (m *ColorMap) Color() color.Color {
	return m.color
}
