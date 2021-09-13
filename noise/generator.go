package noise

import (
	"github.com/fafeitsch/go-infinite-rail-generator/domain"
)

func (g *Generator) Generate(hectometer int) *domain.Tile {

	tile := g.createRndTile(hectometer)

	if tile.station {
		start, stationTiles := stationWidth(hectometer, g.createRndTile)
		g.buildStation(stationTiles)
		tile = stationTiles[hectometer-start]
	}

	right := g.createRndTile(hectometer + 1)
	tile.fixNecessarySwitches(right)
	left := g.createRndTile(hectometer - 1)
	left.fixNecessarySwitches(tile)
	tile.fixLeftSideBumpers(left)

	return tile.Tile
}

func (g *Generator) createRndTile(hectometer int) *rndTile {
	tile := &rndTile{
		Tile: domain.NewTile(g.Seed, g.noise.numberOfTracks(hectometer)),
		seed: g.noise.interpolate(hectometer),
	}
	tile.station = g.derive(2500).isStation(hectometer)
	return tile
}
