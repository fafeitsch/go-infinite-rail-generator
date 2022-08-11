package world

import (
	_ "embed"
	"math/rand"
	"strings"
)

//go:embed default_town_names.txt
var townNames string

func station(start int, values []float64) []Tile {
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
		platform := &tile.Platforms
		for track := first; track < first+tracks; track++ {
			platform.Alpha[track] = true
			platform.Beta[track] = true
			platform.Gamma[track] = true
		}
		result[index] = tile
	}
	platform := &result[len(result)-1].Platforms
	for track := first; track < first+tracks; track++ {
		platform.Gamma[track] = false
	}
	result[len(result)/2].StationName = getStationName(random)
	return result
}

func getStationName(random *rand.Rand) string {
	list := strings.Split(townNames, "\n")
	index := random.Intn(len(list))
	return list[index]
}
