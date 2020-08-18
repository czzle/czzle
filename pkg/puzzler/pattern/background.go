package pattern

import (
	"image/color"
)

type Background struct {
	color color.Color
}

func (Background) Get(x, y int) bool {
	return true
}

func (bg Background) Color() color.Color {
	return bg.color
}
