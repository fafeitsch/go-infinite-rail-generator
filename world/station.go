package world

import (
	"math/rand"
)

type stationBuilder struct {
	nameProvider nameProvider
}

func (s *stationBuilder) build(start int, values []float64) []Tile {
	random := rand.New(rand.NewSource(int64(start)))
	die := random.Float64()
	tracks := 2
	if die > 0.25 && die <= 0.45 {
		tracks = 1
	} else if die > 0.45 && die <= 0.6 {
		tracks = 3
	} else if die > 0.85 && die <= 0.9 {
		tracks = 4
	} else if die > 0.9 && die <= 0.95 {
		tracks = 5
	} else if die > 0.95 {
		tracks = 6
	}
	result := make([]Tile, len(values))
	first := 8 - tracks/2 + random.Intn(2)
	for index := range values {
		tile := NewTile("", tracks)
		for track := first; track < first+tracks; track++ {
			tile.Platforms = append(tile.Platforms, Platform{
				Column: Alpha,
				Track:  track,
			}, Platform{
				Column: Beta,
				Track:  track,
			})
			if index != len(values)-1 {
				tile.Platforms = append(tile.Platforms, Platform{
					Column: Gamma,
					Track:  track,
				})
			}
		}
		result[index] = tile
	}
	result[len(result)/2].StationName = s.nameProvider()
	return result
}
