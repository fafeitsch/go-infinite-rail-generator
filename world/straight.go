package world

import "math/rand"

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
	}
	return result
}

func (s *straightBuilder) buildJunction(random *rand.Rand, tracks int, tiles []Tile) []Tile {
	verticalDirection := random.Intn(2)
	tileIndex := len(tiles) / 2
	connectors := &tiles[tileIndex].Tracks.Alpha
	index := 0
	if verticalDirection == 0 {
		index = len(connectors) - 1
		verticalDirection = -1
	}
	for ; len(connectors[index]) == 0; index = index + verticalDirection {
	}
	if tracks == 2 {
		return s.buildDirectJunction(index, verticalDirection, tiles)
	}
	target := Beta
	for track := 0; track < tracks; track++ {
		connectors[index] = append(connectors[index], &Connector{
			Target: target,
			Track:  index + verticalDirection,
		})
		target = target.Next()
		if track == 0 {
			connectors = &tiles[tileIndex].Tracks.Beta
		} else if track == 1 {
			connectors = &tiles[tileIndex].Tracks.Gamma
		} else if track == 2 {
			connectors = &tiles[tileIndex+1].Tracks.Alpha
		} else if track == 3 {
			connectors = &tiles[tileIndex+1].Tracks.Beta
		}
		index = index + verticalDirection
		if track == tracks-2 {
			verticalDirection = 2 * verticalDirection
		}
	}
	connectors[index] = append(connectors[index], &Connector{
		Target:   target,
		Track:    index,
		Junction: true,
		Label:    "to " + s.nameProvider(),
	})
	return tiles
}

func (s *straightBuilder) buildDirectJunction(index int, direction int, tiles []Tile) []Tile {
	tileIndex := len(tiles) / 2
	multiplier := 2
	if direction == 1 {
		multiplier = 3
	}
	connectors := &tiles[tileIndex].Tracks.Alpha
	aTrack := index + (multiplier * direction)
	connectors[index] = append(connectors[index], &Connector{
		Target: Omega,
		Track:  aTrack,
	})
	bTrack := index + direction + (multiplier * direction)
	connectors[index+direction] = append(connectors[index+direction], &Connector{
		Target: Omega,
		Track:  bTrack,
	})
	connectors = &tiles[tileIndex+1].Tracks.Alpha
	connectors[aTrack] = append(connectors[aTrack], &Connector{
		Target:   Beta,
		Track:    aTrack,
		Junction: true,
	})
	connectors[bTrack] = append(connectors[bTrack], &Connector{
		Target:   Beta,
		Track:    bTrack,
		Junction: true,
		Label:    "to " + s.nameProvider(),
	})
	return tiles
}
