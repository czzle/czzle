package puzzler

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"hash/fnv"
	"image/color"
	"time"

	"github.com/czzle/czzle"
	"github.com/czzle/czzle/pkg/puzzler/pattern"
	"github.com/czzle/czzle/pkg/token"
)

type Puzzler struct {
	secret string
}

func New(secret string) *Puzzler {
	return &Puzzler{secret}
}

func (p *Puzzler) Make(client *czzle.ClientInfo) (*czzle.Puzzle, error) {
	iat := time.Now().UnixNano() / int64(time.Millisecond)
	exp := time.Now().Add(time.Minute*6).UnixNano() / int64(time.Millisecond)
	g := make([]byte, 32)
	_, err := rand.Read(g)
	if err != nil {
		return nil, err
	}
	tkn := token.New(token.Payload{
		ClientID:  client.GetID(),
		ClientIP:  client.GetIP(),
		IssuedAt:  iat,
		ExpiresAt: exp,
		Level:     czzle.Hard,
		Generator: base64.RawStdEncoding.EncodeToString(g),
		Solved:    false,
	})
	err = tkn.Sign(p.secret)
	if err != nil {
		return nil, err
	}
	h := sha256.New()
	h.Write(g)
	h.Write([]byte(p.secret))
	g = h.Sum(nil)
	front, err := pattern.Generate(g, pattern.Options{
		Cuts:     3,
		Style:    pattern.BlobStyle,
		TileType: czzle.FrontTile,
		Pallette: [3]color.Color{
			color.RGBA{
				R: 96,
				G: 162,
				B: 193,
				A: 255,
			},
			color.RGBA{
				R: 42,
				G: 74,
				B: 111,
				A: 255,
			},
			color.RGBA{
				R: 21,
				G: 37,
				B: 55,
				A: 255,
			},
		},
	})
	if err != nil {
		panic(err)
	}
	h = sha256.New()
	h.Write(g)
	g = h.Sum(nil)
	back, err := pattern.Generate(g, pattern.Options{
		Cuts:     3,
		Style:    pattern.BlobStyle,
		TileType: czzle.BackTile,
		Pallette: [3]color.Color{
			color.RGBA{
				R: 96,
				G: 162,
				B: 193,
				A: 255,
			},
			color.RGBA{
				R: 42,
				G: 74,
				B: 111,
				A: 255,
			},
			color.RGBA{
				R: 21,
				G: 37,
				B: 55,
				A: 255,
			},
		},
	})
	if err != nil {
		panic(err)
	}
	fh := fnv.New64()
	fh.Write(g)
	i := int(fh.Sum64() % 9)
	front[i].Data, back[i].Data = back[i].Data, front[i].Data
	tiles := front
	tiles = append(tiles, back...)
	return &czzle.Puzzle{
		Token:     tkn.String(),
		ExpiresAt: exp,
		IssuedAt:  iat,
		Level:     czzle.Hard,
		Client:    client,
		TileMap: &czzle.TileMap{
			Size:  3,
			Tiles: tiles,
		},
	}, nil
}
