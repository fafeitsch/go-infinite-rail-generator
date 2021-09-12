package noise

import (
	"fmt"
	"github.com/fafeitsch/go-infinite-rail-generator/domain"
)

func (n *Noise) Generate(hectometer int) *domain.Tile {

	tile := n.createRndTile(hectometer)

	if tile.Station {
		start, stationTiles := stationWidth(hectometer, n.createRndTile)
		for _, stationTile := range stationTiles {
			fmt.Printf("%f, ", stationTile.seed)
		}
		fmt.Printf("\n")
		buildStation(stationTiles)
		tile = stationTiles[hectometer-start]
	}

	right := n.createRndTile(hectometer + 1)
	tile.fixNecessarySwitches(right)
	left := n.createRndTile(hectometer - 1)
	left.fixNecessarySwitches(tile)
	tile.fixLeftSideBumpers(left)

	return tile.Tile
}

func (n *Noise) createRndTile(hectometer int) *rndTile {
	tile := &rndTile{
		Tile: domain.NewTile(n.Seed, n.numberOfTracks(hectometer)),
		seed: n.interpolate(hectometer),
	}
	tile.Station = n.derive(2500).isStation(hectometer)
	return tile
}
