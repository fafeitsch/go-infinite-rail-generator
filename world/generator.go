package world

import (
	"github.com/fafeitsch/go-infinite-rail-generator/util"
	"math"
	"math/rand"
)

type Generator struct {
	seed  string
	noise *noise
}

func NewGenerator(seed string) *Generator {
	return &Generator{noise: createNoise(seed), seed: seed}
}

func (g *Generator) Seed() string {
	return g.seed
}

func (g *Generator) Generate(tileNumber int) Tile {
	value := g.noise.interpolate(tileNumber)
	values := make([]float64, 0)
	plateauStart := tileNumber
	current := g.noise.interpolate(tileNumber - 1)
	for ; math.Ceil(current*10) == math.Ceil(value*10); current = g.noise.interpolate(plateauStart - 1) {
		values = append(values, current)
		plateauStart = plateauStart - 1
	}
	util.Reverse(values)
	values = append(values, value)
	plateauEnd := tileNumber
	current = g.noise.interpolate(tileNumber + 1)
	for ; math.Ceil(current*10) == math.Ceil(value*10); current = g.noise.interpolate(plateauEnd + 1) {
		values = append(values, current)
		plateauEnd = plateauEnd + 1
	}
	factory := getSegmentFactory(values)
	tiles := factory(plateauStart, values)
	return *tiles[tileNumber-plateauStart]
}

type segmentFactory func(start int, values []float64) []*Tile

func getSegmentFactory(values []float64) segmentFactory {
	sum := 0.0
	for _, value := range values {
		sum = sum + (value * 10000)
	}
	random := rand.New(rand.NewSource(int64(sum)))
	die := random.Float64()
	if die > 0 {
		return straightTrack
	}
	return straightTrack
}
