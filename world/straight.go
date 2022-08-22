package world

import (
	"github.com/fafeitsch/go-infinite-rail-generator/util"
	"math/rand"
)

type straightBuilder struct {
	nameProvider nameProvider
}

func (s *straightBuilder) build(start int, values []float64) []Tile {
	random := rand.New(rand.NewSource(int64(values[0] * 1000)))
	die := random.Float64()
	tracks := 2
	if die > 0.5 && die <= 0.80 {
		tracks = 1
	} else if die > 0.80 && die <= 0.90 {
		tracks = 3
	} else if die > 0.9 {
		tracks = 4
	}
	result := make([]Tile, len(values))
	for index := range values {
		result[index] = NewTile("", tracks)
	}
	die = random.Float64()
	if die < 0.5 && len(result) > 2 {
		result = s.buildJunction(random, tracks, result)
		die = random.Float64()
		if die < 0.5 {
			util.Reverse(result)
			for index := range result {
				result[index].Mirror()
			}
		}
	}
	return result
}

func (s *straightBuilder) buildJunction(random *rand.Rand, tracks int, tiles []Tile) []Tile {
	verticalDirection := random.Intn(2)
	tileIndex := len(tiles) / 2
	index := 0
	if verticalDirection == 0 {
		index = 15
		verticalDirection = -1
	}
	minAlpha := 15
	for _, connector := range tiles[tileIndex].Tracks {
		if connector.SourceColumn == Alpha && connector.SourceTrack < minAlpha {
			minAlpha = connector.SourceTrack
		}
	}
	index = index + minAlpha*verticalDirection
	if tracks == 2 {
		return s.buildDirectJunction(index, verticalDirection, tiles)
	}
	tile := &tiles[tileIndex]
	sources := [3]Column{Alpha, Beta, Gamma}
	targets := [3]Column{Beta, Gamma, Omega}
	for track := 0; track < tracks; track++ {
		tile.Tracks = append(tile.Tracks, Connector{
			SourceColumn: sources[track%len(sources)],
			SourceTrack:  index,
			TargetColumn: targets[track%len(targets)],
			TargetTrack:  index + verticalDirection,
		})
		if track == 2 {
			tile = &tiles[tileIndex+1]
		}
		index = index + verticalDirection
		if track == tracks-2 {
			verticalDirection = 2 * verticalDirection
		}
	}
	tile.Tracks = append(tile.Tracks, Connector{
		SourceTrack:  index,
		SourceColumn: sources[tracks%len(sources)],
		TargetColumn: targets[tracks%len(targets)],
		TargetTrack:  index,
		Junction:     true,
		Label:        "to " + s.nameProvider(),
	})
	return tiles
}

func (s *straightBuilder) buildDirectJunction(index int, direction int, tiles []Tile) []Tile {
	tileIndex := len(tiles) / 2
	multiplier := 2
	if direction == 1 {
		multiplier = 3
	}
	tile := &tiles[tileIndex]
	aTrack := index + (multiplier * direction)
	bTrack := index + direction + (multiplier * direction)
	tile.Tracks = append(tile.Tracks, Connector{
		SourceTrack:  index,
		SourceColumn: Alpha,
		TargetColumn: Omega,
		TargetTrack:  aTrack,
	})
	tile.Tracks = append(tile.Tracks, Connector{
		SourceTrack:  index + direction,
		SourceColumn: Alpha,
		TargetTrack:  bTrack,
		TargetColumn: Omega,
	})
	tile = &tiles[tileIndex+1]
	tile.Tracks = append(tile.Tracks, Connector{
		SourceTrack:  aTrack,
		SourceColumn: Alpha,
		TargetColumn: Beta,
		TargetTrack:  aTrack,
		Junction:     true,
	})
	tile.Tracks = append(tile.Tracks, Connector{
		SourceTrack:  bTrack,
		SourceColumn: Alpha,
		TargetColumn: Beta,
		TargetTrack:  bTrack,
		Junction:     true,
		Label:        "to " + s.nameProvider(),
	})
	return tiles
}
