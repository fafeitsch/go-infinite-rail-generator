package world

import (
	"github.com/fafeitsch/go-infinite-rail-generator/util"
	"math"
	"math/rand"
)

type Generator struct {
	seed      string
	noise     *noise
	TownNames []string
}

func NewGenerator(seed string) *Generator {
	return &Generator{
		noise:     createNoise(seed),
		seed:      seed,
		TownNames: make([]string, 0),
	}
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
	factory := g.getSegmentFactory(values)
	tiles := factory.build(plateauStart, values)
	return plateauStart, tiles
}

type segmentFactory interface {
	build(int, []float64) []Tile
}

type nameProvider func() string

func (g *Generator) getSegmentFactory(values []float64) segmentFactory {
	sum := 0.0
	for _, value := range values {
		sum = sum + (value * 10000)
	}
	random := rand.New(rand.NewSource(int64(sum)))
	// die := random.Float64()
	nameProvider := func() string {
		if len(g.TownNames) == 0 {
			return ""
		}
		index := random.Intn(len(g.TownNames))
		return g.TownNames[index]
	}
	if len(values)%2 == 0 && len(values) <= 10 {
		return &stationBuilder{nameProvider: nameProvider}
	}
	return &straightBuilder{nameProvider: nameProvider}
}

func fixNecessarySwitches(left *Tile, right Tile) {
	minRight := 15.0
	maxRight := 0.0
	for _, connector := range right.Tracks {
		if connector.SourceColumn == Alpha {
			minRight = math.Min(minRight, float64(connector.SourceTrack))
			maxRight = math.Max(maxRight, float64(connector.SourceTrack))
		}
	}
	minLeft := 15.0
	maxLeft := 0.0
	for index, connector := range left.Tracks {
		if connector.TargetColumn != Omega {
			continue
		}
		minLeft = math.Min(minLeft, float64(connector.SourceTrack))
		maxLeft = math.Max(maxLeft, float64(connector.SourceTrack))
		if connector.TargetTrack < int(minRight) {
			left.Tracks[index].TargetTrack = int(minRight)
		}
		if connector.TargetTrack > int(maxRight) {
			left.Tracks[index].TargetTrack = int(maxRight)
		}
	}
	for index := int(minRight); index < int(minLeft); index++ {
		left.Tracks = append(left.Tracks, Connector{
			SourceTrack:  int(minLeft),
			SourceColumn: Gamma,
			TargetTrack:  index,
			TargetColumn: Omega,
		})
	}
	for index := int(maxRight); index > int(maxLeft); index-- {
		left.Tracks = append(left.Tracks, Connector{
			SourceTrack:  int(maxLeft),
			SourceColumn: Gamma,
			TargetTrack:  index,
			TargetColumn: Omega,
		})
	}
}
