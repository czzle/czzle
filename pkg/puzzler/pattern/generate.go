package pattern

import (
	"errors"
	"hash/fnv"
	"image/color"

	"github.com/czzle/czzle"
	"github.com/czzle/czzle/pkg/puzzler/perlin"
)

var (
	ErrUnknownPattern      = errors.New("unknown pattern")
	ErrInvalidCuts         = errors.New("invalid cuts")
	ErrPalleteMissingColor = errors.New("pallette missing colors")
)

func Generate(g []byte, opts Options) ([]*czzle.Tile, error) {
	if opts.Cuts < 2 || opts.Cuts > 3 {
		return nil, ErrInvalidCuts
	}
	for _, c := range opts.Pallette {
		if c == nil {
			return nil, ErrPalleteMissingColor
		}
	}
	var pattern *Pattern
	switch opts.Style {
	case BlobStyle:
		pattern = generateBlob(g, opts.Pallette)
	default:
		return nil, ErrUnknownPattern
	}
	return pattern.Split(opts.Cuts, opts.TileType)
}

type PatternStyle int

const (
	BlobStyle PatternStyle = iota
)

type Options struct {
	Cuts     int
	Style    PatternStyle
	TileType czzle.TileType
	Pallette [3]color.Color
}

func generateBlob(g []byte, pallette [3]color.Color) *Pattern {
	// make new rand seed from generator
	h := fnv.New64a()
	h.Write([]byte(g))
	// consts
	alpha := 1.
	beta := 1.
	n := 2
	scale := 150.
	decrease := 75.
	split := -0.08
	// create perlin
	p := perlin.New(alpha, beta, n, int64(h.Sum64()))
	// set layer set
	layers := make(LayerSet, 3)
	layers[2] = &Background{
		color: pallette[2],
	}
	// paint first layers
	for i := 0; i < 2; i++ {
		l := &ColorMap{
			color: pallette[i],
		}
		for x := 0; x < Size; x++ {
			for y := 0; y < Size; y++ {
				f := p.Noise(
					float64(x)/scale,
					float64(y)/scale,
				)
				l.Set(x, y, f < split)
			}
		}
		layers[i] = l
		scale -= decrease
	}
	return &Pattern{
		layers: layers,
	}
}
