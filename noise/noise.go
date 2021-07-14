package noise

import (
	"hash/fnv"
	"math/rand"
)

type Noise struct {
	source [512]float64
	Seed   string
}

func New(seed string) *Noise {
	hash := fnv.New64a()
	_, _ = hash.Write([]byte(seed))
	hashSum := hash.Sum64()
	source := rand.NewSource(int64(hashSum))
	random := rand.New(source)
	result := Noise{source: [512]float64{}, Seed: seed}
	for i, _ := range result.source {
		result.source[i] = random.Float64()
	}
	return &result
}

func (n *Noise) Interpolate(hectometer int) float64 {
	hectometer = hectometer % 10_000
	if hectometer < 0 {
		hectometer = 10_000 + hectometer
	}
	x := float64(hectometer) * float64(len(n.source)) / 10_000
	xMin := int(x)
	xMax := (xMin + 1) % len(n.source)
	deltaX := x - float64(xMin)
	smoothDeltaX := deltaX * deltaX * (3 - 2*deltaX)

	return n.source[xMin]*(1-smoothDeltaX) + n.source[xMax]*smoothDeltaX
}
