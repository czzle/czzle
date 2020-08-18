package pattern

import (
	"image/color"
)

type Layer interface {
	Get(x, y int) bool
	Color() color.Color
}

type LayerSet []Layer

func (set LayerSet) ColorAt(x, y int) color.Color {
	for _, l := range set {
		if l.Get(x, y) {
			return l.Color()
		}
	}
	return nil
}
