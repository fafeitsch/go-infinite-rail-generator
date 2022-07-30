package noise

import (
	_ "embed"
	"github.com/fafeitsch/go-infinite-rail-generator/domain"
	"math"
	"strings"
)

//go:embed default_town_names.txt
var defaultTownNames string

func (g *Generator) buildStation(tiles []*rndTile) {
	if len(tiles) == 0 {
		return
	}
	for index, tile := range tiles {
		minConnected := len(tile.Tracks.Alpha) + 1
		maxConnected := 0
		for i, connected := range tiles[0].Tracks.AlphaTracks() {
			if connected {
				minConnected = int(math.Min(float64(minConnected), float64(i)))
				maxConnected = int(math.Max(float64(maxConnected), float64(i)))
				tile.Tracks.Alpha[i] = []*domain.Connector{
					{
						Target: domain.Beta,
						Slot:   i,
					},
				}
				tile.Tracks.Beta[i] = []*domain.Connector{
					{
						Target: domain.Gamma,
						Slot:   i,
					},
				}
				tile.Tracks.Gamma[i] = []*domain.Connector{
					{
						Target: domain.Omega,
						Slot:   i,
					},
				}
			}
		}
		for i, connected := range tiles[0].Tracks.AlphaTracks() {
			if !connected {
				continue
			}
			if !tile.Platforms.Alpha[i] && !tile.Platforms.Alpha[i+1] {
				if tiles[0].createRandom(98).Float64() < 0.5 && i%2 == 1 {
					tile.Platforms.Alpha[i] = true
					tile.Platforms.Beta[i] = true
					tile.Platforms.Gamma[i] = index != len(tiles)
				} else {
					tile.Platforms.Alpha[i+1] = true
					tile.Platforms.Beta[i+1] = true
					tile.Platforms.Gamma[i+1] = index != len(tiles)
				}
			}
		}
		if index == len(tiles)/2 {
			tile.StationName = g.getStationName(tile)
		}
	}
}

func (g *Generator) getStationName(tile *rndTile) string {
	if len(g.TownNames) == 0 {
		g.TownNames = strings.Split(defaultTownNames, "\n")
	}
	index := tile.createRandom(99).Intn(len(g.TownNames))
	return g.TownNames[index]
}

func stationWidth(start int, generator func(int) *rndTile) (int, []*rndTile) {
	tiles := make([]*rndTile, 0, 0)
	i := start - 1
	tile := generator(i)
	for tile.station {
		i = i - 1
		tiles = append(tiles, tile)
		tile = generator(i)
	}
	for a, b := 0, len(tiles)-1; a < b; a, b = a+1, b-1 {
		tiles[a], tiles[b] = tiles[b], tiles[a]
	}
	tiles = append(tiles, generator(start))
	j := start + 1
	tile = generator(j)
	for tile.station {
		j = j + 1
		tiles = append(tiles, tile)
		tile = generator(j)
	}
	return i + 1, tiles
}
