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
	return result
}
