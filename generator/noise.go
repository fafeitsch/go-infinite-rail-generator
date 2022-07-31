package generator

import (
	"hash/fnv"
	"math/rand"
)

type noise struct {
	source [512]float64
}

type Generator struct {
	*noise
	Seed      string
	TownNames []string
}

func New(seed string) *Generator {
	hash := fnv.New64a()
	_, _ = hash.Write([]byte(seed))
	hashSum := hash.Sum64()
	return &Generator{noise: createNoise(hashSum), Seed: seed}
}

func createNoise(hashSum uint64) *noise {
	source := rand.NewSource(int64(hashSum))
	random := rand.New(source)
	result := noise{source: [512]float64{}}
	for i, _ := range result.source {
		result.source[i] = random.Float64()
	}
	return &result
}

func (g *Generator) derive() *noise {
	hash := fnv.New64a()
	_, _ = hash.Write([]byte(g.Seed))
	hashSum := hash.Sum64()
	result := createNoise(hashSum + 43)
	return result
}

func (n *noise) interpolate(hectometer int, sampling int) float64 {
	hectometer = hectometer % sampling
	if hectometer < 0 {
		hectometer = sampling + hectometer
	}
	x := float64(hectometer) * float64(len(n.source)) / float64(sampling)
	xMin := int(x)
	xMax := (xMin + 1) % len(n.source)
	deltaX := x - float64(xMin)
	smoothDeltaX := deltaX * deltaX * (3 - 2*deltaX)

	return n.source[xMin]*(1-smoothDeltaX) + n.source[xMax]*smoothDeltaX
}
