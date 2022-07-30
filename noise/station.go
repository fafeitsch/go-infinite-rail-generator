package noise

import (
	_ "embed"
	"github.com/fafeitsch/go-infinite-rail-generator/domain"
	"strings"
)

//go:embed default_town_names.txt
var defaultTownNames string

func (g *Generator) buildStation(hectometer int) *rndTile {
	start, tiles := g.stationWidth(hectometer)
	index := hectometer - start
	relevantTile := tiles[hectometer-start]
	relevantTile.replaceTrackLayout(tiles[0].Tracks.AlphaTracks())
	for i, connected := range tiles[0].Tracks.AlphaTracks() {
		if !connected {
			continue
		}
		if !relevantTile.Platforms.Alpha[i] && !relevantTile.Platforms.Alpha[i+1] {
			if tiles[0].createRandom(98).Float64() < 0.5 && i%2 == 1 {
				relevantTile.Platforms.Alpha[i] = true
				relevantTile.Platforms.Beta[i] = true
				relevantTile.Platforms.Gamma[i] = index != len(tiles)
			} else {
				relevantTile.Platforms.Alpha[i+1] = true
				relevantTile.Platforms.Beta[i+1] = true
				relevantTile.Platforms.Gamma[i+1] = index != len(tiles)
			}
		}
	}
	if index == len(tiles)/2 {
		relevantTile.StationName = g.getStationName(relevantTile)
	}
	return relevantTile
}

func (g *Generator) stationWidth(start int) (int, []*rndTile) {
	generator := g.createRndTile
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

func (r *rndTile) replaceTrackLayout(template [16]bool) {
	for i, connected := range template {
		if !connected {
			r.Tracks.Alpha[i] = nil
			r.Tracks.Beta[i] = nil
			r.Tracks.Gamma[i] = nil
			continue
		}
		r.Tracks.Alpha[i] = []*domain.Connector{
			{
				Target: domain.Beta,
				Slot:   i,
			},
		}
		r.Tracks.Beta[i] = []*domain.Connector{
			{
				Target: domain.Gamma,
				Slot:   i,
			},
		}
		r.Tracks.Gamma[i] = []*domain.Connector{
			{
				Target: domain.Omega,
				Slot:   i,
			},
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
