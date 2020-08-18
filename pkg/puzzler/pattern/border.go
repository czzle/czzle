package pattern

import (
	"image/color"
)

type Border struct {
	color color.Color
	bsize int
}

func (b Border) Get(x, y int) bool {
	bX := x % b.bsize
	bY := y % b.bsize
	if bX < 5 || bX > b.bsize-5 || bY < 5 || bY > b.bsize-5 {
		return true
	}
	return false
}

func (b Border) Color() color.Color {
	return b.color
}
