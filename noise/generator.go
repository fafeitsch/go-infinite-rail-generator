package noise

import (
	"github.com/fafeitsch/go-infinite-rail-generator/domain"
	"math"
	"math/rand"
)

type tileGenerator struct {
	random         rand.Rand
	tracks         int
	nextTileTracks int
}

func (n *Noise) Generate(hectometer int) domain.Tile {
	seed := n.Interpolate(hectometer)
	source := rand.NewSource(int64(seed * 10e10))
	random := rand.New(source)

	neighborSeed := n.Interpolate(hectometer + 1)
	generator := tileGenerator{
		random:         *random,
		tracks:         computeTracks(seed),
		nextTileTracks: computeTracks(neighborSeed),
	}
	switches := generator.mandatorySwitches()
	tracks := make([]domain.Track, 0, generator.tracks)
	for _, sw := range switches {
		tracks = append(tracks, domain.Track{Switches: sw})
	}
	return domain.Tile{Tracks: tracks}
}

func computeTracks(seed float64) int {
	if seed < 0.2 {
		return 1
	}
	if seed < 0.6 {
		return 2
	}
	if seed < 0.7 {
		return 3
	}
	return int(seed*10 - 3)
}

func (t *tileGenerator) mandatorySwitches() [][]*domain.Switch {
	if t.tracks > t.nextTileTracks {
		return mergingSwitches(t.tracks, t.nextTileTracks)
	} else {
		return divergingSwitches(t.tracks, t.nextTileTracks)
	}
}

func mergingSwitches(left int, right int) [][]*domain.Switch {
	result := make([][]*domain.Switch, left)
	for track := 0; track < left; track++ {
		result[track] = make([]*domain.Switch, 0, 0)
	}
	offsetLeft := left / 2
	offsetRight := -right / 2
	for track := 0; track < left; track++ {
		position := track - offsetLeft
		if position < offsetRight {
			result[track] = append(result[track], &domain.Switch{Direction: domain.Merging, TrackSpan: int(math.Abs(float64(offsetRight - position)))})
		} else if position > offsetRight+right-1 {
			result[track] = append(result[track], &domain.Switch{Direction: domain.Merging, TrackSpan: -(position - (offsetRight + right - 1))})
		}
	}
	return result
}

func divergingSwitches(left int, right int) [][]*domain.Switch {
	result := make([][]*domain.Switch, left)
	for track := 0; track < left; track++ {
		result[track] = make([]*domain.Switch, 0, 0)
	}
	offsetLeft := -left / 2
	offsetRight := right / 2
	for track := 0; track < right; track++ {
		position := track - offsetRight
		if position < offsetLeft {
			result[0] = append(result[0], &domain.Switch{Direction: domain.Diverging, TrackSpan: -int(math.Abs(float64(offsetLeft) - float64(position)))})
		} else if position > offsetLeft+left-1 {
			result[left-1] = append(result[left-1], &domain.Switch{Direction: domain.Diverging, TrackSpan: position - (offsetLeft + left - 1)})
		}
	}
	return result
}
