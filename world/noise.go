package world

import (
	"hash/fnv"
	"math/rand"
)

type noise struct {
	source [512]float64
}

func createNoise(seed string) *noise {
	hash := fnv.New64a()
	_, _ = hash.Write([]byte(seed))
	hashSum := hash.Sum64()
	source := rand.NewSource(int64(hashSum))
	random := rand.New(source)
	result := noise{source: [512]float64{}}
	for i, _ := range result.source {
		result.source[i] = random.Float64()
	}
	return &result
}

// sampling denotes the number of maximal different tiles in the world.
const sampling = 30_000

func (n *noise) interpolate(tile int) float64 {
	tile = tile % sampling
	if tile < 0 {
		tile = sampling + tile
	}
	x := float64(tile) * float64(len(n.source)) / float64(sampling)
	xMin := int(x)
	xMax := (xMin + 1) % len(n.source)
	deltaX := x - float64(xMin)
	smoothDeltaX := deltaX * deltaX * (3 - 2*deltaX)

	return n.source[xMin]*(1-smoothDeltaX) + n.source[xMax]*smoothDeltaX
}
