package noise

import (
	"math/rand"
)

type Noise struct {
	source [512]float64
}

func New(seed int64) *Noise {
	source := rand.NewSource(seed)
	random := rand.New(source)
	result := Noise{source: [512]float64{}}
	for i, _ := range result.source {
		result.source[i] = random.Float64()
	}
	return &result
}

func (n *Noise) NumberOfTracks(hectometer int) int {
	hectometer = hectometer % 10_000
	if hectometer < 0 {
		hectometer = 10_000 + hectometer
	}
	x := float64(hectometer) * float64(len(n.source)) / 10_000
	xMin := int(x)
	xMax := (xMin + 1) % len(n.source)
	deltaX := x - float64(xMin)
	smoothDeltaX := deltaX * deltaX * (3 - 2*deltaX)

	result := n.source[xMin]*(1-smoothDeltaX) + n.source[xMax]*smoothDeltaX
	if result < 0.2 {
		return 1
	}
	if result < 0.6 {
		return 2
	}
	return int(result*10 - 3)
}
