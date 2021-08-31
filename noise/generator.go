package noise

import (
	"github.com/fafeitsch/go-infinite-rail-generator/domain"
)

func (n *Noise) Generate(hectometer int) *domain.Tile {
	tile := rndTile{
		Tile: domain.NewTile(n.Seed, n.numberOfTracks(hectometer)),
		seed: n.interpolate(hectometer),
	}
	right := rndTile{
		Tile: domain.NewTile(n.Seed, n.numberOfTracks(hectometer+1)),
		seed: n.interpolate(hectometer + 1),
	}
	tile.fixNecessarySwitches(right)
	left := rndTile{
		Tile: domain.NewTile(n.Seed, n.numberOfTracks(hectometer-1)),
		seed: n.interpolate(hectometer - 1),
	}
	left.fixNecessarySwitches(tile)
	tile.fixLeftSideBumpers(left)
	return tile.Tile
}
