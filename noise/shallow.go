package noise

import (
	"github.com/fafeitsch/go-infinite-rail-generator/domain"
	"math/rand"
)

type rndTile struct {
	*domain.Tile
	seed    float64
	station bool
}

type side int

const (
	leftSide  side = 43
	rightSide      = 73
)

func (r *rndTile) createRandom(shift int) *rand.Rand {
	return rand.New(rand.NewSource(int64(r.seed*10e10) + int64(shift)))
}

func (r *rndTile) fixNecessarySwitches(right *rndTile) {
	rightConnectors := make(map[int]bool)
	for i, connectors := range right.Tracks.Alpha {
		rightConnectors[i] = len(connectors) > 0
	}
	for i, track := range r.Tracks.Gamma {
		underTest := track.FindConnector(domain.Omega, i)
		var j int
		if underTest != nil && !rightConnectors[i] {
			if r.createRandom(i+rightSide).Float64() < 0.3 {
				r.Tracks.Gamma[i] = track[:0]
				continue
			} else if i <= len(r.Tracks.Gamma)/2 {
				j = i + 1
				for !rightConnectors[j] && j <= len(r.Tracks.Gamma)-1 {
					j = j + 1
				}
			} else {
				j = i - 1
				for !rightConnectors[j] && j >= 0 {
					j = j - 1
				}
			}
			underTest.Slot = j
		} else if underTest == nil && rightConnectors[i] {
			if right.createRandom(i+int(leftSide)).Float64() < 0.3 {
				continue
			}
			if i <= len(r.Tracks.Gamma)/2 {
				j = i + 1
				for len(r.Tracks.Gamma[j]) == 0 {
					j = j + 1
				}
			} else {
				j = i - 1
				for len(r.Tracks.Gamma[j]) == 0 {
					j = j - 1
				}
			}
			r.Tracks.Gamma[j] = append(r.Tracks.Gamma[j], &domain.Connector{Target: domain.Omega, Slot: i})
		}
	}
}

func (r *rndTile) fixLeftSideBumpers(left *rndTile) {
	leftConnectors := left.Tracks.BuildConnectorMap(domain.Gamma, domain.Omega)
	for i, connectors := range r.Tracks.Alpha {
		hasRightConnector := len(connectors) > 0
		if !leftConnectors[i] && hasRightConnector {
			r.Tracks.Alpha[i] = connectors[:0]
		}
	}
}
