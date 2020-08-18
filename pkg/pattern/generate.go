package pattern

import (
	crand "crypto/rand"
	"math"
	"math/big"
	"math/rand"

	"github.com/czzle/czzle/pkg/perlin"
)

const (
	alpha = 1.
	beta  = 1.
	n     = 2
)

func randPerlin() (*perlin.Perlin, error) {
	ri, err := crand.Int(crand.Reader, big.NewInt(math.MaxInt64))
	if err != nil {
		return nil, err
	}
	r := rand.NewSource(ri.Int64())
	p := perlin.NewRandSource(alpha, beta, n, r)
	return p, nil
}

func Generate(scale, split float64, m *Map) error {
	ri, err := crand.Int(crand.Reader, big.NewInt(math.MaxInt64))
	if err != nil {
		return err
	}
	p := perlin.New(alpha, beta, n, ri.Int64())
	for x := 0; x < m.size; x++ {
		for y := 0; y < m.size; y++ {
			f := p.Noise(
				float64(x)/scale,
				float64(y)/scale,
			)
			m.Set(x, y, f < split)
		}
	}
	return nil
}
