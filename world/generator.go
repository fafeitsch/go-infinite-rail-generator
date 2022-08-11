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
	plateauStart, tiles := g.buildSegment(tileNumber)
	tileIndex := tileNumber - plateauStart
	tile := tiles[tileIndex]
	if tileIndex == len(tiles)-1 {
		_, nextSegment := g.buildSegment(plateauStart + len(tiles) + 1)
		fixNecessarySwitches(&tile, nextSegment[0])
	}
	return tile
}

func (g *Generator) buildSegment(tileNumber int) (int, []Tile) {
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
	return plateauStart, tiles
}

type segmentFactory func(start int, values []float64) []Tile

func getSegmentFactory(values []float64) segmentFactory {
	sum := 0.0
	for _, value := range values {
		sum = sum + (value * 10000)
	}
	random := rand.New(rand.NewSource(int64(sum)))
	die := random.Float64()
	if len(values)%2 == 0 && die < 0.5 {
		return station
	}
	return straightTrack
}

func fixNecessarySwitches(left *Tile, right Tile) {
	rightConnectors := right.Tracks.AlphaTracks()
	for i, track := range left.Tracks.Gamma {
		underTest := track.FindConnector(Omega, i)
		var j int
		if underTest != nil && !rightConnectors[i] {
			if i <= len(left.Tracks.Gamma)/2 {
				j = i + 1
				for !rightConnectors[j] && j <= len(left.Tracks.Gamma)-1 {
					j = j + 1
				}
			} else {
				j = i - 1
				for !rightConnectors[j] && j >= 0 {
					j = j - 1
				}
			}
			underTest.Slot = j
		} else if underTest == nil && rightConnectors[i] {
			if i <= len(left.Tracks.Gamma)/2 {
				j = i + 1
				for len(left.Tracks.Gamma[j]) == 0 {
					j = j + 1
				}
			} else {
				j = i - 1
				for len(left.Tracks.Gamma[j]) == 0 {
					j = j - 1
				}
			}
			left.Tracks.Gamma[j] = append(left.Tracks.Gamma[j], &Connector{
				Target: Omega,
				Slot:   i,
			})
		}
	}
}
