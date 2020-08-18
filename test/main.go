package main

import (
	"crypto/rand"
	"github.com/czzle/czzle/pkg/pattern"
	"encoding/base64"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math/big"
	"os"
	"time"
)

const (
	size     = 300
	split    = -0.08
	scale    = 150.
	increase = -75.
)

func main() {
	start := time.Now()
	maps := []*pattern.Map{
		pattern.NewMap(size, color.RGBA{
			R: 96,
			G: 162,
			B: 193,
			A: 255,
		}),
		pattern.NewMap(size, color.RGBA{
			R: 42,
			G: 74,
			B: 111,
			A: 255,
		}),
	}
	for i, m := range maps {
		s := scale + (float64(i) * increase)
		err := pattern.Generate(s, split, m)
		if err != nil {
			panic(err)
		}
	}
	cuts := 3
	block := size / cuts
	bg := color.RGBA{
		R: 21,
		G: 37,
		B: 55,
		A: 255,
	}
	am := make([]image.Image, cuts*cuts)
	bm := make([]image.Image, cuts*cuts)

	for i := 0; i < cuts; i++ {
		for j := 0; j < cuts; j++ {

			img := image.NewRGBA(image.Rectangle{
				Min: image.Point{0, 0},
				Max: image.Point{block, block},
			})
			offsetX := block * i
			offsetY := block * j
			for x := 0; x < block; x++ {
			rowa:
				for y := 0; y < block; y++ {
					for _, m := range maps {
						if m.Get(x+offsetX, y+offsetY) {
							img.Set(x, y, m.Color)
							continue rowa
						}
					}
					img.Set(x, y, bg)
				}
			}
			am[j*cuts+i] = img
		}
	}
	// regenerate patterns
	for i, m := range maps {
		s := scale + (float64(i) * increase)
		err := pattern.Generate(s, split, m)
		if err != nil {
			panic(err)
		}
	}

	for i := 0; i < cuts; i++ {
		for j := 0; j < cuts; j++ {

			img := image.NewRGBA(image.Rectangle{
				Min: image.Point{0, 0},
				Max: image.Point{block, block},
			})
			offsetX := block * i
			offsetY := block * j
			for x := 0; x < block; x++ {
			rowb:
				for y := 0; y < block; y++ {
					for _, m := range maps {
						if m.Get(x+offsetX, y+offsetY) {
							img.Set(x, y, m.Color)
							continue rowb
						}
					}
					img.Set(x, y, bg)
				}
			}
			bm[j*cuts+i] = img
		}
	}
	bi, err := rand.Int(rand.Reader, big.NewInt(int64(cuts*cuts)))
	if err != nil {
		panic(err)
	}
	si := int(bi.Int64())
	am[si], bm[si] = bm[si], am[si]
	for i, img := range am {
		y := i / cuts
		x := i % cuts
		f, _ := os.Create(fmt.Sprintf("web/%dx%d.png", x, y))
		png.Encode(f, img)
		f.Close()
	}
	for i, img := range bm {
		y := i / cuts
		x := i % cuts
		f, _ := os.Create(fmt.Sprintf("web/r%dx%d.png", x, y))
		png.Encode(f, img)
		f.Close()
	}
	fmt.Println("took:", time.Since(start))
	tmp := make([]byte, 8+32+32)
	rand.Read(tmp)
	fmt.Println(base64.RawURLEncoding.EncodeToString(tmp))
}

func rotate(times, offset, x, y int) (int, int) {
	times = times % 3
	if times < 0 {
		times++
		tmp := (x-offset)*(-1) + offset*2
		x, y = y, tmp
	} else if times > 0 {
		times--
		tmp := (y-offset)*(-1) + offset*2
		x, y = tmp, x
	} else {
		return x, y
	}
	return rotate(times, offset, x, y)
}
