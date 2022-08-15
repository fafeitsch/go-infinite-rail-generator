package world

import "math/rand"

func straightTrack(start int, values []float64) []Tile {
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
		result = buildJunction(random, tracks, result)
	}
	return result
}

func buildJunction(random *rand.Rand, tracks int, tiles []Tile) []Tile {
	size := random.Intn(2) + 1
	direction := random.Intn(2)
	tileIndex := len(tiles) / 2
	connectors := &tiles[tileIndex].Tracks.Alpha
	index := 0
	if direction == 0 {
		index = len(connectors) - 1
		direction = -1
	}
	for ; len(connectors[index]) == 0; index = index + direction {
	}
	if size == tracks && tracks == 2 {
		connectors[index] = append(connectors[index], &Connector{
			Target: Gamma,
			Slot:   index + (direction * tracks),
		})
		connectors[index+direction] = append(connectors[index+direction], &Connector{
			Target: Gamma,
			Slot:   index + direction + (direction * tracks),
		})
		return tiles
	}
	target := Beta
	for track := 0; track < tracks; track++ {
		connectors[index] = append(connectors[index], &Connector{
			Target: target,
			Slot:   index + direction,
		})
		if size == 2 {
			connectors[index+direction] = append(connectors[index+direction], &Connector{
				Target: target,
				Slot:   index + direction + direction,
			})
		}
		target = target.Next()
		if track == 0 {
			connectors = &tiles[tileIndex].Tracks.Beta
		} else if track == 1 {
			connectors = &tiles[tileIndex].Tracks.Gamma
		} else if track == 2 {
			connectors = &tiles[tileIndex+1].Tracks.Alpha
		}
		index = index + direction
	}
	return tiles
}
