package pattern

import (
	"bytes"
	"image"
	"image/color"
	"image/png"

	"github.com/czzle/czzle"
)

type Pattern struct {
	layers LayerSet
}

func (p Pattern) Split(cuts int, tt czzle.TileType) ([]*czzle.Tile, error) {
	bsize := Size / cuts
	dst := make([]*czzle.Tile, cuts*cuts)
	at := 0
	layers := p.layers

	layers = append(LayerSet{
		&Border{
			color: color.RGBA{
				A: 0,
			},
			bsize: bsize,
		},
	}, layers...)
	img := image.NewRGBA(image.Rectangle{
		Min: image.Point{0, 0},
		Max: image.Point{bsize, bsize},
	})
	var err error
	for blockX := 0; blockX < cuts; blockX++ {
		for blockY := 0; blockY < cuts; blockY++ {
			buf := bytes.NewBuffer(nil)
			offsetX, offsetY := blockX*bsize, blockY*bsize
			for x := 0; x < bsize; x++ {
				for y := 0; y < bsize; y++ {
					img.Set(x, y, layers.ColorAt(
						x+offsetX,
						y+offsetY,
					))
				}
			}
			err = png.Encode(buf, img)
			if err != nil {
				return nil, err
			}
			dst[at] = &czzle.Tile{
				Type: tt,
				Pos: &czzle.Pos{
					X: blockX,
					Y: blockY,
				},
				Data: buf.Bytes(),
			}
			at++
		}
	}
	return dst, nil
}
