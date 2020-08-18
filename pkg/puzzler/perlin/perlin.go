package perlin

import (
	"math"
	"math/rand"
)

// General constants
const (
	B  = 0x100
	N  = 0x1000
	BM = 0xff
)

// Perlin is the noise generator
type Perlin struct {
	alpha float64
	beta  float64
	n     int

	p  [B + B + 2]int
	g3 [B + B + 2][3]float64
	g2 [B + B + 2][2]float64
	g1 [B + B + 2]float64
}

func New(alpha, beta float64, n int, seed int64) *Perlin {
	return NewRandSource(alpha, beta, n, rand.NewSource(seed))
}

func NewRandSource(alpha, beta float64, n int, source rand.Source) *Perlin {
	var p Perlin
	var i, j, k int

	p.alpha = alpha
	p.beta = beta
	p.n = n

	r := rand.New(source)

	for i = 0; i < B; i++ {
		p.p[i] = i
		p.g1[i] = float64((r.Int()%(B+B))-B) / B

		for j = 0; j < 2; j++ {
			p.g2[i][j] = float64((r.Int()%(B+B))-B) / B
		}

		normalize2(&p.g2[i])

		for j = 0; j < 3; j++ {
			p.g3[i][j] = float64((r.Int()%(B+B))-B) / B
		}
		normalize3(&p.g3[i])
	}

	for ; i > 0; i-- {
		k = p.p[i]
		j = r.Int() % B
		p.p[i] = p.p[j]
		p.p[j] = k
	}

	for i = 0; i < B+2; i++ {
		p.p[B+i] = p.p[i]
		p.g1[B+i] = p.g1[i]
		for j = 0; j < 2; j++ {
			p.g2[B+i][j] = p.g2[i][j]
		}
		for j = 0; j < 3; j++ {
			p.g3[B+i][j] = p.g3[i][j]
		}
	}

	return &p
}

func normalize2(v *[2]float64) {
	s := math.Sqrt(v[0]*v[0] + v[1]*v[1])
	v[0] = v[0] / s
	v[1] = v[1] / s
}

func normalize3(v *[3]float64) {
	s := math.Sqrt(v[0]*v[0] + v[1]*v[1] + v[2]*v[2])
	v[0] = v[0] / s
	v[1] = v[1] / s
	v[2] = v[2] / s
}

func at2(rx, ry float64, q [2]float64) float64 {
	return rx*q[0] + ry*q[1]
}

func at3(rx, ry, rz float64, q [3]float64) float64 {
	return rx*q[0] + ry*q[1] + rz*q[2]
}

func sCurve(t float64) float64 {
	return t * t * (3. - 2.*t)
}

func lerp(t, a, b float64) float64 {
	return a + t*(b-a)
}

func (p *Perlin) noise1(arg float64) float64 {
	var vec [1]float64
	vec[0] = arg

	t := vec[0] + N
	bx0 := int(t) & BM
	bx1 := (bx0 + 1) & BM
	rx0 := t - float64(int(t))
	rx1 := rx0 - 1.

	sx := sCurve(rx0)
	u := rx0 * p.g1[p.p[bx0]]
	v := rx1 * p.g1[p.p[bx1]]

	return lerp(sx, u, v)
}

func (p *Perlin) noise2(vec [2]float64) float64 {
	t := vec[0] + N
	bx0 := int(t) & BM
	bx1 := (bx0 + 1) & BM
	rx0 := t - float64(int(t))
	rx1 := rx0 - 1.

	t = vec[1] + N
	by0 := int(t) & BM
	by1 := (by0 + 1) & BM
	ry0 := t - float64(int(t))
	ry1 := ry0 - 1.

	i := p.p[bx0]
	j := p.p[bx1]

	b00 := p.p[i+by0]
	b10 := p.p[j+by0]
	b01 := p.p[i+by1]
	b11 := p.p[j+by1]

	sx := sCurve(rx0)
	sy := sCurve(ry0)

	q := p.g2[b00]
	u := at2(rx0, ry0, q)
	q = p.g2[b10]
	v := at2(rx1, ry0, q)
	a := lerp(sx, u, v)

	q = p.g2[b01]
	u = at2(rx0, ry1, q)
	q = p.g2[b11]
	v = at2(rx1, ry1, q)
	b := lerp(sx, u, v)

	return lerp(sy, a, b)
}

func (p *Perlin) noise3(vec [3]float64) float64 {
	t := vec[0] + N
	bx0 := int(t) & BM
	bx1 := (bx0 + 1) & BM
	rx0 := t - float64(int(t))
	rx1 := rx0 - 1.

	t = vec[1] + N
	by0 := int(t) & BM
	by1 := (by0 + 1) & BM
	ry0 := t - float64(int(t))
	ry1 := ry0 - 1.

	t = vec[2] + N
	bz0 := int(t) & BM
	bz1 := (bz0 + 1) & BM
	rz0 := t - float64(int(t))
	rz1 := rz0 - 1.

	i := p.p[bx0]
	j := p.p[bx1]

	b00 := p.p[i+by0]
	b10 := p.p[j+by0]
	b01 := p.p[i+by1]
	b11 := p.p[j+by1]

	t = sCurve(rx0)
	sy := sCurve(ry0)
	sz := sCurve(rz0)

	q := p.g3[b00+bz0]
	u := at3(rx0, ry0, rz0, q)
	q = p.g3[b10+bz0]
	v := at3(rx1, ry0, rz0, q)
	a := lerp(t, u, v)

	q = p.g3[b01+bz0]
	u = at3(rx0, ry1, rz0, q)
	q = p.g3[b11+bz0]
	v = at3(rx1, ry1, rz0, q)
	b := lerp(t, u, v)

	c := lerp(sy, a, b)

	q = p.g3[b00+bz1]
	u = at3(rx0, ry0, rz1, q)
	q = p.g3[b10+bz1]
	v = at3(rx1, ry0, rz1, q)
	a = lerp(t, u, v)

	q = p.g3[b01+bz1]
	u = at3(rx0, ry1, rz1, q)
	q = p.g3[b11+bz1]
	v = at3(rx1, ry1, rz1, q)
	b = lerp(t, u, v)

	d := lerp(sy, a, b)

	return lerp(sz, c, d)
}

func (p *Perlin) Noise(x, y float64) float64 {
	var scale float64 = 1
	var sum, val float64
	var px [2]float64

	px[0] = x
	px[1] = y

	for i := 0; i < p.n; i++ {
		val = p.noise2(px)
		sum += val / scale
		scale *= p.alpha
		px[0] *= p.beta
		px[1] *= p.beta
	}
	return sum
}
